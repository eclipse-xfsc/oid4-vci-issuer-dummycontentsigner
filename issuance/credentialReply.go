package issuance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	"github.com/eclipse-xfsc/nats-message-library/common"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/config"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/metadata"
	issuance "github.com/eclipse-xfsc/oid4-vci-issuer-service/pkg/messaging"
)

func signCredential(credential map[string]interface{}, tenantId string, signerkey string, url string, origin string, nonce string, format string) (any, error) {

	env := os.Getenv("DUMMYCONTENTSIGNER_STATUS")
	var err error
	status := false

	if env != "" {
		status, err = strconv.ParseBool(env)

		if err != nil {
			status = false
		}
	}

	credential["namespace"] = tenantId
	credential["group"] = ""
	credential["key"] = signerkey
	credential["status"] = status
	credential["nonce"] = nonce

	body, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-origin", origin)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil || res.StatusCode != 200 {
		if res != nil && res.StatusCode != 200 {
			b, _ := io.ReadAll(res.Body)
			return nil, errors.New("signer service returned no 200: " + string(b))
		}
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if format == "ldp_vc" {
		var post map[string]interface{}

		derr := json.NewDecoder(bytes.NewBuffer(b)).Decode(&post)

		if derr != nil {
			return nil, err
		}

		if post == nil {
			return nil, errors.New("no content could be signed")
		}

		return post, nil
	}

	return strings.Trim(strings.Replace(string(b), "\"", "", -1), "\n"), nil
}

func CredentialReply(conf config.Config, storage IssuanceStorage) {

	client, err := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypeRep,
		metadata.Registration.Issuer.CredentialConfigurationsSupported[metadata.Credential_Identifier].Subject+".issue",
	)
	if err != nil {
		panic(err)
	}

	for {
		if err := client.ReplyCtx(context.Background(), func(ctx context.Context, event event.Event) (*event.Event, error) {
			log.Printf("Event received %+v", event)
			var req issuance.IssuanceModuleReq
			err := json.Unmarshal(event.DataEncoded, &req)

			if err != nil {
				return nil, err
			}

			reply := issuance.IssuanceModuleRep{
				Reply: common.Reply{
					TenantId:  req.TenantId,
					RequestId: req.RequestId,
					GroupId:   req.GroupId,
				},
				Format: req.Format,
			}

			cred, err := storage.GetCredential(req.Code)

			if err != nil || cred == nil {
				log.Printf("Error %+v", err)
				reply.Error = &common.Error{
					Id:     "no credential found",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			if req.Format == "" {
				reply.Format = cred["format"].(string)
			}

			if err != nil {
				log.Printf("Error %+v", err)
				reply.Error = &common.Error{
					Id:     "credential-load-error",
					Status: 400,
					Msg:    err.Error(),
				}
			} else {

				if req.Holder != "" {
					cred["holder"] = req.Holder
				}

				c, err := signCredential(cred, req.TenantId, conf.SignerKey, conf.SignerUrl, conf.Origin, req.Code, reply.Format)

				if err != nil {
					return nil, err
				}

				if c == nil {
					reply.Error = &common.Error{
						Id:     "credential-load-error",
						Status: 400,
						Msg:    err.Error(),
					}
				} else {
					reply.Credential = c
				}
			}

			b, err := json.Marshal(reply)

			if err != nil {
				return nil, err
			}

			event, err = cloudeventprovider.NewEvent("test-issuer", "dummycontentsigner", b)
			if err != nil {
				return nil, err
			}

			return &event, nil
		}); err != nil {
			log.Printf("%+v", err)
			continue
		}
	}
}
