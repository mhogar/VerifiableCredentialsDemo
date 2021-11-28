package common

import "errors"

type VerifiableCredential struct {
	Credentials map[string]string `json:"credentials"`

	SubjectDID       string `json:"subject"`
	SubjectSignature string `json:"subject_signature,omitempty"`

	IssuerDID string `json:"issuer,omitempty"`
	IssuerSig string `json:"issuer_signature,omitempty"`
}

func ChainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}
