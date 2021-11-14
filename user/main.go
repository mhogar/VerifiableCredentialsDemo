package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"vcd/common"
)

func createVerifiableCredential() error {
	//load DID
	DID, err := common.LoadKeyFromFile("DID.cert")
	if err != nil {
		return common.ChainError("error loading DID certificate", err)
	}

	cred := common.VerifiableCredential{
		SubjectDID: string(DID),
		FirstName:  "Bob",
		LastName:   "ADobDob",
	}

	//sign request
	sig, err := common.SignCredential("wallet/private.key", &cred)
	if err != nil {
		return common.ChainError("error signing request", err)
	}
	cred.SubjectSignature = hex.EncodeToString(sig)

	//send request
	buffer, _ := common.EncodeJSON(&cred)
	res, err := http.Post("http://localhost:8082/creds", "application/json", buffer)
	if err != nil {
		return common.ChainError("error sending POST creds request", err)
	}

	defer res.Body.Close()

	//handle error response
	if res.StatusCode != http.StatusOK {
		result := common.ErrorResponse{}
		common.DecodeJSON(res.Body, &result)

		return common.ChainError("error from issuer service", errors.New(result.Error))
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
