package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	messaging "github.com/eclipse-xfsc/nats-message-library"
	"github.com/eclipse-xfsc/nats-message-library/common"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/metadata"
	issuance "github.com/eclipse-xfsc/oid4-vci-issuer-service/pkg/messaging"
	"github.com/eclipse-xfsc/oid4-vci-vp-library/model/credential"
	"github.com/google/uuid"
)

var createCredentialClient cloudeventprovider.CloudEventProvider

func createCredential() (offer *credential.CredentialOffer, err error) {
	if createCredentialClient == nil {

		createCredentialClient, err = cloudeventprovider.New(
			cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: cloudeventprovider.NatsConfig{
				Url:          "nats://localhost:4222",
				TimeoutInSec: time.Hour,
			}},
			cloudeventprovider.ConnectionTypeReq,
			"issuer.dummycontentsigner.request",
		)

		if err != nil {
			panic(err)
		}
	}

	var req = messaging.IssuanceRequest{
		Request: common.Request{
			TenantId:  "tenant_space",
			RequestId: uuid.NewString(),
		},
		Payload: map[string]interface{}{
			"given_name":  "test",
			"family_name": "test",
		},
		Identifier: metadata.Credential_Identifier,
	}

	b, _ := json.Marshal(req)

	testEvent, _ := cloudeventprovider.NewEvent("test-issuer", "issuance", b)

	ev, err := createCredentialClient.RequestCtx(context.Background(), testEvent)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println(ev)

	var rep messaging.IssuanceReply

	err = json.Unmarshal(ev.Data(), &rep)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println(rep.Offer.CredentialOffer)
	return &rep.Offer, nil
}

var issueCredentialClient cloudeventprovider.CloudEventProvider

func issueCredential(offering *credential.CredentialOffer) (err error) {

	if issueCredentialClient == nil {

		issueCredentialClient, err = cloudeventprovider.New(
			cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: cloudeventprovider.NatsConfig{
				Url:          "nats://localhost:4222",
				TimeoutInSec: time.Hour,
			}},
			cloudeventprovider.ConnectionTypeReq,
			"issuer.dummycontentsigner.issue",
		)

		if err != nil {
			panic(err)
		}
	}

	param, err := offering.GetOfferParameters()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	req := issuance.IssuanceModuleReq{
		Request: common.Request{
			TenantId:  "tenant_space",
			RequestId: uuid.NewString(),
		},
		Code:   param.Grants.PreAuthorizedCode.PreAuthorizationCode,
		Holder: "test",
	}

	b, _ := json.Marshal(req)

	testEvent, _ := cloudeventprovider.NewEvent("test-issuer", "issuance", b)

	ev, err := issueCredentialClient.RequestCtx(context.Background(), testEvent)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var rep issuance.IssuanceModuleRep

	err = json.Unmarshal(ev.Data(), &rep)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(rep.Credential)

	return nil

}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		//place credential
		_, err := createCredential()

		if err != nil {
			reader.ReadString('\n')
			continue
		}

		//issue it

		/*err = issueCredential(offer)
		if err != nil {
			reader.ReadString('\n')
			continue
		}*/
		reader.ReadString('\n')
	}
}
