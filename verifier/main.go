package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"vcd/common"
)

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
		return
	}
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
	err = common.VerifyCredentialSignature([]byte(cred.SubjectDID), sig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	//load issuer DID
	issuerDID, err := common.LoadDIDFromBlockchain(cred.IssuerDID)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	//strip issuer signature and DID from credential
	sig = cred.IssuerSig
	cred.IssuerSig = ""

	//verify issuer signature
	err = common.VerifyCredentialSignature(issuerDID, sig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying issuer signature")
		return
	}

	//TODO: do something with verified credential

	common.SendSuccessResponse(w)
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
