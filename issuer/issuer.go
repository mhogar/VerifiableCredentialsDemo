package issuer

import (
	"log"
	"net/http"
	"vcd/common"
)

type IssueRequest struct {
	Fields map[string]string `json:"fields"`

	Subject common.Signature `json:"subject"`
}

type Issuer interface {
	CreateVerifiableCredentials(iss *IssueRequest) (*common.VerifiableCredential, error)
}

type IssuerService struct {
	Issuer        Issuer
	DID           string
	PrivateKeyURI string
}

func (s IssuerService) PostIssueHandler(w http.ResponseWriter, req *http.Request) {
	iss := IssueRequest{}

	err := common.DecodeJSON(req.Body, &iss)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err = common.VerifyStructSignature([]byte(iss.Subject.DID), &iss.Subject, &iss)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	cred, err := s.Issuer.CreateVerifiableCredentials(&iss)
	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cred.Issuer = common.Signature{
		DID: s.DID,
	}

	err = common.SignStruct(s.PrivateKeyURI, &cred.Issuer, &cred)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, cred)
}
