package common

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
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

func LoadJSONFromFile(uri string, v interface{}) error {
	f, err := os.Open(uri)
	if err != nil {
		return ChainError("error opening JSON file", err)
	}
	defer f.Close()

	return DecodeJSON(f, v)
}

func WriteJSONToFile(uri string, v interface{}) error {
	f, err := os.Create(uri)
	if err != nil {
		return ChainError("error creating JSON file", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(v)
	if err != nil {
		return ChainError("error encoding/writing JSON file", err)
	}

	return nil
}
