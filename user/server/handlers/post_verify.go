package handlers

import (
	"log"
	"net/http"
	"vcd/common"
)

const PRIVATE_KEY_URI = "wallet/private.key"

type PostVerifyBody struct {
	ServiceURL   string `json:"service_url"`
	CredentialID string `json:"credential_id"`
}

func PostVerifyHandler(w http.ResponseWriter, req *http.Request) {
	body := PostVerifyBody{}

	err := common.DecodeJSON(req.Body, &body)
	if err != nil {
		common.LogChainError("error decoding post verify body", err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cerr := postVerify(&body)
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

func postVerify(body *PostVerifyBody) CustomError {
	creds, err := loadVerifiableCredentials()
	if err != nil {
		common.LogChainError("error loading verifiable credentials", err)
		return InternalError()
	}

	cred, ok := (*creds)[body.CredentialID]
	if !ok {
		log.Println("credential with id", body.CredentialID, "no found")
		return ClientError("No credential found for ID.")
	}

	err = common.SignStruct(PRIVATE_KEY_URI, &cred.Subject, &cred)
	if err != nil {
		common.LogChainError("error signing credential", err)
		return InternalError()
	}

	_, cerr, err := sendRequest(http.MethodPost, body.ServiceURL, &cred)
	if err != nil {
		log.Println(err)
	}
	if cerr.Type != TypeNoError {
		return cerr
	}

	return NoError()
}
