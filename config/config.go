package config

import cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"

type Config struct {
	Nats                 cloudeventprovider.NatsConfig `envconfig:"NATS"`
	Origin               string                        `envconfig:"ORIGIN"`
	SignerUrl            string                        `envconfig:"SIGNERURL"`
	SignerKey            string                        `envconfig:"SIGNERKEY"`
	Credential_Issuer    string                        `envconfig:"CREDENTIAL_ISSUER"`
	Authorization_Server []string                      `envconfig:"AUTHORIZATION_SERVER"`
	Credential_Endpoint  string                        `envconfig:"CREDENTIAL_ENDPOINT"`
}
