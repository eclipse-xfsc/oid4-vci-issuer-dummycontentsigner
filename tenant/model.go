package tenant

type Config struct {
	TenantID             string   `json:"tenantId"`
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
	StatusEndpoint       *string  `json:"statusEndpoint"`
}

type UpsertRequest struct {
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
	StatusEndpoint       *string  `json:"statusEndpoint"`
}
