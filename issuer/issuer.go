package issuer

import (
	"log"
	"net/http"
	"vcd/common"
)

type Issuer interface {
	CreatePresentationRequest() common.PresentationRequest
	CreateVerifiableCredentials(creds *common.VerifiableCredential) (*common.VerifiableCredential, error)
}

type IssuerService struct {
	Issuer        Issuer
	DID           string
	PrivateKeyURI string
}

func (s IssuerService) GetIssueHandler(w http.ResponseWriter, _ *http.Request) {
	pres := s.Issuer.CreatePresentationRequest()

	err := common.SignStruct(s.PrivateKeyURI, &pres.Entity, &pres)
	if err != nil {
		common.LogChainError("error signing presentation request", err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, &pres)
}

func (s IssuerService) PostIssueHandler(w http.ResponseWriter, req *http.Request) {
	cred := &common.VerifiableCredential{}

	err := common.DecodeJSON(req.Body, cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err = common.VerifyStructSignature([]byte(cred.Subject.DID), &cred.Subject.Signature, cred)
	if err != nil {
		common.LogChainError("error verifying subject signature", err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "Subject signature could not be verified.")
		return
	}

	if cred.Issuer.DID != "" {
		bytes, err := common.LoadPublicKeyFromURI(cred.Issuer.DID)
		if err != nil {
			common.LogChainError("error loading issuer public key", err)
			common.SendInternalErrorResponse(w)
			return
		}

		err = common.VerifyStructSignature(bytes, &cred.Issuer.Signature, cred)
		if err != nil {
			common.LogChainError("error verifying issuer signature", err)
			common.SendErrorResponse(w, http.StatusUnauthorized, "Issuer signature could not be verified.")
			return
		}
	}

	cred, err = s.Issuer.CreateVerifiableCredentials(cred)
	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cred.Issuer = common.Signature{
		DID: s.DID,
	}

	err = common.SignStruct(s.PrivateKeyURI, &cred.Issuer, cred)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, cred)
}
