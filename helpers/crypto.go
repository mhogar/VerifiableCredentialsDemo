package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func SignResponse(keyURI string, res interface{}) ([]byte, error) {
	//load private key
	key, err := loadPrivateKeyFromFile(keyURI)
	if err != nil {
		return nil, ChainError("error loading private key from file", err)
	}

	//marshal response as json
	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, ChainError("error marshaling json", err)
	}

	//hash and sign
	sum := sha256.Sum256(bytes)
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		return nil, ChainError("error signing hash", err)
	}

	return sig, nil
}

func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	//read file
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ChainError("error reading private key file", err)
	}

	//parse PEM block
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("error parsing PEM block")
	}

	//parse private key
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, ChainError("error parsing private key from PEM bytes", err)
	}

	return key, nil
}
