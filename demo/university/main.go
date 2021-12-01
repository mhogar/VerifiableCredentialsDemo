package main

import (
	"errors"
	"fmt"
	"log"
	"vcd/common"
	"vcd/demo"
	"vcd/issuer"
	"vcd/verifier"
)

const ISSUER_DID = "did:example:e98e0ae2-5096-4de5-8096-97df8e50cf41"
const EXAM_VERIFIER_DID = "did:example:189a2384-cc88-4a72-8c70-c2b7aedda6b8"
const EVENT_VERIFIER_DID = "did:example:383997b0-b9ec-49a9-99cc-0793c9ca4a90"
const CRED_TYPE = "Student ID Card"

const port = 8084

type Issuer struct{}

func (Issuer) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		Type:        "iss:form",
		ServiceURL:  fmt.Sprintf("http://localhost:%d/issue", port),
		EntityName:  "University Student Card Issuer",
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

func (Issuer) CreateVerifiableCredentials(cred *common.VerifiableCredential) (*common.VerifiableCredential, error) {
	log.Println("(Issuer) Login attempt:", cred.Credentials["Username"])

	//NOTE: in real environment would verify login properly
	if cred.Credentials["Username"] != "username" || cred.Credentials["Password"] != "password" {
		return nil, errors.New("invalid username and/or password")
	}

	cred.CredType = CRED_TYPE
	cred.Credentials = map[string]string{
		"First Name":     "Alice",
		"Last Name":      "Student",
		"Student Number": "0123456",
		"Email":          "alice@university.ca",
	}
	return cred, nil
}

type ExamVerifier struct{}

func (ExamVerifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		ServiceURL:  fmt.Sprintf("http://localhost:%d/verify/exam", port),
		EntityName:  "University Exam Attendance Tracker",
		CredType:    CRED_TYPE,
		Description: "Records student information for the purpose of attendance tracking for the exam at 9:00am on December 15th, 2021 in room 101.",
		Issuer:      ISSUER_DID,
		Entity: common.Signature{
			DID: EXAM_VERIFIER_DID,
		},
	}
}

func (ExamVerifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Printf("(Exam Verifier) Verified: %s %s, %s %s", cred.Credentials["First Name"], cred.Credentials["Last Name"], cred.Credentials["Student Number"], cred.Credentials["Email"])
	return nil
}

type EventVerifier struct{}

func (EventVerifier) CreatePresentationRequest() common.PresentationRequest {
	return common.PresentationRequest{
		ServiceURL:  fmt.Sprintf("http://localhost:%d/verify/event", port),
		EntityName:  "University Event Registrar",
		CredType:    CRED_TYPE,
		Description: "Registers the student for the job fair on January 15th, 2022. Upon registration, the student will receive a confirmation by email.",
		Issuer:      ISSUER_DID,
		Entity: common.Signature{
			DID: EVENT_VERIFIER_DID,
		},
	}
}

func (EventVerifier) VerifyCredentials(cred *common.VerifiableCredential) error {
	log.Printf("(Event Verifier) Registered: %s %s, %s %s", cred.Credentials["First Name"], cred.Credentials["Last Name"], cred.Credentials["Student Number"], cred.Credentials["Email"])
	return nil
}

func main() {
	server := demo.DemoServer{
		PublicURL: "./university/public",
		IssuerService: issuer.IssuerService{
			Issuer:        Issuer{},
			DID:           ISSUER_DID,
			PrivateKeyURI: "university/keys/issuer.private.key",
		},
		VerifierServices: map[string]verifier.VerifierService{
			"exam": {
				Verifier:      ExamVerifier{},
				PrivateKeyURI: "university/keys/exam-verifier.private.key",
			},
			"event": {
				Verifier:      EventVerifier{},
				PrivateKeyURI: "university/keys/event-verifier.private.key",
			},
		},
	}
	server.RunServer(port)
}
