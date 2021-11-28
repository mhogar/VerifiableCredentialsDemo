package main

import (
	"flag"
	"fmt"
	"net/http"
	"vcd/user/handlers"
)

func main() {
	//log.Fatal(createVerifiableCredential())

	//parse flags
	port := flag.Int("port", 8080, "port to run the server on")
	flag.Parse()

	//setup routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/verify", handlers.VerifyHandler)

	//run the server
	fmt.Printf("listening on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
