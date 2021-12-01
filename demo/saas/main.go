package main

import (
	"fmt"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

const ISSUER_DID = "did:example:bd395203-9b81-4808-b259-7ff410aa7f73"
const VERIFIER_DID = "did:example:41766f26-de13-4c9f-b9f2-aa51f189f6d1"
const CRED_TYPE = "Account Credentials"

const port = 8085

type Verifier struct{}

func (Verifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		ServiceURL:  fmt.Sprintf("http://localhost:%d/verify", port),
		EntityName:  "SaaS Login",
		CredType:    CRED_TYPE,
		Description: "Creates a new session for the provided account",
		Issuer:      ISSUER_DID,
		Entity: common.Signature{
			DID: VERIFIER_DID,
		},
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Println("(Verifier) New session:", cred.Credentials["Username"])
	return nil
}

type Issuer struct{}

func (Issuer) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Type:        "iss:form",
		ServiceURL:  fmt.Sprintf("http://localhost:%d/issue", port),
		EntityName:  "SaaS Account Creator",
		CredType:    CRED_TYPE,
		Description: "Create a new account for SaaS.",
		Fields: []common.PresentationField{
			{
				Name: "Username",
				Type: "text",
			},
			{
				Name: "First Name",
				Type: "text",
			},
			{
				Name: "Last Name",
				Type: "text",
			},
		},
		Entity: common.Signature{
			DID: ISSUER_DID,
		},
	}
}

func (Issuer) CreateVerifiableCredentials(cred *common.VerifiableCredential) error {
	log.Println("(Issuer) Account Created:", cred.Credentials["Username"])

	cred.CredType = CRED_TYPE
	return nil
}

func main() {
	server := demo.DemoServer{
		PublicURL: "./saas/public",
		VerifierService: verifier.VerifierService{
			Verifier:      Verifier{},
			PrivateKeyURI: "saas/keys/verifier.private.key",
		},
		IssuerService: issuer.IssuerService{
			Issuer:        Issuer{},
			DID:           ISSUER_DID,
			PrivateKeyURI: "saas/keys/issuer.private.key",
		},
	}
	server.RunServer(port)
}
