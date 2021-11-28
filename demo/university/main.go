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
		Name:        "Sample Verifier",
		Purpose:     "Logs first and last name.",
		VerifierDID: "university-verifier.json",
	}
}

func (Verifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Println("Verified:", cred.Credentials["FirstName"], cred.Credentials["LastName"])
	return nil
}

type Issuer struct{}

func (Issuer) CreateVerifiableCredentials(iss *issuer.IssueRequest) (*common.VerifiableCredential, error) {
	return &common.VerifiableCredential{
		SubjectDID:  iss.SubjectDID,
		Credentials: iss.Fields,
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
			DID:           "university-issuer.json",
			PrivateKeyURI: "university/keys/issuer.private.key",
		},
	}
	server.RunServer(*port)
}
