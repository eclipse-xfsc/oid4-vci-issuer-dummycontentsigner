package tenant

import "testing"

func TestRegistryPutGetDelete(t *testing.T) {
	registry := NewRegistry()
	config := Config{
		TenantID:             "tenant-a",
		CredentialIssuer:     "https://tenant-a.example.org",
		AuthorizationServers: []string{"https://tenant-a.example.org/auth"},
		CredentialEndpoint:   "https://tenant-a.example.org/api/issuance/credential",
	}

	registry.Put(config)

	stored, ok := registry.Get("tenant-a")
	if !ok {
		t.Fatal("expected tenant to exist")
	}
	if stored.CredentialIssuer != config.CredentialIssuer {
		t.Fatalf("unexpected issuer: %s", stored.CredentialIssuer)
	}
	if len(registry.List()) != 1 {
		t.Fatalf("expected one tenant, got %d", len(registry.List()))
	}
	if !registry.Delete("tenant-a") {
		t.Fatal("expected delete to return true")
	}
	if _, ok := registry.Get("tenant-a"); ok {
		t.Fatal("expected tenant to be deleted")
	}
}

func TestRegistryPutIsIdempotentByTenantID(t *testing.T) {
	registry := NewRegistry()
	registry.Put(Config{TenantID: "tenant-a", CredentialIssuer: "https://old.example.org"})
	registry.Put(Config{TenantID: "tenant-a", CredentialIssuer: "https://new.example.org"})

	if len(registry.List()) != 1 {
		t.Fatalf("expected one tenant after upsert, got %d", len(registry.List()))
	}
	stored, _ := registry.Get("tenant-a")
	if stored.CredentialIssuer != "https://new.example.org" {
		t.Fatalf("expected updated issuer, got %s", stored.CredentialIssuer)
	}
}
