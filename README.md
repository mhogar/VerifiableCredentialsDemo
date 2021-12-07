# VerifiableCredentialsDemo

## Setup

### Dependencies
- Node.js (https://nodejs.org/en/)
- Golang 1.17+ (https://go.dev/)

### Installing
- `cd` into "user/client" and run `npm install`

## Running

The demo is composed of two types of applications: the user app and the verifiers/issuers

### User Application

Run the following commands in separate terminals:
- `cd` into "user" and run `go run server.main.go`. This will start the user application backend on port 8082
- `cd` into "user/client" and run `npm run serve`. This will start the user application front-end on port 8080

### Demo Applications

This demo contains three verifier/issuer services, all of which can be found under the "demo" directory. To run a demo service, for example "university", use the following command:
- `cd` into "demo" and run `go run university/main.go`. Upon running the service, it will print all of its accessible endpoints.

#### Application Specific Notes
- __SaaS__: When prompted, any non-empty values for the account fields are valid. These will be the values used in the created credential
- __University__: When prompted, the login credentials are "username" and "password". The created credential will always have the same values, hardcoded in the back-end
- __Bus__: This service involves verifying the university's signature and therefore requires that the university service is running as it hosts its own public key

## Using the Application
- Once the user application and desired demo services are running, navigate to http://localhost:8080 in a browser
- The home page of the application shows all verifiable credentials owned by the user, which are loaded from "user/wallet/verifiable-credentials.json". A fresh run application will have no credentials
- Enter the url from one of the demo services in the query field to start a request
- All DID documents for services can be found in the "blockchain" directory. This serves as a local replacement for an actual blockchain that would be used in a production environment
