# SigGraph Smart Contract
## Description
This repo contains the hyperledger implementation as according to the specification from [SigGraph repo link].

The development will be divided into N phases:
1. Support for asset creation, query and transfer.
2. Support for asset split, merge and fusion.
3. Support for certificate authority and certificate.
4. Support for remaining types
5. Support for versioning
6. Version voting.
7. Cross network connection and asset transfer.

## Architecture
The project is divided into 3 layers: view, controller and service

- View: entry point for request from client. Normalize data if necessary before passing to controller
- Controller: main logic layer.
- Service: external service, data access layer or common functionality.

