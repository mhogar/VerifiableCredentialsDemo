package handlers

import (
	"log"
	"net/http"
	"vcd/common"
)

const PRIVATE_KEY_URI = "wallet/private.key"

type PresentationRequestResponse struct {
	Name            string `json:"name"`
	Domain          string `json:"domain"`
	Purpose         string `json:"purpose"`
	Issuer          string `json:"issuer"`
	TrustedByIssuer bool   `json:"trusted_by_issuer"`
}

func VerifyHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	// case http.MethodGet:
	// 	getVerifyHandler(w, req)
	case http.MethodPost:
		postVerifyHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

// func getVerifyHandler(w http.ResponseWriter, req *http.Request) {
// 	pres, cerr := getVerify(req.URL.Query().Get("url"))
// 	if cerr.Type == TypeClientError {
// 		common.SendErrorResponse(w, http.StatusBadRequest, cerr.Message)
// 		return
// 	}
// 	if cerr.Type == TypeInternalError {
// 		common.SendInternalErrorResponse(w)
// 		return
// 	}

// 	common.SendJSONResponse(w, http.StatusOK, pres)
// }

// func getVerify(url string) (*PresentationRequestResponse, CustomError) {
// 	body, err := sendRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		common.LogChainError("error sending get verify request", err)
// 		return nil, ClientError("Invalid URL.")
// 	}
// 	defer body.Close()

// 	pres := verifier.PresentationRequest{}
// 	err = common.DecodeJSON(body, &pres)
// 	if err != nil {
// 		common.LogChainError("error decoding get verify response", err)
// 		return nil, ClientError("Invalid verify endpoint.")
// 	}

// 	doc, err := common.LoadDIDDocumentFromURI(pres.Verifier.DID)
// 	if err != nil {
// 		common.LogChainError("error loading DID doc", err)
// 		return nil, InternalError()
// 	}

// 	DID, err := common.LoadPublicKeyFromDocument(doc)
// 	if err != nil {
// 		common.LogChainError("error loading public key from DID doc", err)
// 		return nil, InternalError()
// 	}

// 	err = common.VerifyStructSignature(DID, &pres.Verifier.Signature, &pres)
// 	if err != nil {
// 		common.LogChainError("error verifying verifier signature", err)
// 		return nil, ClientError("Verifier has invalid signature.")
// 	}

// 	res := PresentationRequestResponse{
// 		Name:            doc.Name,
// 		Purpose:         pres.Purpose,
// 		Domain:          doc.Domain,
// 		Issuer:          pres.Issuer,
// 		TrustedByIssuer: false,
// 	}

// 	err = common.VerifyDIDDocumentSignature(doc, pres.Issuer)
// 	if err == nil {
// 		res.TrustedByIssuer = true
// 	}

// 	return &res, NoError()
// }

func postVerifyHandler(w http.ResponseWriter, req *http.Request) {
	body := map[string]string{}

	err := common.DecodeJSON(req.Body, &body)
	if err != nil {
		common.LogChainError("error decoding post verify body", err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	url, ok := body["url"]
	if !ok {
		common.SendErrorResponse(w, http.StatusBadRequest, "missing required \"url\" parameter")
		return
	}

	cerr := postVerify(url, body["issuer"])
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

func postVerify(url string, issuer string) CustomError {
	creds := map[string]common.VerifiableCredential{}
	err := common.LoadJSONFromFile(VC_URI, &creds)
	if err != nil {
		common.LogChainError("error loading verifiable credentials", err)
		return InternalError()
	}

	cred, ok := creds[issuer]
	if !ok {
		return ClientError("No credentials for issuer.")
	}

	err = common.SignStruct(PRIVATE_KEY_URI, &cred.Subject, &cred)
	if err != nil {
		common.LogChainError("error signing credential", err)
		return InternalError()
	}

	_, err = sendRequest(http.MethodPost, url, &cred)
	if err != nil {
		log.Println(err)
		return InternalError()
	}

	return NoError()
}
