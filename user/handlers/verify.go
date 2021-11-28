package handlers

import (
	"log"
	"net/http"
	"os"
	"vcd/common"
	"vcd/verifier"
)

const PRIVATE_KEY_URI = "wallet/private.key"

const VERIFIER_URL = "http://localhost:8082/verify"
const ISSUER_URL = "http://localhost:8082/issue"

type PresentationRequestResponse struct {
	Name            string `json:"name"`
	Domain          string `json:"domain"`
	Purpose         string `json:"purpose"`
	TrustedByIssuer bool   `json:"trusted_by_issuer"`
}

func VerifyHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getVerifyHandler(w, req)
	case http.MethodPost:
		postVerifyHandler(w)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func getVerifyHandler(w http.ResponseWriter, req *http.Request) {
	pres, cerr := getVerify(req.URL.Query().Get("url"))
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusOK, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, pres)
}

func getVerify(url string) (*PresentationRequestResponse, CustomError) {
	body, err := sendRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogChainError("error sending get verify request", err)
		return nil, ClientError("Invalid URL.")
	}
	defer body.Close()

	pres := verifier.PresentationRequest{}
	err = common.DecodeJSON(body, &pres)
	if err != nil {
		common.LogChainError("error decoding get verify response", err)
		return nil, ClientError("Invalid verify endpoint.")
	}

	doc, err := common.LoadDIDDocumentFromURI(pres.Verifier.DID)
	if err != nil {
		common.LogChainError("error loading DID doc", err)
		return nil, InternalError()
	}

	DID, err := common.LoadPublicKeyFromDocument(doc)
	if err != nil {
		common.LogChainError("error loading public key from DID doc", err)
		return nil, InternalError()
	}

	err = common.VerifyStructSignature(DID, &pres.Verifier.Signature, &pres)
	if err != nil {
		common.LogChainError("error verifying verifier signature", err)
		return nil, ClientError("Verifier has invalid signature.")
	}

	res := PresentationRequestResponse{
		Name:            doc.Name,
		Purpose:         pres.Purpose,
		Domain:          doc.Domain,
		TrustedByIssuer: false,
	}

	err = common.VerifyDIDDocumentSignature(doc, pres.Issuer)
	if err == nil {
		res.TrustedByIssuer = true
	}

	return &res, NoError()
}

func postVerifyHandler(w http.ResponseWriter) {
	cerr := postVerify()
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusOK, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendSuccessResponse(w)
}

func postVerify() CustomError {
	f, err := os.Open("wallet/vc.json")
	if err != nil {
		common.LogChainError("error opening vc file", err)
		return InternalError()
	}
	defer f.Close()

	cred := common.VerifiableCredential{}
	err = common.DecodeJSON(f, &cred)
	if err != nil {
		common.LogChainError("error decoding credential JSON", err)
		return InternalError()
	}

	err = common.SignStruct(PRIVATE_KEY_URI, &cred.Subject, &cred)
	if err != nil {
		common.LogChainError("error signing credential", err)
		return InternalError()
	}

	_, err = sendRequest(http.MethodPost, VERIFIER_URL, &cred)
	if err != nil {
		log.Println(err)
		return InternalError()
	}

	return NoError()
}
