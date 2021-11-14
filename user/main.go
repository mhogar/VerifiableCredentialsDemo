package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"vcd/common"
	"vcd/helpers"
)

func createVerifiableCredential() error {
	//load DID
	DID, err := helpers.LoadKeyFromFile("DID.cert")
	if err != nil {
		return helpers.ChainError("error loading DID certificate", err)
	}

	cred := common.VerifiableCredential{
		SubjectDID: string(DID),
		FirstName:  "Bob",
		LastName:   "ADobDob",
	}

	//sign request
	sig, err := helpers.SignCredential("wallet/private.key", &cred)
	if err != nil {
		return helpers.ChainError("error signing request", err)
	}
	cred.SubjectSignature = hex.EncodeToString(sig)

	//send request
	buffer, _ := helpers.EncodeJSON(&cred)
	res, err := http.Post("http://localhost:8082/creds", "application/json", buffer)
	if err != nil {
		return helpers.ChainError("error sending POST creds request", err)
	}

	defer res.Body.Close()

	//handle error response
	if res.StatusCode != http.StatusOK {
		result := common.ErrorResponse{}
		helpers.DecodeJSON(res.Body, &result)

		return helpers.ChainError("error from issuer service", errors.New(result.Error))
	}

	//-- save credential --
	fmt.Println("Saving verifiable credential...")

	f, _ := os.Create("wallet/vc.json")
	defer f.Close()

	f.ReadFrom(res.Body)

	return nil
}

func main() {
	err := createVerifiableCredential()
	if err != nil {
		log.Fatal(err)
	}
}
