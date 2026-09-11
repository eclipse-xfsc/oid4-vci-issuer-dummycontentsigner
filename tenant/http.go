package tenant

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const tenantPathPrefix = "/internal/tenants/"

func RegisterHTTPHandlers(mux *http.ServeMux, registry *Registry) {
	mux.HandleFunc(tenantPathPrefix, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		tenantID := strings.TrimPrefix(r.URL.Path, tenantPathPrefix)

		slog.Info("tenant request received",
			"method", r.Method,
			"path", r.URL.Path,
			"tenant_id", tenantID,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		if tenantID == "" || strings.Contains(tenantID, "/") {
			slog.Warn("invalid tenant path",
				"method", r.Method,
				"path", r.URL.Path,
				"tenant_id", tenantID,
			)

			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "tenant not found",
			})
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
			slog.Warn("unsupported tenant request method",
				"method", r.Method,
				"tenant_id", tenantID,
			)

			w.Header().Set("Allow", "PUT, GET, DELETE")
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "method not allowed",
			})
		}

		slog.Info("tenant request finished",
			"method", r.Method,
			"path", r.URL.Path,
			"tenant_id", tenantID,
			"duration", time.Since(start),
		)
	})
}

func upsertTenant(
	w http.ResponseWriter,
	r *http.Request,
	registry *Registry,
	tenantID string,
) {
	defer r.Body.Close()

	slog.Info("upserting tenant configuration",
		"tenant_id", tenantID,
	)

	var request UpsertRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		slog.Warn("failed to decode tenant configuration",
			"tenant_id", tenantID,
			"error", err,
		)

		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	slog.Debug("tenant configuration decoded",
		"tenant_id", tenantID,
		"credential_issuer", request.CredentialIssuer,
		"credential_endpoint", request.CredentialEndpoint,
		"nonce_endpoint", request.NonceEndpoint,
		"status_endpoint", request.StatusEndpoint,
		"authorization_servers", request.AuthorizationServers,
	)

	if request.CredentialIssuer == "" ||
		request.CredentialEndpoint == "" ||
		len(request.AuthorizationServers) == 0 {

		slog.Warn("tenant configuration validation failed",
			"tenant_id", tenantID,
			"credential_issuer_set", request.CredentialIssuer != "",
			"credential_endpoint_set", request.CredentialEndpoint != "",
			"authorization_server_count", len(request.AuthorizationServers),
		)

		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "credentialIssuer, authorizationServers and credentialEndpoint are required",
		})
		return
	}

	config := Config{
		TenantID:             tenantID,
		CredentialIssuer:     request.CredentialIssuer,
		AuthorizationServers: request.AuthorizationServers,
		CredentialEndpoint:   request.CredentialEndpoint,
		NonceEndpoint:        request.NonceEndpoint,
		StatusEndpoint:       request.StatusEndpoint,
		SchemaEndpoint:       request.SchemaEndpoint,
	}

	registry.Put(config)

	slog.Info("tenant configuration stored",
		"tenant_id", tenantID,
		"credential_issuer", config.CredentialIssuer,
		"credential_endpoint", config.CredentialEndpoint,
		"authorization_servers", config.AuthorizationServers,
	)

	writeJSON(w, http.StatusOK, config)
}

func getTenant(
	w http.ResponseWriter,
	registry *Registry,
	tenantID string,
) {
	slog.Debug("looking up tenant configuration",
		"tenant_id", tenantID,
	)

	config, ok := registry.Get(tenantID)
	if !ok {
		slog.Warn("tenant configuration not found",
			"tenant_id", tenantID,
		)

		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "tenant not found",
		})
		return
	}

	slog.Info("tenant configuration found",
		"tenant_id", tenantID,
		"credential_issuer", config.CredentialIssuer,
		"credential_endpoint", config.CredentialEndpoint,
	)

	writeJSON(w, http.StatusOK, config)
}

func deleteTenant(
	w http.ResponseWriter,
	registry *Registry,
	tenantID string,
) {
	slog.Info("deleting tenant configuration",
		"tenant_id", tenantID,
	)

	if !registry.Delete(tenantID) {
		slog.Warn("tenant configuration could not be deleted because it does not exist",
			"tenant_id", tenantID,
		)

		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "tenant not found",
		})
		return
	}

	slog.Info("tenant configuration deleted",
		"tenant_id", tenantID,
	)

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("failed to write JSON response",
			"status", status,
			"error", err,
		)
	}
}
