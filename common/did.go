package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

type DIDDocument struct {
	Name       string            `json:"name"`
	Domain     string            `json:"domain"`
	Route      string            `json:"route"`
	Signatures map[string]string `json:"signatures,omitempty"`
}

func getFullURI(uri string) string {
	return path.Join("..", "blockchain", uri+".json")
}

func LoadDIDDocumentFromURI(uri string) (*DIDDocument, error) {
	f, err := os.Open(getFullURI(uri))
	if err != nil {
		return nil, ChainError("error opening DID document file", err)
	}
	defer f.Close()

	doc := DIDDocument{}
	err = DecodeJSON(f, &doc)
	if err != nil {
		return nil, ChainError("error decoding JSON", err)
	}

	return &doc, nil
}

func LoadPublicKeyFromDocument(doc *DIDDocument) ([]byte, error) {
	url := "http://" + path.Join(doc.Domain, doc.Route)

	res, err := http.Get(url)
	if err != nil {
		return nil, ChainError("error sending DID Get request", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error getting DID from url: " + url)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ChainError("error reading request body", err)
	}

	return bytes, nil
}

func LoadPublicKeyFromURI(uri string) ([]byte, error) {
	doc, err := LoadDIDDocumentFromURI(uri)
	if err != nil {
		return nil, ChainError("error loading DID document", err)
	}

	return LoadPublicKeyFromDocument(doc)
}

func SaveDIDDocument(uri string, doc *DIDDocument) error {
	f, err := os.Create(getFullURI(uri))
	if err != nil {
		return ChainError("error creating DID document file", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(doc)
	if err != nil {
		return ChainError("error encoding/writing DID document", err)
	}

	return nil
}
