package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/config"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/issuance"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/metadata"
	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/tenant"
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
	tenantRegistry := tenant.NewRegistry()

	go startHTTPServer(tenantRegistry)
	go metadata.Publish(conf, tenantRegistry)
	go issuance.CredentialReply(conf, storage, tenantRegistry)
	go issuance.CredentialRequest(conf, tenantRegistry, storage)

	wg.Wait()
}

func startHTTPServer(tenantRegistry *tenant.Registry) {
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

	tenant.RegisterHTTPHandlers(mux, tenantRegistry)

	addr := fmt.Sprintf("%s:%d", conf.HttpHost, conf.HttpPort)
	log.Printf("starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
