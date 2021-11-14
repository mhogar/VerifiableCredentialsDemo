package main

import (
	"encoding/hex"
	"log"
	"net/http"
	"vcd/helpers"
)

type RequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CredentialResponse struct {
	RequestBody
	IssuerDID string `json:"issuer"`
	IssuerSig string `json:"issuer_sig,omitempty"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	body := RequestBody{}

	//parse request body
	err := helpers.DecodeJSON(req.Body, &body)
	if err != nil {
		log.Println(err)
		helpers.SendBadRequestResponse(w, "invalid JSON body")
		return
	}

	//TODO: verify request fields

	//create response body
	res := CredentialResponse{
		RequestBody: body,
		IssuerDID:   "Issuer DID",
	}

	//sign response
	sig, err := helpers.SignResponse("keys/private.key", res)
	if err != nil {
		log.Println(err)
		helpers.SendInternalErrorResponse(w)
		return
	}

	//send response
	res.IssuerSig = hex.EncodeToString(sig)
	helpers.SendJSONResponse(w, http.StatusOK, res)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
