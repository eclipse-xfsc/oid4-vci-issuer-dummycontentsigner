# Dynamic tenant registration patch

This patch removes tenant-specific issuer registration data from the static application and Helm configuration.

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
  "credentialEndpoint": "https://tenant-a.example.org/api/issuance/credential"
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
    "credentialEndpoint": "https://demo.example.org/api/issuance/credential"
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

## Removed static configuration

The following application environment variables are no longer used:

- `CREDENTIAL_ISSUER`
- `AUTHORIZATION_SERVER`
- `CREDENTIAL_ENDPOINT`

The matching Helm values have also been removed.
