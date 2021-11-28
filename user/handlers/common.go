package handlers

import (
	"errors"
	"io"
	"net/http"
	"vcd/common"
)

func sendRequest(method string, url string, body interface{}) (io.ReadCloser, error) {
	var buffer io.Reader = nil

	if body != nil {
		buffer, _ = common.EncodeJSON(body)
	}

	//create request
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, common.ChainError("error creating request", err)
	}

	//send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, common.ChainError("error sending request", err)
	}

	//handle error response
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		result := common.ErrorResponse{}
		common.DecodeJSON(res.Body, &result)

		return nil, common.ChainError("error from server", errors.New(result.Error))
	}

	return res.Body, nil
}
