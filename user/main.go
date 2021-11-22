package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"vcd/common"
)

func sendRequest(cred *common.VerifiableCredential, url string) (io.ReadCloser, error) {
	//sign credential
	sig, err := common.SignCredential("wallet/private.key", cred)
	if err != nil {
		return nil, common.ChainError("error signing credential", err)
	}
	cred.SubjectSignature = hex.EncodeToString(sig)

	//send request
	buffer, _ := common.EncodeJSON(&cred)
	res, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return nil, common.ChainError("error sending POST creds request", err)
	}

	//handle error response
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		result := common.ErrorResponse{}
		common.DecodeJSON(res.Body, &result)

		return nil, common.ChainError("error from server", errors.New(result.Error))
	}

	return res.Body, nil
}

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

	//send issue vc request
	body, err := sendRequest(&cred, "http://localhost:8082/creds")
	if err != nil {
		return err
	}
	defer body.Close()

	//-- save credential --
	fmt.Println("Saving verifiable credential...")

	f, _ := os.Create("wallet/vc.json")
	defer f.Close()

	f.ReadFrom(body)
	return nil
}

func handleVerifyCredential(w http.ResponseWriter, _ *http.Request) {
	//open file
	f, err := os.Open("wallet/vc.json")
	if err != nil {
		log.Println(common.ChainError("error opening vc file", err))
		common.SendInternalErrorResponse(w)
		return
	}
	defer f.Close()

	//decode json
	cred := common.VerifiableCredential{}
	err = common.DecodeJSON(f, &cred)
	if err != nil {
		log.Println(common.ChainError("error decoding JSON", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//send verify request
	_, err = sendRequest(&cred, "http://localhost:8083/verify")
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendSuccessResponse(w)
}

func main() {
	//parse flags
	port := flag.Int("port", 8080, "port to run the server on")
	flag.Parse()

	//setup routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/verify", handleVerifyCredential)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
