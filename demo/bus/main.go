package main

import (
	"fmt"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

const ISSUER_DID = "did:example:d2f54564-cbf4-4574-904f-a49e3a6a2f1f"
const VERIFIER_DID = "did:example:c6970460-f6b0-4eaa-9e96-75418fb8c4f9"
const CRED_TYPE = "Bus Pass"

const port = 8086

type Issuer struct{}

func (Issuer) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Type:        "iss:cred",
		ServiceURL:  fmt.Sprintf("http://localhost:%d/issue", port),
		EntityName:  "Bus Pass Creator",
		CredType:    CRED_TYPE,
		Description: "Create a new bus pass using your student ID card.",
		Issuer:      "did:example:e98e0ae2-5096-4de5-8096-97df8e50cf41", //university issuer
		Entity: common.Signature{
			DID: ISSUER_DID,
		},
	}
}

func (Issuer) CreateVerifiableCredentials(cred *common.VerifiableCredential) (*common.VerifiableCredential, error) {
	firstName := cred.Credentials["First Name"]
	lastName := cred.Credentials["Last Name"]

	log.Println("(Issuer) Bus Pass Created:", firstName, lastName)

	busCreds := common.VerifiableCredential{
		CredType: CRED_TYPE,
		Credentials: map[string]string{
			"First Name":      firstName,
			"Last Name":       lastName,
			"Expiration Date": "05-01-2022",
		},
		Subject: cred.Subject,
		Issuer: common.Signature{
			DID: ISSUER_DID,
		},
	}

	return &busCreds, nil
}

type Verifier struct{}

func (Verifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		ServiceURL:  fmt.Sprintf("http://localhost:%d/verify/check", port),
		EntityName:  "Bus Pass Check",
		CredType:    CRED_TYPE,
		Description: "Verifies the bus pass is valid and not expired.",
		Issuer:      ISSUER_DID,
		Entity: common.Signature{
			DID: VERIFIER_DID,
		},
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Println("(Verifier) Bus pass checked:", cred.Credentials["First Name"], cred.Credentials["Last Name"])
	return nil
}

func main() {
	server := demo.DemoServer{
		PublicURL: "./bus/public",
		IssuerService: issuer.IssuerService{
			Issuer:        Issuer{},
			DID:           ISSUER_DID,
			PrivateKeyURI: "bus/keys/issuer.private.key",
		},
		VerifierServices: map[string]verifier.VerifierService{
			"check": {
				Verifier:      Verifier{},
				PrivateKeyURI: "bus/keys/verifier.private.key",
			},
		},
	}
	server.RunServer(port)
}
