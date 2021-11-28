package common

import (
	"errors"
	"log"
)

type Signature struct {
	DID       string `json:"did"`
	Signature string `json:"signature,omitempty"`
}

type VerifiableCredential struct {
	Credentials map[string]string `json:"credentials"`

	Subject Signature `json:"subject"`
	Issuer  Signature `json:"issuer"`
}

func ChainError(message string, err error) error {
	return errors.New(message + "\n\t" + err.Error())
}

func LogChainError(message string, err error) {
	log.Println(ChainError(message, err))
}
