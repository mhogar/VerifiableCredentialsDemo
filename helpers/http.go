package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

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
	SendJSONResponse(w, status, ErrorResponse{
		Error: err,
	})
}

func SendBadRequestResponse(w http.ResponseWriter, err string) {
	SendErrorResponse(w, http.StatusBadRequest, err)
}

func SendInternalErrorResponse(w http.ResponseWriter) {
	SendErrorResponse(w, http.StatusInternalServerError, "an internal error occurred")
}
