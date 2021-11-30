package handlers

import (
	"net/http"
	"vcd/common"
)

const VC_URI = "wallet/verifiable-credentials.json"

func GetCredsHandler(w http.ResponseWriter, _ *http.Request) {
	creds := CredentialsMap{}
	err := common.LoadJSONFromFile(VC_URI, &creds)
	if err != nil {
		common.LogChainError("error loading verifiable credentials", err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, &creds)
}
