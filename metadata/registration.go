package metadata

import (
	"encoding/json"
	"log"
	"time"

	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	messaging "github.com/eclipse-xfsc/nats-message-library"
	"github.com/eclipse-xfsc/nats-message-library/common"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/config"
	"github.com/eclipse-xfsc/oid4-vci-vp-library/model/credential"
	"github.com/google/uuid"
)

const Credential_Identifier = "DeveloperCredential"
const Credential_Identifier2 = "SDJWTCredential"

var vct = "SD_JWT_DEVELOPER_CREDENTIAL"

var Registration = messaging.IssuerRegistration{
	Request: common.Request{
		TenantId:  "tenant_space",
		RequestId: uuid.NewString(),
	},
	Issuer: credential.IssuerMetadata{
		CredentialIssuer:     "https://cloud-wallet.xfsc.dev",
		AuthorizationServers: []string{"https://auth-cloud-wallet.xfsc.dev/realms/master", "https://cloud-wallet.xfsc.dev"},
		CredentialEndpoint:   "https://cloud-wallet.xfsc.dev/api/issuance/credential",
		CredentialResponseEncryption: credential.CredentialRespEnc{
			EncryptionRequired: false,
		},
		Display: []credential.LocalizedCredential{
			{Name: "Example Issuer", Locale: "en-US"},
			{Name: "Beispiel Issuer", Locale: "de-DE"},
		},
		CredentialIdentifiersSupported: true,
		CredentialConfigurationsSupported: map[string]credential.CredentialConfiguration{
			Credential_Identifier: {
				Format:                               "ldp_vc",
				CryptographicBindingMethodsSupported: []string{"did:jwk"},
				CredentialSigningAlgValuesSupported:  []string{"ES256"},
				CredentialDefinition: credential.CredentialDefinition{
					Context: []string{"https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"},
					Type:    []string{"VerifiableCredential", "UniversityDegreeCredential"},
					CredentialSubject: map[string]credential.CredentialSubject{
						"given_name": {
							Display: []credential.Display{credential.Display{
								Name:   "Given Name",
								Locale: "en-US",
							}},
						},
						"family_name": {
							Display: []credential.Display{credential.Display{
								Name:   "Surname",
								Locale: "en-US",
							},
							}},
					},
				},
				ProofTypesSupported: map[credential.ProofVariant]credential.ProofType{
					"ldp_vc": {
						ProofSigningAlgValuesSupported: []string{"ES256"},
					},
				},
				Display: []credential.LocalizedCredential{
					{
						Name:   "Developer Credential",
						Locale: "en-US",
						Logo: credential.DescriptiveURL{
							URL:             "https://www.eclipse.org/eclipse.org-common/themes/solstice/public/images/logo/eclipse-foundation-grey-orange.svg",
							AlternativeText: "Eclipse Foundation Logo",
						},
						BackgroundColor: "#FFFFFF",
						TextColor:       "#000000",
					},
					{
						Name:   "Developer Credential",
						Locale: "de-DE",
						Logo: credential.DescriptiveURL{
							URL:             "https://www.eclipse.org/eclipse.org-common/themes/solstice/public/images/logo/eclipse-foundation-grey-orange.svg",
							AlternativeText: "Eclipse Foundation Logo",
						},
						BackgroundColor: "#FFFFFF",
						TextColor:       "#000000",
					},
				},
				Schema: map[string]interface{}{
					"data": map[string]interface{}{
						"$schema":     "https://json-schema.org/draft/2020-12/schema",
						"$id":         "https://example.com/developercredential.schema.json",
						"title":       "Developer Credential",
						"description": "A product from Acme's catalog",
						"type":        "object",
						"properties": map[string]interface{}{
							"given_name": map[string]interface{}{
								"description": "The unique identifier for a product",
								"type":        "string",
							},
							"family_name": map[string]interface{}{
								"description": "Name of the product",
								"type":        "string",
							},
						},
					},
					"ui": map[string]interface{}{
						"ui:order": []string{"given_name", "family_name"},
					},
				},
				Subject: "issuer.dummycontentsigner",
			},
			Credential_Identifier2: {
				Format:                               "vc+sd-jwt",
				CryptographicBindingMethodsSupported: []string{"did:jwk"},
				CredentialSigningAlgValuesSupported:  []string{"ES256"},
				CredentialDefinition: credential.CredentialDefinition{
					Type: []string{"VerifiableCredential", "SDJWTCredential"},
					CredentialSubject: map[string]credential.CredentialSubject{
						"given_name": {
							Display: []credential.Display{credential.Display{
								Name:   "Given Name",
								Locale: "en-US",
							}},
						},
						"family_name": {
							Display: []credential.Display{credential.Display{
								Name:   "Surname",
								Locale: "en-US",
							},
							}},
					},
				},
				Vct:                 &vct,
				ProofTypesSupported: map[credential.ProofVariant]credential.ProofType{},
				Display: []credential.LocalizedCredential{
					{
						Name:   "SDJWT Credential",
						Locale: "en-US",
						Logo: credential.DescriptiveURL{
							URL:             "https://www.eclipse.org/eclipse.org-common/themes/solstice/public/images/logo/eclipse-foundation-grey-orange.svg",
							AlternativeText: "Eclipse Foundation Logo",
						},
						BackgroundColor: "#FFFFFF",
						TextColor:       "#000000",
					},
					{
						Name:   "SDJWT Credential",
						Locale: "de-DE",
						Logo: credential.DescriptiveURL{
							URL:             "https://www.eclipse.org/eclipse.org-common/themes/solstice/public/images/logo/eclipse-foundation-grey-orange.svg",
							AlternativeText: "Eclipse Foundation Logo",
						},
						BackgroundColor: "#FFFFFF",
						TextColor:       "#000000",
					},
				},
				Schema: map[string]interface{}{
					"data": map[string]interface{}{
						"$schema":     "https://json-schema.org/draft/2020-12/schema",
						"$id":         "https://example.com/developercredential.schema.json",
						"title":       "SDJWT Credential",
						"description": "A product from Acme's catalog",
						"type":        "object",
						"properties": map[string]interface{}{
							"given_name": map[string]interface{}{
								"description": "The unique identifier for a product",
								"type":        "string",
							},
							"family_name": map[string]interface{}{
								"description": "Name of the product",
								"type":        "string",
							},
						},
					},
					"ui": map[string]interface{}{
						"ui:order": []string{"given_name", "family_name"},
					},
				},
				Subject: "issuer.dummycontentsigner",
			},
		},
	},
}

func Publish(conf config.Config) {

	if conf.Credential_Issuer != "" {
		Registration.Issuer.CredentialIssuer = conf.Credential_Issuer
	}

	if conf.Authorization_Server != nil && len(conf.Authorization_Server) > 0 {
		Registration.Issuer.AuthorizationServers = conf.Authorization_Server
	}

	if conf.Credential_Endpoint != "" {
		Registration.Issuer.CredentialEndpoint = conf.Credential_Endpoint
	}

	client, err := cloudeventprovider.New(
		cloudeventprovider.Config{Protocol: cloudeventprovider.ProtocolTypeNats, Settings: conf.Nats},
		cloudeventprovider.ConnectionTypePub,
		messaging.TopicIssuerRegistration,
	)
	if err != nil {
		panic(err)
	}

	interval := time.NewTicker(time.Second * 30)

	data, err := json.Marshal(Registration)
	if err != nil {
		panic(err)
	}

	event, err := cloudeventprovider.NewEvent("test-issuer", messaging.EventTypeIssuerRegistration, data)
	if err != nil {
		panic(err)
	}

	for {
		<-interval.C

		if err := client.Pub(event); err != nil {
			log.Printf("%+v", err)
			continue
		}

		//log.Printf("send event: %s", Registration.Issuer.CredentialIssuer)
	}
}
