package main

import (
	"fmt"
	"log"
	"net/http"
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

	go startHealthServer()

	//publish metadata
	go metadata.Publish(conf)

	//reply to credential request
	go issuance.CredentialReply(conf, storage)

	go issuance.CredentialRequest(conf, storage)

	wg.Wait()
}

func startHealthServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"UP"}`)); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})
	addr := fmt.Sprintf("%s:%d", conf.HttpHost, conf.HttpPort)
	log.Printf("starting health server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("health server failed: %v", err)
	}

}
