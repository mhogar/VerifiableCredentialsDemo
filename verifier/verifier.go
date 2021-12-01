package verifier

import (
	"log"
	"net/http"
	"vcd/common"
)

type Verifier interface {
	CreatePresentationRequest() common.PresentationRequest
	VerifyCredentials(cred *common.VerifiableCredential) error
}

type VerifierService struct {
	Verifier      Verifier
	PrivateKeyURI string
}

func (s VerifierService) GetVerifyHandler(w http.ResponseWriter, _ *http.Request) {
	pres := s.Verifier.CreatePresentationRequest()
	pres.Type = "verify"

	err := common.SignStruct(s.PrivateKeyURI, &pres.Entity, &pres)
	if err != nil {
		common.LogChainError("error signing presentation request", err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, pres)
}

func (s VerifierService) PostVerifyHandler(w http.ResponseWriter, req *http.Request) {
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

	issuerDID, err := common.LoadPublicKeyFromURI(cred.Issuer.DID)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	err = common.VerifyStructSignature(issuerDID, &cred.Issuer.Signature, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying issuer signature")
		return
	}

	err = s.Verifier.VerifyCredentials(&cred)
	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	common.SendSuccessResponse(w)
}
