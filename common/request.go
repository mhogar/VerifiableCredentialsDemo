package common

type ErrorResponse struct {
	Error string `json:"error"`
}

type VerifiableCredential struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	SubjectDID       string `json:"subject"`
	SubjectSignature string `json:"subject_signature,omitempty"`

	IssuerDID string `json:"issuer,omitempty"`
	IssuerSig string `json:"issuer_signature,omitempty"`
}
