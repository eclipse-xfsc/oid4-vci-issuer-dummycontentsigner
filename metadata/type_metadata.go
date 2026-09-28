package metadata

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/tenant"
)

// TypeMetadata is the SD-JWT VC Type Metadata representation exposed at the
// HTTPS URL used as vct for a tenant.
type TypeMetadata struct {
	VCT         string                `json:"vct"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Display     []TypeMetadataDisplay `json:"display,omitempty"`
	Claims      []TypeMetadataClaim   `json:"claims,omitempty"`
}

type TypeMetadataDisplay struct {
	Locale      string `json:"locale"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type TypeMetadataClaim struct {
	Path []any `json:"path"`
}

// BuildTypeMetadata builds SD-JWT VC Type Metadata from the same credential
// configuration that is published in the OID4VCI issuer metadata. The tenant
// configuration is applied first, so the returned vct is tenant-specific.
func BuildTypeMetadata(
	tenantConfig tenant.Config,
	credentialConfigurationID string,
) (*TypeMetadata, error) {
	registration := BuildRegistration(tenantConfig)

	configuration, ok := registration.Issuer.CredentialConfigurationsSupported[credentialConfigurationID]
	if !ok {
		return nil, errors.New("credential configuration not found")
	}

	if configuration.Vct == nil || strings.TrimSpace(*configuration.Vct) == "" {
		return nil, errors.New("credential configuration has no vct")
	}

	result := &TypeMetadata{
		VCT: *configuration.Vct,
	}

	if configuration.CredentialMetadata != nil {
		for _, display := range configuration.CredentialMetadata.Display {
			// Type Metadata only needs the normalized locale/name fields here.
			// Issuance-specific rendering fields stay in issuer metadata.
			if display.Locale == "" || display.Name == "" {
				continue
			}

			if result.Name == "" {
				result.Name = display.Name
			}

			result.Display = append(result.Display, TypeMetadataDisplay{
				Locale: display.Locale,
				Name:   display.Name,
			})
		}

		for _, claim := range configuration.CredentialMetadata.Claims {
			if len(claim.Path) == 0 {
				continue
			}

			path := make([]any, len(claim.Path))
			copy(path, claim.Path)

			result.Claims = append(result.Claims, TypeMetadataClaim{
				Path: path,
			})
		}
	}

	return result, nil
}

func RegisterTypeMetadataHTTPHandlers(mux *http.ServeMux, registry *tenant.Registry) {
	// The concrete public path is tenant-configured via schemaEndpoint. A root
	// fallback lets this service serve arbitrary configured schema paths while
	// the more specific /health and /internal/tenants/ handlers keep precedence.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeTypeMetadataJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "method not allowed",
			})
			return
		}

		tenantConfig, credentialConfigurationID, ok := matchTypeMetadataRequest(r, registry)
		if !ok {
			writeTypeMetadataJSON(w, http.StatusNotFound, map[string]string{
				"error": "type metadata not found",
			})
			return
		}

		metadata, err := BuildTypeMetadata(tenantConfig, credentialConfigurationID)
		if err != nil {
			writeTypeMetadataJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeTypeMetadataJSON(w, http.StatusOK, metadata)
	})
}

const (
	headerTenantID = "x-tenantid"
	headerOrigin   = "x-origin"
)

func matchTypeMetadataRequest(
	r *http.Request,
	registry *tenant.Registry,
) (tenant.Config, string, bool) {

	tenantID := strings.TrimSpace(r.Header.Get(headerTenantID))
	if tenantID == "" {
		return tenant.Config{}, "", false
	}

	tenantConfig, ok := registry.Get(tenantID)
	if !ok {
		return tenant.Config{}, "", false
	}

	origin := strings.TrimRight(
		strings.TrimSpace(r.Header.Get(headerOrigin)),
		"/",
	)

	requestPath := cleanURLPath(r.URL.Path)

	registration := BuildRegistration(tenantConfig)

	for id, configuration := range registration.Issuer.CredentialConfigurationsSupported {

		if configuration.Vct == nil ||
			strings.TrimSpace(*configuration.Vct) == "" {
			continue
		}

		vctURL, err := url.Parse(*configuration.Vct)
		if err != nil {
			continue
		}

		if cleanURLPath(vctURL.Path) != requestPath {
			continue
		}

		if origin != "" {
			expectedOrigin :=
				vctURL.Scheme + "://" + vctURL.Host

			if !strings.EqualFold(
				strings.TrimRight(origin, "/"),
				expectedOrigin,
			) {
				continue
			}
		}

		return tenantConfig, id, true
	}

	return tenant.Config{}, "", false
}

func cleanURLPath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}

func writeTypeMetadataJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
