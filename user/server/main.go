package main

import (
	"flag"
	"fmt"
	"net/http"
	"vcd/common"
	"vcd/user/server/handlers"
)

func createHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length")

		switch req.Method {
		case http.MethodOptions:
			return
		case method:
			handler(w, req)
		default:
			common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
		}
	}
}

func main() {
	//parse flags
	port := flag.Int("port", 8082, "port to run the server on")
	flag.Parse()

	//setup routes
	http.HandleFunc("/creds", createHandler(http.MethodGet, handlers.GetCredsHandler))
	http.HandleFunc("/cred", createHandler(http.MethodGet, handlers.GetCredHandler))
	http.HandleFunc("/query", createHandler(http.MethodGet, handlers.GetQueryHandler))
	http.HandleFunc("/verify", createHandler(http.MethodPost, handlers.PostVerifyHandler))
	http.HandleFunc("/issue", createHandler(http.MethodPost, handlers.PostIssueHandler))

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
