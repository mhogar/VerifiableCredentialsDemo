package main

import (
	"encoding/hex"
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

	//strip signature from credential
	subjectSig := cred.SubjectSignature
	cred.SubjectSignature = ""

	//verify subject signature
	err = common.VerifyCredentialSignature(cred.SubjectDID, subjectSig, &cred)
	if err != nil {
		log.Println(err)
		common.SendErrorResponse(w, http.StatusUnauthorized, "error verifying subject signature")
		return
	}

	//TODO: verify request fields?

	//add DID and create signature
	cred.IssuerDID = "blockchain/issuer.json"
	sig, err := common.SignCredential("keys/private.key", cred)
	if err != nil {
		log.Println(err)
		common.SendInternalErrorResponse(w)
		return
	}

	//add signature and send response
	cred.IssuerSig = hex.EncodeToString(sig)
	common.SendJSONResponse(w, http.StatusOK, cred)
}

func main() {
	//parse flags
	port := flag.Int("port", 8082, "port to run the server on")
	flag.Parse()

	//setup routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/creds", handler)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
