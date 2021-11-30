package main

import (
	"flag"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

const issuerDID = "university-issuer"
const verifierDID = "university-verifier"

type Verifier struct{}

func (Verifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Purpose: "Logs first and last name.",
		Issuer:  issuerDID,
		Entity: common.Signature{
			DID: verifierDID,
		},
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Printf("Verified: %s %s, %s %s", cred.Credentials["FirstName"], cred.Credentials["LastName"], cred.Credentials["StudentNumber"], cred.Credentials["Email"])
	return nil
}

type Issuer struct{}

func (Issuer) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Type:    "iss:form",
		Purpose: "Create ID token from login credentials",
		Fields: map[string]string{
			"Username": "",
			"Password": "",
		},
		Entity: common.Signature{
			DID: issuerDID,
		},
	}
}

func (Issuer) CreateVerifiableCredentials(cred *common.VerifiableCredential) error {
	//NOTE: in real environment would verify login
	log.Println("Username:", cred.Credentials["Username"])

	cred.Credentials = map[string]string{
		"FirstName":     "Alice",
		"LastName":      "Student",
		"StudentNumber": "123456",
		"Email":         "alice@university.ca",
	}
	return nil
}

func main() {
	port := flag.Int("port", 8084, "port to run the server on")
	flag.Parse()

	server := demo.DemoServer{
		PublicURL: "./university/public",
		VerifierService: verifier.VerifierService{
			Verifier:      Verifier{},
			PrivateKeyURI: "university/keys/verifier.private.key",
		},
		IssuerService: issuer.IssuerService{
			Issuer:        Issuer{},
			DID:           "university-issuer",
			PrivateKeyURI: "university/keys/issuer.private.key",
		},
	}
	server.RunServer(*port)
}
