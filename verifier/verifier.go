package verifier

import (
	"encoding/base64"
	"log"
	"net/http"
	"vcd/common"
)

type PresentationRequest struct {
	Name    string `json:"string"`
	Purpose string `json:"purpose"`

	VerifierDID       string `json:"verifier"`
	VerifierSignature string `json:"verifier_signature,omitempty"`
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

	sig, err := common.SignStruct(s.PrivateKeyURI, &pres)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	pres.VerifierSignature = base64.RawStdEncoding.EncodeToString(sig)
	common.SendJSONResponse(w, http.StatusOK, pres)
}

func (s VerifierService) PostVerifyHandler(w http.ResponseWriter, req *http.Request) {
	cred := common.VerifiableCredential{}

	//parse request body
	err := common.DecodeJSON(req.Body, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	//strip subject signature from credential
	sig := cred.SubjectSignature
	cred.SubjectSignature = ""

	//verify subject signature
	err = common.VerifyStructSignature([]byte(cred.SubjectDID), sig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	//load issuer DID
	issuerDID, err := DIDLoader.LoadPublicKeyFromURI(cred.IssuerDID)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	//strip issuer signature and DID from credential
	sig = cred.IssuerSig
	cred.IssuerSig = ""

	//verify issuer signature
	err = common.VerifyStructSignature(issuerDID, sig, &cred)
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
