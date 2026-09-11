package tenant

type Config struct {
	TenantID             string   `json:"tenantId"`
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
	StatusEndpoint       *string  `json:"statusEndpoint"`
	SchemaEndpoint       *string  `json:"schemaEndpoint"`
}

type UpsertRequest struct {
	CredentialIssuer     string   `json:"credentialIssuer"`
	AuthorizationServers []string `json:"authorizationServers"`
	CredentialEndpoint   string   `json:"credentialEndpoint"`
	NonceEndpoint        *string  `json:"nonceEndpoint"`
	StatusEndpoint       *string  `json:"statusEndpoint"`
	SchemaEndpoint       *string  `json:"schemaEndpoint"`
}
