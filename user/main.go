package main

import (
	"fmt"
	"log"
	"net/http"
	"vcd/common"
	"vcd/helpers"
)

func main() {
	credFields := common.VerifiableCredential{
		SubjectDID: "Subject DID",
		FirstName:  "Bob",
		LastName:   "ADobDob",
	}

	//encode request body
	buffer, err := helpers.EncodeJSON(&credFields)
	if err != nil {
		log.Fatal(err)
	}

	//send request
	res, err := http.Post("http://localhost:8082/creds", "application/json", buffer)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//select result struct type
	var result interface{}
	if res.StatusCode == http.StatusOK {
		result = common.VerifiableCredential{}
	} else {
		result = common.ErrorResponse{}
	}

	//decode result
	err = helpers.DecodeJSON(res.Body, &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
