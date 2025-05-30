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
	"github.com/google/uuid"
)

func createCredential(id string, payload map[string]interface{}, storage IssuanceStorage, identifier string) error {
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

	err := storage.AddCredential(id, credJson)

	if err != nil {
		return err
	}

	return nil
}

func CredentialRequest(conf config.Config, storage IssuanceStorage) {
	authclient, err := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypeReq,
		issumsg.TopicOffering,
	)
	if err != nil {
		log.Printf("failed to create authclient: %v", err)
		return
	}

	client, err := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypeRep,
		metadata.Registration.Issuer.CredentialConfigurationsSupported[metadata.Credential_Identifier].Subject+".request",
	)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return
	}

	for {
		if err := client.ReplyCtx(context.Background(), func(ctx context.Context, event event.Event) (*event.Event, error) {
			var req messaging.IssuanceRequest
			err := json.Unmarshal(event.DataEncoded, &req)
			if err != nil {
				log.Printf("failed to unmarshal request: %v", err)
				return nil, err
			}

			reply := messaging.IssuanceReply{
				Reply: common.Reply{
					TenantId:  req.TenantId,
					RequestId: req.RequestId,
				},
			}

			id := uuid.NewString()
			nonce := uuid.NewString()

			offerReq := issumsg.OfferingURLReq{
				Request: common.Request{
					TenantId:  req.TenantId,
					RequestId: req.RequestId,
				},
				Params: issumsg.AuthorizationReq{
					CredentialType:       req.Identifier,
					CredentialIdentifier: []string{id},
					GrantType:            "urn:ietf:params:oauth:grant-type:pre-authorized_code",
					TwoFactor: issumsg.TwoFactor{
						Enabled: false,
					},
					Nonce: nonce,
				},
			}

			r, err := json.Marshal(offerReq)
			if err != nil {
				log.Printf("failed to marshal offer request: %v", err)
				reply.Error = &common.Error{
					Id:     "marshal-offerreq-error",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			authevent, err := cloudeventprovider.NewEvent("test-issuer", issumsg.EventTypeOffering, r)
			if err != nil {
				log.Printf("failed to create auth event: %v", err)
				reply.Error = &common.Error{
					Id:     "auth-event-error",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			authrep, err := authclient.RequestCtx(ctx, authevent)
			if err != nil {
				log.Printf("authclient.RequestCtx failed: %v", err)
				reply.Error = &common.Error{
					Id:     "auth-request-error",
					Status: 400,
					Msg:    err.Error(),
				}
			}

			if authrep != nil {
				var resp issumsg.OfferingURLResp
				err = json.Unmarshal(authrep.Data(), &resp)
				if err != nil {
					log.Printf("failed to unmarshal auth response: %v", err)
					reply.Error = &common.Error{
						Id:     "auth-response-unmarshal-error",
						Status: 400,
						Msg:    err.Error(),
					}
				} else {
					err = createCredential(id, req.Payload, storage, req.Identifier)
					if err != nil {
						log.Printf("failed to create credential: %v", err)
						reply.Error = &common.Error{
							Id:     "create-credential-error",
							Status: 400,
							Msg:    err.Error(),
						}
					} else {
						reply.Offer = resp.CredentialOffer
					}
				}
			} else if reply.Error == nil {
				// brak odpowiedzi i brak wcześniejszego błędu
				log.Printf("auth response is nil without prior error")
				reply.Error = &common.Error{
					Id:     "auth-response-missing",
					Status: 400,
					Msg:    "no result from auth request",
				}
			}

			b, err := json.Marshal(reply)
			if err != nil {
				log.Printf("failed to marshal reply: %v", err)
				return nil, err
			}

			event, err = cloudeventprovider.NewEvent("test-issuer", "dummycontentsigner", b)
			if err != nil {
				log.Printf("failed to create final event: %v", err)
				return nil, err
			}

			return &event, nil
		}); err != nil {
			log.Printf("ReplyCtx loop error: %+v", err)
			continue
		}
	}
}
