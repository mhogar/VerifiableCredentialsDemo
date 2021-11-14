package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"vcd/common"
	"vcd/helpers"
)

func createVerifiableCredential() (*common.VerifiableCredential, error) {
	//load DID
	DID, err := helpers.LoadKeyFromFile("keys/public.cert")
	if err != nil {
		log.Fatal(err)
	}

	cred := common.VerifiableCredential{
		SubjectDID: string(DID),
		FirstName:  "Bob",
		LastName:   "ADobDob",
	}

	//sign request
	sig, err := helpers.SignCredential("keys/private.key", &cred)
	if err != nil {
		return nil, helpers.ChainError("error signing request", err)
	}
	cred.SubjectSignature = hex.EncodeToString(sig)

	//send request
	buffer, _ := helpers.EncodeJSON(&cred)
	res, err := http.Post("http://localhost:8082/creds", "application/json", buffer)
	if err != nil {
		return nil, helpers.ChainError("error sending POST creds request", err)
	}

	defer res.Body.Close()

	//handle error response
	if res.StatusCode != http.StatusOK {
		result := common.ErrorResponse{}
		helpers.DecodeJSON(res.Body, &result)

		return nil, helpers.ChainError("error from issuer service", errors.New(result.Error))
	}

	//parse credential
	cred = common.VerifiableCredential{}
	helpers.DecodeJSON(res.Body, &cred)

	return &cred, nil
}

func main() {
	cred, err := createVerifiableCredential()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*cred)
}
