package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"vcd/common"
)

const DID = "verifier.json"

var DIDLoader = common.DIDFileLoader{}

func getVerifyHandler(w http.ResponseWriter, _ *http.Request) {
	pres := common.PresentationRequest{
		Name:        "Sample Verifier",
		Purpose:     "Logs first and last name.",
		VerifierDID: DID,
	}

	sig, err := common.SignStruct("keys/private.key", &pres)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	pres.VerifierSignature = hex.EncodeToString(sig)
	common.SendJSONResponse(w, http.StatusOK, pres)
}

func postVerifyHandler(w http.ResponseWriter, req *http.Request) {
	cred := common.VerifiableCredential{}

	//parse request body
	err := common.DecodeJSON(req.Body, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	//strip subject signature from credential
	sig := cred.SubjectSignature
	cred.SubjectSignature = ""

	//verify subject signature
	err = common.VerifyStructSignature([]byte(cred.SubjectDID), sig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	//load issuer DID
	issuerDID, err := DIDLoader.LoadPublicKeyFromURI(cred.IssuerDID)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	//strip issuer signature and DID from credential
	sig = cred.IssuerSig
	cred.IssuerSig = ""

	//verify issuer signature
	err = common.VerifyStructSignature(issuerDID, sig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying issuer signature")
		return
	}

	//TODO: do something with verified credential
	log.Println("Verified:", cred.FirstName, cred.LastName)

	common.SendSuccessResponse(w)
}

func handler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getVerifyHandler(w, req)
	case http.MethodPost:
		postVerifyHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func main() {
	//parse flags
	port := flag.Int("port", 8083, "port to run the server on")
	flag.Parse()

	//setup routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/verify", handler)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
