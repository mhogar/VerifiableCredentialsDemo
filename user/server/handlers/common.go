package handlers

import (
	"errors"
	"io"
	"net/http"
	"vcd/common"
)

const VC_URI = "wallet/verifiable-credentials.json"

type CredentialsMap map[string]common.VerifiableCredential

func sendRequest(method string, url string, body interface{}) (io.ReadCloser, CustomError, error) {
	var buffer io.Reader = nil

	if body != nil {
		buffer, _ = common.EncodeJSON(body)
	}

	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, InternalError(), common.ChainError("error creating request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, InternalError(), common.ChainError("error sending request", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		result := common.ErrorResponse{}
		common.DecodeJSON(res.Body, &result)

		err := common.ChainError("error from server", errors.New(result.Error))

		if res.StatusCode == http.StatusInternalServerError {
			return nil, InternalError(), err
		}

		return nil, ClientError(result.Error), err
	}

	return res.Body, NoError(), nil
}

func loadVerifiableCredentials() (*CredentialsMap, error) {
	creds := CredentialsMap{}

	err := common.LoadJSONFromFile(VC_URI, &creds)
	if err != nil {
		return nil, common.ChainError("error loading JSON file", err)
	}

	return &creds, nil
}
