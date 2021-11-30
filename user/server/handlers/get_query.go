package handlers

import (
	"net/http"
	"vcd/common"
)

type QueryResponse struct {
	Type   string            `json:"type"`
	Fields map[string]string `json:"fields,omitempty"`

	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Purpose string `json:"purpose"`

	Issuer          string `json:"issuer,omitempty"`
	TrustedByIssuer bool   `json:"trusted_by_issuer"`
}

func GetQueryHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")
	if url == "" {
		common.SendErrorResponse(w, http.StatusBadGateway, "missing required parameter 'url'")
		return
	}

	res, cerr := getQuery(url)
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusBadRequest, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, res)
}

func getQuery(url string) (*QueryResponse, CustomError) {
	body, err := sendRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogChainError("error sending query request", err)
		return nil, ClientError("Invalid URL.")
	}
	defer body.Close()

	pres := common.PresentationRequest{}
	err = common.DecodeJSON(body, &pres)
	if err != nil {
		common.LogChainError("error decoding query response", err)
		return nil, ClientError("Invalid URL.")
	}

	doc, err := common.LoadDIDDocumentFromURI(pres.Entity.DID)
	if err != nil {
		common.LogChainError("error loading DID doc", err)
		return nil, InternalError()
	}

	DID, err := common.LoadPublicKeyFromDocument(doc)
	if err != nil {
		common.LogChainError("error loading public key from DID doc", err)
		return nil, InternalError()
	}

	err = common.VerifyStructSignature(DID, &pres.Entity.Signature, &pres)
	if err != nil {
		common.LogChainError("error verifying entity signature", err)
		return nil, ClientError("Entity cannot be verified.")
	}

	res := QueryResponse{
		Type:    pres.Type,
		Fields:  pres.Fields,
		Name:    doc.Name,
		Domain:  doc.Domain,
		Purpose: pres.Purpose,
	}

	if pres.Issuer != "" {
		res.TrustedByIssuer = (common.VerifyDIDDocumentSignature(doc, pres.Issuer) == nil)
	}

	return &res, NoError()
}
