package common

import (
	"bytes"
	"encoding/json"
	"io"
)

func DecodeJSON(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func EncodeJSON(v interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	//encode json to buffer
	err := json.NewEncoder(buffer).Encode(v)
	if err != nil {
		return nil, ChainError("error encoding JSON", err)
	}

	return buffer, nil
}
