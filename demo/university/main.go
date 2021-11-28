package main

import (
	"flag"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

type Verifier struct{}

func (Verifier) CreatePresentationRequest() verifier.PresentationRequest {
	return verifier.PresentationRequest{
		Purpose: "Logs first and last name.",
		Issuer:  "university-issuer",
		Verifier: common.Signature{
			DID: "university-verifier",
		},
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Println("Verified:", cred.Credentials["FirstName"], cred.Credentials["LastName"])
	return nil
}

type Issuer struct{}

func (Issuer) CreateVerifiableCredentials(iss *issuer.IssueRequest) (*common.VerifiableCredential, error) {
	return &common.VerifiableCredential{
		Credentials: iss.Fields,
		Subject:     iss.Subject,
	}, nil
}

func main() {
	port := flag.Int("port", 8082, "port to run the server on")
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
