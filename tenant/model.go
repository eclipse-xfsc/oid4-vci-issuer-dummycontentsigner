package tenant

type Config struct {
	TenantID             string   `json:"tenantId"`
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
}

type UpsertRequest struct {
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
}
