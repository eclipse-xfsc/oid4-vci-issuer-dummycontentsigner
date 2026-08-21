package tenant

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPutGetDeleteTenant(t *testing.T) {
	registry := NewRegistry()
	mux := http.NewServeMux()
	RegisterHTTPHandlers(mux, registry)

	body := `{"credentialIssuer":"https://tenant-a.example.org","authorizationServers":["https://tenant-a.example.org/auth"],"credentialEndpoint":"https://tenant-a.example.org/api/issuance/credential"}`
	put := httptest.NewRequest(http.MethodPut, "/internal/tenants/tenant-a", strings.NewReader(body))
	putResult := httptest.NewRecorder()
	mux.ServeHTTP(putResult, put)
	if putResult.Code != http.StatusOK {
		t.Fatalf("PUT returned %d: %s", putResult.Code, putResult.Body.String())
	}

	get := httptest.NewRequest(http.MethodGet, "/internal/tenants/tenant-a", nil)
	getResult := httptest.NewRecorder()
	mux.ServeHTTP(getResult, get)
	if getResult.Code != http.StatusOK {
		t.Fatalf("GET returned %d: %s", getResult.Code, getResult.Body.String())
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/internal/tenants/tenant-a", nil)
	deleteResult := httptest.NewRecorder()
	mux.ServeHTTP(deleteResult, deleteRequest)
	if deleteResult.Code != http.StatusNoContent {
		t.Fatalf("DELETE returned %d", deleteResult.Code)
	}
}

func TestPutTenantValidatesCoreData(t *testing.T) {
	registry := NewRegistry()
	mux := http.NewServeMux()
	RegisterHTTPHandlers(mux, registry)

	request := httptest.NewRequest(http.MethodPut, "/internal/tenants/tenant-a", strings.NewReader(`{"credentialIssuer":"https://tenant-a.example.org"}`))
	result := httptest.NewRecorder()
	mux.ServeHTTP(result, request)

	if result.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", result.Code)
	}
}
