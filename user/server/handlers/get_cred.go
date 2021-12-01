package handlers

import (
	"net/http"
	"vcd/common"
)

func GetCredHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	creds := CredentialsMap{}
	err := common.LoadJSONFromFile(VC_URI, &creds)
	if err != nil {
		common.LogChainError("error loading verifiable credentials", err)
		common.SendInternalErrorResponse(w)
		return
	}

	cred, ok := creds[id]
	if !ok {
		common.SendSuccessResponse(w)
	}

	common.SendJSONResponse(w, http.StatusOK, &cred)
}
