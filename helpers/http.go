package helpers

import (
	"encoding/json"
	"net/http"
	"vcd/common"
)

func SendJSONResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	if res != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}
}

func SendErrorResponse(w http.ResponseWriter, status int, err string) {
	SendJSONResponse(w, status, common.ErrorResponse{
		Error: err,
	})
}

func SendInternalErrorResponse(w http.ResponseWriter) {
	SendErrorResponse(w, http.StatusInternalServerError, "an internal error occurred")
}
