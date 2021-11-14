package helpers

import (
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
