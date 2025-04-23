package main

import (
	"fmt"
	"sync"

	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/config"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/issuance"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/metadata"
	"github.com/kelseyhightower/envconfig"
)

var conf config.Config

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	if err := envconfig.Process("", &conf); err != nil {
		panic(fmt.Sprintf("failed to load config from env: %+v", err))
	}

	storage := new(issuance.DummyStorage)

	//publish metadata
	go metadata.Publish(conf)

	//reply to credential request
	go issuance.CredentialReply(conf, storage)

	go issuance.CredentialRequest(conf, storage)

	wg.Wait()
}
