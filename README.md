# Description

This is a dummy issuing module to demonstrate how a issuing functionality can be constructed. This module serves an issuing endpoint over nats to the issuing frame and an creation nats endpoint to the cPCM. The module contains an in memory credential storage which prepares the credentials according to authorization code of the offering for later pickup. The credential itself will be signed by the signer service. 

# Capabilities

- Prepares dummy credentials in internal storage for issuance. 
- Uses TSA Signer Service to sign credentials
- Provides metadata for two credential types, one for JSON-LD one for SD-JWT
- Provides Nats interface to pickup offering links

## Internal API

Create or update a tenant:

```http
PUT /internal/tenants/{tenantId}
Content-Type: application/json
```

```json
{
  "credentialIssuer": "https://tenant-a.example.org",
  "authorizationServers": [
    "https://tenant-a.example.org/auth"
  ],
  "credentialEndpoint": "https://tenant-a.example.org/api/issuance/credential",
  "nonceEndpoint": "https://tenant-a.example.org/api/nonce",
  "schemaEndpoint": "https://tenant-a.example.org/schema"
}
```

curl: 

```
curl -X PUT http://localhost:8080/internal/tenants/demo \
  -H "Content-Type: application/json" \
  -d '{
    "credentialIssuer": "https://demo.example.org",
    "authorizationServers": [
      "https://demo.example.org/auth"
    ],
    "credentialEndpoint": "https://demo.example.org/api/issuance/credential",
    "nonceEndpoint": "https://tenant-a.example.org/api/issuance/credential"
  }'
```

Read a tenant:

```http
GET /internal/tenants/{tenantId}
```

Delete a tenant:

```http
DELETE /internal/tenants/{tenantId}
```

The existing registration publisher wakes up every 30 seconds, reads the current tenant registry, combines each tenant's core data with the static credential metadata and publishes one issuer registration per tenant to NATS.

The registry is intentionally in-memory. The tenant control plane/provisioner must therefore reconcile existing tenants again after a pod restart.

## SD-JWT VC Type Metadata

For the SD-JWT credential configuration the tenant-specific `vct` is built from `schemaEndpoint` and the configured type name. For example, a tenant configuration containing:

```json
{
  "schemaEndpoint": "https://demo.example.org/schema"
}
```

results in the issuer metadata value:

```json
{
  "vct": "https://demo.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL"
}
```

The HTTP server exposes the same URL path as SD-JWT VC Type Metadata. The response is generated from the existing SD-JWT credential configuration, so claim paths and display names stay aligned with the published OID4VCI issuer metadata:

```http
GET /schema/SD_JWT_DEVELOPER_CREDENTIAL
Host: demo.example.org
```

```json
{
  "vct": "https://demo.example.org/schema/SD_JWT_DEVELOPER_CREDENTIAL",
  "name": "SDJWT Credential",
  "display": [
    {"locale": "en-US", "name": "SDJWT Credential"},
    {"locale": "de-DE", "name": "SDJWT Credential"}
  ],
  "claims": [
    {"path": ["given_name"]},
    {"path": ["family_name"]}
  ]
}
```

The endpoint returns `Content-Type: application/json`. Host matching is tenant-aware. If a reverse proxy rewrites the host, `X-Forwarded-Host` is used. A path-only fallback is accepted only when the path uniquely identifies one configured tenant.

During SD-JWT issuance the same tenant-specific `vct` is added to the credential data passed to the signer. This keeps issuer metadata, the issued credential and the Type Metadata document on the same identifier.

## Removed static configuration

The following application environment variables are no longer used:

- `CREDENTIAL_ISSUER`
- `AUTHORIZATION_SERVER`
- `CREDENTIAL_ENDPOINT`

The matching Helm values have also been removed.
