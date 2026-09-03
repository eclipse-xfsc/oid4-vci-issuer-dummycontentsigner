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

func TestBuildRegistrationDoesNotLeakVCTBetweenTenants(t *testing.T) {
	schemaA := "https://tenant-a.example.org/schema"
	schemaB := "https://tenant-b.example.org/types"

	registrationA := BuildRegistration(tenant.Config{
		TenantID:             "tenant-a",
		CredentialIssuer:     "https://tenant-a.example.org",
		AuthorizationServers: []string{"https://tenant-a.example.org"},
		CredentialEndpoint:   "https://tenant-a.example.org/api/credential",
		SchemaEndpoint:       &schemaA,
	})

	registrationB := BuildRegistration(tenant.Config{
		TenantID:             "tenant-b",
		CredentialIssuer:     "https://tenant-b.example.org",
		AuthorizationServers: []string{"https://tenant-b.example.org"},
		CredentialEndpoint:   "https://tenant-b.example.org/api/credential",
		SchemaEndpoint:       &schemaB,
	})

	vctA := registrationA.Issuer.CredentialConfigurationsSupported[CredentialIdentifier2].Vct
	vctB := registrationB.Issuer.CredentialConfigurationsSupported[CredentialIdentifier2].Vct

	if vctA == nil || *vctA != "https://tenant-a.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL" {
		t.Fatalf("unexpected tenant A vct: %v", vctA)
	}
	if vctB == nil || *vctB != "https://tenant-b.example.org/types/SD_JWT_DEVELOPER_CREDENTIAL" {
		t.Fatalf("unexpected tenant B vct: %v", vctB)
	}
}
