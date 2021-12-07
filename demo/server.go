package demo

import (
	"fmt"
	"net/http"
	"vcd/common"
	"vcd/issuer"
	"vcd/verifier"
)

type DemoServer struct {
	PublicURL        string
	IssuerService    issuer.IssuerService
	VerifierServices map[string]verifier.VerifierService
}

func (s DemoServer) issueHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.IssuerService.GetIssueHandler(w, req)
	case http.MethodPost:
		s.IssuerService.PostIssueHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func (DemoServer) createVerifyHandler(v verifier.VerifierService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			v.GetVerifyHandler(w, req)
		case http.MethodPost:
			v.PostVerifyHandler(w, req)
		default:
			common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
		}
	}
}

func (s DemoServer) RunServer(port int) {
	http.Handle("/", http.FileServer(http.Dir(s.PublicURL)))

	http.HandleFunc("/issue", s.issueHandler)
	fmt.Printf("- http://localhost:%d/issue\n", port)

	for key, val := range s.VerifierServices {
		http.HandleFunc("/verify/"+key, s.createVerifyHandler(val))
		fmt.Printf("- http://localhost:%d/verify/%s\n", port, key)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
