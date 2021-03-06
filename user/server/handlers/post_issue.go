package handlers

import (
	"log"
	"net/http"
	"os"
	"vcd/common"
)

const DID_URI = "wallet/DID.cert"

type IssuePostBody struct {
	ServiceURL   string            `json:"service_url"`
	Type         string            `json:"type"`
	Fields       map[string]string `json:"fields,omitempty"`
	CredentialID string            `json:"credential_id,omitempty"`
}

func PostIssueHandler(w http.ResponseWriter, req *http.Request) {
	body := IssuePostBody{}

	err := common.DecodeJSON(req.Body, &body)
	if err != nil {
		common.LogChainError("error decoding post issue body", err)
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cerr := postIssue(&body)
	if cerr.Type == TypeClientError {
		common.SendErrorResponse(w, http.StatusBadRequest, cerr.Message)
		return
	}
	if cerr.Type == TypeInternalError {
		common.SendInternalErrorResponse(w)
		return
	}

	common.SendSuccessResponse(w)
}

func postIssue(body *IssuePostBody) CustomError {
	cred := common.VerifiableCredential{}

	if body.Type == "iss:form" {
		bytes, err := os.ReadFile(DID_URI)
		if err != nil {
			common.LogChainError("error reading DID file", err)
			return InternalError()
		}
		cred.Subject.DID = string(bytes)

		cred.Credentials = body.Fields

	} else { //iss:cred
		creds, err := loadVerifiableCredentials()
		if err != nil {
			common.LogChainError("error loading verifiable credentials", err)
			return InternalError()
		}

		var ok bool
		cred, ok = (*creds)[body.CredentialID]
		if !ok {
			log.Println("credential with id", body.CredentialID, "no found")
			return ClientError("No credential found for ID.")
		}
	}

	err := common.SignStruct(PRIVATE_KEY_URI, &cred.Subject, cred)
	if err != nil {
		common.LogChainError("error signing issue request", err)
		return InternalError()
	}

	res, cerr, err := sendRequest(http.MethodPost, body.ServiceURL, &cred)
	if err != nil {
		log.Println(err)
	}
	if cerr.Type != TypeNoError {
		return cerr
	}
	defer res.Close()

	cred = common.VerifiableCredential{}
	err = common.DecodeJSON(res, &cred)
	if err != nil {
		common.LogChainError("error decoding verifiable credential", err)
		return InternalError()
	}

	creds := CredentialsMap{}
	err = common.LoadJSONFromFile(VC_URI, &creds)
	if err != nil {
		common.LogChainError("error loading verifiable credentials", err)
		return InternalError()
	}

	creds[cred.Issuer.DID] = cred
	err = common.WriteJSONToFile(VC_URI, &creds)
	if err != nil {
		common.LogChainError("error saving verifiable credentials", err)
		return InternalError()
	}

	return NoError()
}
