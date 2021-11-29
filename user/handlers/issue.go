package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"vcd/common"
	"vcd/issuer"
)

type IssueRequestResponse struct {
	Name    string            `json:"name"`
	Domain  string            `json:"domain"`
	Purpose string            `json:"purpose"`
	Fields  map[string]string `json:"fields"`
}

type IssueRequest struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

func IssueHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getIssueHandler(w, req)
	case http.MethodPost:
		postIssueHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func getIssueHandler(w http.ResponseWriter, req *http.Request) {
	pres, cerr := getIssue(req.URL.Query().Get("url"))
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusBadRequest, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, pres)
}

func getIssue(url string) (*IssueRequestResponse, CustomError) {
	body, err := sendRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogChainError("error sending get issue request", err)
		return nil, ClientError("Invalid URL.")
	}
	defer body.Close()

	iss := issuer.IssueRequest{}
	err = common.DecodeJSON(body, &iss)
	if err != nil {
		common.LogChainError("error decoding get issue response", err)
		return nil, ClientError("Invalid verify endpoint.")
	}

	doc, err := common.LoadDIDDocumentFromURI(iss.Issuer.DID)
	if err != nil {
		common.LogChainError("error loading DID doc", err)
		return nil, InternalError()
	}

	DID, err := common.LoadPublicKeyFromDocument(doc)
	if err != nil {
		common.LogChainError("error loading public key from DID doc", err)
		return nil, InternalError()
	}

	err = common.VerifyStructSignature(DID, &iss.Issuer.Signature, &iss)
	if err != nil {
		common.LogChainError("error verifying issuer signature", err)
		return nil, ClientError("Issuer has invalid signature.")
	}

	res := IssueRequestResponse{
		Name:    doc.Name,
		Purpose: iss.Purpose,
		Domain:  doc.Domain,
		Fields:  iss.Fields,
	}

	return &res, NoError()
}

func postIssueHandler(w http.ResponseWriter, req *http.Request) {
	iss := IssueRequest{}

	err := common.DecodeJSON(req.Body, &iss)
	if err != nil {
		common.LogChainError("error decoding post issue body", err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cerr := postIssue(&iss)
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusBadRequest, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendSuccessResponse(w)
}

func postIssue(iss *IssueRequest) CustomError {
	cred := common.VerifiableCredential{}

	bytes, err := os.ReadFile("DID.cert")
	if err != nil {
		common.LogChainError("error reading DID file", err)
		return InternalError()
	}
	cred.Subject.DID = string(bytes)

	cred.Credentials = iss.Fields

	err = common.SignStruct(PRIVATE_KEY_URI, &cred.Subject, cred)
	if err != nil {
		common.LogChainError("error signing issue request", err)
		return InternalError()
	}

	body, err := sendRequest(http.MethodPost, iss.URL, &cred)
	if err != nil {
		log.Println(err)
		return InternalError()
	}
	defer body.Close()

	cred = common.VerifiableCredential{}
	err = common.DecodeJSON(body, &cred)
	if err != nil {
		common.LogChainError("error decoding verifiable credential", err)
		return InternalError()
	}

	err = common.WriteJSONToFile(fmt.Sprintf("wallet/%s.json", cred.Issuer.DID), &cred)
	if err != nil {
		common.LogChainError("error saving verifiable credential", err)
		return InternalError()
	}

	return NoError()
}
