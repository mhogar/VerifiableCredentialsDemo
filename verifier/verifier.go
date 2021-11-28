package verifier

import (
	"log"
	"net/http"
	"vcd/common"
)

type PresentationRequest struct {
	Purpose  string           `json:"purpose"`
	Issuer   string           `json:"issuer"`
	Verifier common.Signature `json:"verifier"`
}

var DIDLoader = common.DIDFileLoader{}

type Verifier interface {
	CreatePresentationRequest() PresentationRequest
	VerifyCredentials(cred *common.VerifiableCredential) error
}

type VerifierService struct {
	Verifier      Verifier
	PrivateKeyURI string
}

func (s VerifierService) GetVerifyHandler(w http.ResponseWriter, _ *http.Request) {
	pres := s.Verifier.CreatePresentationRequest()

	err := common.SignStruct(s.PrivateKeyURI, &pres.Verifier, &pres)
	if err != nil {
		log.Println(err)
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

	err = common.VerifyStructSignature([]byte(cred.Subject.DID), &cred.Subject, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	issuerDID, err := DIDLoader.LoadPublicKeyFromURI(cred.Issuer.DID)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	err = common.VerifyStructSignature(issuerDID, &cred.Issuer, &cred)
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
