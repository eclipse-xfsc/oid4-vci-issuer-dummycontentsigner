package tenant

import (
	"encoding/json"
	"net/http"
	"strings"
)

const tenantPathPrefix = "/internal/tenants/"

func RegisterHTTPHandlers(mux *http.ServeMux, registry *Registry) {
	mux.HandleFunc(tenantPathPrefix, func(w http.ResponseWriter, r *http.Request) {
		tenantID := strings.TrimPrefix(r.URL.Path, tenantPathPrefix)
		if tenantID == "" || strings.Contains(tenantID, "/") {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "tenant not found"})
			return
		}

		switch r.Method {
		case http.MethodPut:
			upsertTenant(w, r, registry, tenantID)
		case http.MethodGet:
			getTenant(w, registry, tenantID)
		case http.MethodDelete:
			deleteTenant(w, registry, tenantID)
		default:
			w.Header().Set("Allow", "PUT, GET, DELETE")
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})
}

func upsertTenant(w http.ResponseWriter, r *http.Request, registry *Registry, tenantID string) {
	defer r.Body.Close()

	var request UpsertRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if request.CredentialIssuer == "" ||
		request.CredentialEndpoint == "" ||
		len(request.AuthorizationServers) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "credentialIssuer, authorizationServers and credentialEndpoint are required"})
		return
	}

	config := Config{
		TenantID:             tenantID,
		CredentialIssuer:     request.CredentialIssuer,
		AuthorizationServers: request.AuthorizationServers,
		CredentialEndpoint:   request.CredentialEndpoint,
		NonceEndpoint:        request.NonceEndpoint,
	}

	registry.Put(config)
	writeJSON(w, http.StatusOK, config)
}

func getTenant(w http.ResponseWriter, registry *Registry, tenantID string) {
	config, ok := registry.Get(tenantID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "tenant not found"})
		return
	}
	writeJSON(w, http.StatusOK, config)
}

func deleteTenant(w http.ResponseWriter, registry *Registry, tenantID string) {
	if !registry.Delete(tenantID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "tenant not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
