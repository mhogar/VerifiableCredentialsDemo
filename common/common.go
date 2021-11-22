package common

import "errors"

type VerifiableCredential struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	SubjectDID       string `json:"subject"`
	SubjectSignature string `json:"subject_signature,omitempty"`

	IssuerDID string `json:"issuer,omitempty"`
	IssuerSig string `json:"issuer_signature,omitempty"`
}

type PresentationRequest struct {
	Name    string `json:"string"`
	Purpose string `json:"purpose"`

	VerifierDID       string `json:"verifier"`
	VerifierSignature string `json:"verifier_signature,omitempty"`
}

func ChainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}
