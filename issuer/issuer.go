package issuer

import (
	"log"
	"net/http"
	"vcd/common"
)

type IssueRequest struct {
	Purpose string            `json:"purpose"`
	Fields  map[string]string `json:"fields"`
	Issuer  common.Signature  `json:"issuer"`
	Subject common.Signature  `json:"subject"`
}

type Issuer interface {
	CreateIssueRequest() IssueRequest
	CreateVerifiableCredentials(creds *common.VerifiableCredential) error
}

type IssuerService struct {
	Issuer        Issuer
	DID           string
	PrivateKeyURI string
}

func (s IssuerService) GetIssueHandler(w http.ResponseWriter, _ *http.Request) {
	iss := s.Issuer.CreateIssueRequest()

	err := common.SignStruct(s.PrivateKeyURI, &iss.Issuer, &iss)
	if err != nil {
		common.LogChainError("error signing issuer request", err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, iss)
}

func (s IssuerService) PostIssueHandler(w http.ResponseWriter, req *http.Request) {
	cred := common.VerifiableCredential{}

	err := common.DecodeJSON(req.Body, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err = common.VerifyStructSignature([]byte(cred.Subject.DID), &cred.Subject.Signature, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	err = s.Issuer.CreateVerifiableCredentials(&cred)
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
