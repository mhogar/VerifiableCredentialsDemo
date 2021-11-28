package demo

import (
	"fmt"
	"net/http"
	"vcd/common"
	"vcd/issuer"
	"vcd/verifier"
)

type DemoServer struct {
	PublicURL       string
	VerifierService verifier.VerifierService
	IssuerService   issuer.IssuerService
}

func (s DemoServer) verifyHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.VerifierService.GetVerifyHandler(w, req)
	case http.MethodPost:
		s.VerifierService.PostVerifyHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func (s DemoServer) issueHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		s.IssuerService.PostIssueHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func (s DemoServer) RunServer(port int) {
	//setup routes
	http.Handle("/", http.FileServer(http.Dir(s.PublicURL)))
	http.HandleFunc("/verify", s.verifyHandler)
	http.HandleFunc("/issue", s.issueHandler)

	//run the server
	fmt.Printf("listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
