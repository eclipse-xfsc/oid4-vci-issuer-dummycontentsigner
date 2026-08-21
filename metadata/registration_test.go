package metadata

import (
	"testing"

	"github.com/eclipse-xfsc/oid4-vci-issuer-dummycontentsigner/tenant"
)

func TestBuildRegistrationUsesTenantCoreData(t *testing.T) {
	config := tenant.Config{
		TenantID:             "tenant-a",
		CredentialIssuer:     "https://tenant-a.example.org",
		AuthorizationServers: []string{"https://tenant-a.example.org/auth"},
		CredentialEndpoint:   "https://tenant-a.example.org/api/issuance/credential",
	}

	registration := BuildRegistration(config)

	if registration.Request.TenantId != config.TenantID {
		t.Fatalf("unexpected tenant id: %s", registration.Request.TenantId)
	}
	if registration.Issuer.CredentialIssuer != config.CredentialIssuer {
		t.Fatalf("unexpected credential issuer: %s", registration.Issuer.CredentialIssuer)
	}
	if registration.Issuer.CredentialEndpoint != config.CredentialEndpoint {
		t.Fatalf("unexpected credential endpoint: %s", registration.Issuer.CredentialEndpoint)
	}
	if len(registration.Issuer.CredentialConfigurationsSupported) != 2 {
		t.Fatalf("expected static credential metadata to remain available")
	}
}
