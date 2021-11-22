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

type PresentationRequestResponse struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Purpose string `json:"purpose"`
}

const VERIFIER_URL = "http://localhost:8083/verify"

var DIDLoader = common.DIDFileLoader{}

func sendRequest(method string, url string, body *common.VerifiableCredential) (io.ReadCloser, error) {
	var buffer io.Reader = nil

	if body != nil {
		//sign body
		sig, err := common.SignStruct("wallet/private.key", body)
		if err != nil {
			return nil, common.ChainError("error signing credential", err)
		}
		body.SubjectSignature = hex.EncodeToString(sig)

		buffer, _ = common.EncodeJSON(body)
	}

	//send request
	req, _ := http.NewRequest(method, url, buffer)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, common.ChainError("error sending request", err)
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
	body, err := sendRequest("http://localhost:8082/creds", http.MethodPost, &cred)
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

func getVerify() (*PresentationRequestResponse, error) {
	body, err := sendRequest(http.MethodGet, VERIFIER_URL, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	pres := common.PresentationRequest{}
	err = common.DecodeJSON(body, &pres)
	if err != nil {
		return nil, err
	}

	doc, err := DIDLoader.LoadDIDDocumentFromURI(pres.VerifierDID)
	if err != nil {
		return nil, err
	}

	DID, err := DIDLoader.LoadPublicKeyFromDocument(doc)
	if err != nil {
		return nil, err
	}

	err = common.VerifyStructSignature(DID, pres.VerifierSignature, &pres)
	if err != nil {
		//TODO: fix verify issue
		//return nil, common.ChainError("error verifying verifier signature", err)
	}

	//TODO: receive issuer DID to determine what type of creds to accept
	//TODO: validate issuer signature on verifier DID document

	return &PresentationRequestResponse{
		Name:    pres.Name,
		Purpose: pres.Purpose,
		Domain:  doc.Domain,
	}, nil
}

func getVerifyHandler(w http.ResponseWriter) {
	pres, err := getVerify()
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendJSONResponse(w, http.StatusOK, pres)
}

func postVerifyHandler(w http.ResponseWriter) {
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
	_, err = sendRequest(http.MethodPost, VERIFIER_URL, &cred)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendSuccessResponse(w)
}

func verifyHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getVerifyHandler(w)
	case http.MethodPost:
		postVerifyHandler(w)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func main() {
	// createVerifiableCredential()
	// return

	//parse flags
	port := flag.Int("port", 8080, "port to run the server on")
	flag.Parse()

	//setup routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/verify", verifyHandler)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
