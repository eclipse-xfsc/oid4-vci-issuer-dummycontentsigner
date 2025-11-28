package issuance

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	messaging "github.com/eclipse-xfsc/nats-message-library"
	"github.com/eclipse-xfsc/nats-message-library/common"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/config"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/metadata"
	issumsg "github.com/eclipse-xfsc/oid4-vci-issuer-service/pkg/messaging"
	"github.com/eclipse-xfsc/oid4-vci-vp-library/model/credential"
	"github.com/google/uuid"
)

func createCredential(code string, payload map[string]interface{}, storage IssuanceStorage, identifier string) error {
	var credJson = make(map[string]interface{})

	credJson = map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://w3id.org/security/suites/jws-2020/v1",
			"https://schema.org",
		},
		"issuanceDate": "2022-06-02T17:24:05.032533+03:00",
	}

	credJson["credentialSubject"] = payload

	credJson["issuer"] = metadata.Registration.Issuer.CredentialIssuer

	if identifier == metadata.Credential_Identifier2 {
		credJson["format"] = "vc+sd-jwt"
		credJson["type"] = []string{"VerifiableCredential", "SDJWTCredential"}
	} else {
		credJson["format"] = "ldp_vc"
		credJson["issuanceDate"] = time.Now().Format(time.RFC3339)
		credJson["type"] = []string{"VerifiableCredential", "DeveloperCredential"}
	}

	err := storage.AddCredential(code, credJson)

	if err != nil {
		return err
	}

	return nil
}

func CredentialRequest(conf config.Config, storage IssuanceStorage) {

	authclient, _ := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypeReq,
		issumsg.TopicOffering,
	)

	client, _ := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypeRep,
		metadata.Registration.Issuer.CredentialConfigurationsSupported[metadata.Credential_Identifier].Subject+".request",
	)

	for {
		if err := client.ReplyCtx(context.Background(), func(ctx context.Context, event event.Event) (*event.Event, error) {

			var req messaging.IssuanceRequest
			err := json.Unmarshal(event.DataEncoded, &req)

			if err != nil {
				return nil, err
			}

			reply := messaging.IssuanceReply{
				Reply: common.Reply{
					TenantId:  req.TenantId,
					RequestId: req.RequestId,
					GroupId:   req.GroupId,
				},
			}

			nonce := uuid.NewString()
			offerReq := issumsg.OfferingURLReq{
				Request: common.Request{
					TenantId:  req.TenantId,
					RequestId: req.RequestId,
					GroupId:   reply.GroupId,
				},
				Params: issumsg.AuthorizationReq{
					CredentialConfigurations: []credential.CredentialConfigurationIdentifier{
						credential.CredentialConfigurationIdentifier{
							Id: req.Identifier,
						},
					},
					GrantType: "urn:ietf:params:oauth:grant-type:pre-authorized_code",
					TwoFactor: issumsg.TwoFactor{
						Enabled: false,
					},
					Nonce: nonce,
				},
			}

			r, _ := json.Marshal(offerReq)

			authevent, err := cloudeventprovider.NewEvent("test-issuer", issumsg.EventTypeOffering, r)

			if err != nil {
				reply.Error = &common.Error{
					Id:     "auth-req-error",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			authrep, err := authclient.RequestCtx(ctx, authevent)

			if err != nil {
				reply.Error = &common.Error{
					Id:     "credential-req-error",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			if authrep != nil {

				var resp issumsg.OfferingURLResp

				err = json.Unmarshal(authrep.Data(), &resp)

				if err == nil {
					err = createCredential(resp.Code, req.Payload, storage, req.Identifier)
				}

				if err != nil {
					reply.Error = &common.Error{
						Id:     "credential-req-error",
						Status: 400,
						Msg:    err.Error(),
					}
				} else {
					reply.Offer = resp.CredentialOffer
				}
			} else {
				reply.Error = &common.Error{
					Id:     "credential-req-error",
					Status: 400,
					Msg:    "no result",
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
