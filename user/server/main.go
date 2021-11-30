package main

import (
	"flag"
	"fmt"
	"net/http"
	"vcd/common"
	"vcd/user/server/handlers"
)

func credsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

	switch req.Method {
	case http.MethodGet:
		handlers.GetCredsHandler(w, req)
	default:
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid request method")
	}
}

func main() {
	//parse flags
	port := flag.Int("port", 8082, "port to run the server on")
	flag.Parse()

	//setup routes
	http.HandleFunc("/creds", credsHandler)
	http.HandleFunc("/verify", handlers.VerifyHandler)
	http.HandleFunc("/issue", handlers.IssueHandler)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
