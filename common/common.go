package common

import (
	"errors"
	"log"
)

type Signature struct {
	DID       string `json:"did"`
	Signature string `json:"signature,omitempty"`
}

type PresentationField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type PresentationRequest struct {
	Type       string `json:"type"`
	ServiceURL string `json:"service_url"`

	EntityName  string              `json:"entity_name"`
	CredType    string              `json:"cred_type"`
	Description string              `json:"description"`
	Fields      []PresentationField `json:"fields,omitempty"`

	Issuer string    `json:"issuer,omitempty"`
	Entity Signature `json:"entity"`
}

type VerifiableCredential struct {
	CredType    string            `json:"cred_type"`
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
