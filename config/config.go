package config

import cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"

type Config struct {
	Nats                 cloudeventprovider.NatsConfig `envconfig:"NATS"`
	SignerCredentialUrl  string                        `envconfig:"SIGNER_CREDENTIAL_URL"`
	Credential_Issuer    string                        `envconfig:"CREDENTIAL_ISSUER"`
	Authorization_Server []string                      `envconfig:"AUTHORIZATION_SERVER"`
	Credential_Endpoint  string                        `envconfig:"CREDENTIAL_ENDPOINT"`
}
