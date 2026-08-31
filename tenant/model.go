package tenant

type Config struct {
	TenantID             string   `json:"tenantId"`
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
}

type UpsertRequest struct {
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
}
