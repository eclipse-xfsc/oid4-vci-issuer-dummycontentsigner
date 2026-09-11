package metadata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/tenant"
)

func TestBuildTypeMetadataUsesTenantVCTAndCredentialMetadata(t *testing.T) {
	schemaEndpoint := "https://tenant-a.example.org/schema"
	config := tenant.Config{
		TenantID:             "tenant-a",
		CredentialIssuer:     "https://tenant-a.example.org",
		AuthorizationServers: []string{"https://tenant-a.example.org"},
		CredentialEndpoint:   "https://tenant-a.example.org/api/credential",
		SchemaEndpoint:       &schemaEndpoint,
	}

	metadata, err := BuildTypeMetadata(config, CredentialIdentifier2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedVCT := "https://tenant-a.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL"
	if metadata.VCT != expectedVCT {
		t.Fatalf("unexpected vct: %s", metadata.VCT)
	}
	if metadata.Name != "SDJWT Credential" {
		t.Fatalf("unexpected name: %s", metadata.Name)
	}
	if len(metadata.Display) != 2 {
		t.Fatalf("expected 2 display entries, got %d", len(metadata.Display))
	}
	if len(metadata.Claims) != 2 {
		t.Fatalf("expected 2 claims, got %d", len(metadata.Claims))
	}
}

func TestTypeMetadataEndpointUsesConfiguredSchemaEndpoint(t *testing.T) {
	registry := tenant.NewRegistry()
	schemaEndpoint := "https://tenant-a.example.org/schema"
	registry.Put(tenant.Config{
		TenantID:             "tenant-a",
		CredentialIssuer:     "https://tenant-a.example.org",
		AuthorizationServers: []string{"https://tenant-a.example.org"},
		CredentialEndpoint:   "https://tenant-a.example.org/api/credential",
		SchemaEndpoint:       &schemaEndpoint,
	})

	mux := http.NewServeMux()
	RegisterTypeMetadataHTTPHandlers(mux, registry)

	req := httptest.NewRequest(
		http.MethodGet,
		"https://tenant-a.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL",
		nil,
	)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d, body: %s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: %s", got)
	}

	var response TypeMetadata
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	expectedVCT := "https://tenant-a.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL"
	if response.VCT != expectedVCT {
		t.Fatalf("unexpected vct: %s", response.VCT)
	}
}

func TestTypeMetadataEndpointResolvesTenantsByHost(t *testing.T) {
	registry := tenant.NewRegistry()

	for _, host := range []string{"tenant-a.example.org", "tenant-b.example.org"} {
		schemaEndpoint := "https://" + host + "/schema"
		registry.Put(tenant.Config{
			TenantID:             host,
			CredentialIssuer:     "https://" + host,
			AuthorizationServers: []string{"https://" + host},
			CredentialEndpoint:   "https://" + host + "/api/credential",
			SchemaEndpoint:       &schemaEndpoint,
		})
	}

	mux := http.NewServeMux()
	RegisterTypeMetadataHTTPHandlers(mux, registry)

	req := httptest.NewRequest(
		http.MethodGet,
		"https://tenant-b.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL",
		nil,
	)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d, body: %s", recorder.Code, recorder.Body.String())
	}

	var response TypeMetadata
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if response.VCT != "https://tenant-b.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL" {
		t.Fatalf("wrong tenant metadata returned: %s", response.VCT)
	}
}
