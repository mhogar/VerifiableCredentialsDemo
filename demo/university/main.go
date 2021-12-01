package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

const ISSUER_DID = "did:example:e98e0ae2-5096-4de5-8096-97df8e50cf41"
const VERIFIER_DID = "did:example:189a2384-cc88-4a72-8c70-c2b7aedda6b8"
const CRED_TYPE = "Student ID Card"

var port int

type Verifier struct{}

func (Verifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		ServiceURL:  fmt.Sprintf("http://localhost:%d/verify", port),
		EntityName:  "University Verifier",
		CredType:    CRED_TYPE,
		Description: "Records student information for the purpose of attendance tracking for the exam.",
		Issuer:      ISSUER_DID,
		Entity: common.Signature{
			DID: VERIFIER_DID,
		},
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Printf("(Verifier) Verified: %s %s, %s %s", cred.Credentials["First Name"], cred.Credentials["Last Name"], cred.Credentials["Student Number"], cred.Credentials["Email"])
	return nil
}

type Issuer struct{}

func (Issuer) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Type:        "iss:form",
		ServiceURL:  fmt.Sprintf("http://localhost:%d/issue", port),
		EntityName:  "University Issuer",
		CredType:    CRED_TYPE,
		Description: "Authenticate using login information to create a student ID card.",
		Fields: []common.PresentationField{
			{
				Name: "Username",
				Type: "text",
			},
			{
				Name: "Password",
				Type: "password",
			},
		},
		Entity: common.Signature{
			DID: ISSUER_DID,
		},
	}
}

func (Issuer) CreateVerifiableCredentials(cred *common.VerifiableCredential) error {
	log.Println("(Issuer) Login attempt:", cred.Credentials["Username"])

	//NOTE: in real environment would verify login properly
	if cred.Credentials["Username"] != "username" || cred.Credentials["Password"] != "password" {
		return errors.New("invalid username and/or password")
	}

	cred.CredType = CRED_TYPE
	cred.Credentials = map[string]string{
		"First Name":     "Alice",
		"Last Name":      "Student",
		"Student Number": "0123456",
		"Email":          "alice@university.ca",
	}
	return nil
}

func main() {
	port_ptr := flag.Int("port", 8084, "port to run the server on")
	flag.Parse()

	port = *port_ptr

	server := demo.DemoServer{
		PublicURL: "./university/public",
		VerifierService: verifier.VerifierService{
			Verifier:      Verifier{},
			PrivateKeyURI: "university/keys/verifier.private.key",
		},
		IssuerService: issuer.IssuerService{
			Issuer:        Issuer{},
			DID:           ISSUER_DID,
			PrivateKeyURI: "university/keys/issuer.private.key",
		},
	}
	server.RunServer(port)
}
