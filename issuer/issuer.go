package issuer

import (
	"encoding/hex"
	"log"
	"net/http"
	"vcd/common"
)

type IssueRequest struct {
	Fields map[string]string `json:"fields"`

	SubjectDID       string `json:"subject"`
	SubjectSignature string `json:"subject_signature,omitempty"`
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

	//parse request body
	err := common.DecodeJSON(req.Body, &iss)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	//strip signature
	subjectSig := iss.SubjectSignature
	iss.SubjectSignature = ""

	//verify subject signature
	err = common.VerifyStructSignature([]byte(iss.SubjectDID), subjectSig, &iss)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	//create credentials
	cred, err := s.Issuer.CreateVerifiableCredentials(&iss)
	if err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//add DID and create signature
	cred.IssuerDID = s.DID
	sig, err := common.SignStruct(s.PrivateKeyURI, &cred)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	//add signature and send response
	cred.IssuerSig = hex.EncodeToString(sig)
	common.SendJSONResponse(w, http.StatusOK, cred)
}
