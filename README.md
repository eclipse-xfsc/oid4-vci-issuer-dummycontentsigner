# Description

This is a dummy issuing module to demonstrate how a issuing functionality can be constructed. This module serves an issuing endpoint over nats to the issuing frame and an creation nats endpoint to the cPCM. The module contains an in memory credential storage which prepares the credentials according to authorization code of the offering for later pickup. The credential itself will be signed by the signer service. 

# Capabilities

- Prepares dummy credentials in internal storage for issuance. 
- Uses TSA Signer Service to sign credentials
- Provides metadata for two credential types, one for JSON-LD one for SD-JWT
- Provides Nats interface to pickup offering links