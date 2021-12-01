package common

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

type DIDDocument struct {
	Domain     string            `json:"domain"`
	Routes     map[string]string `json:"routes"`
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
	keyRoute, ok := doc.Routes["key"]
	if !ok {
		return nil, errors.New("DID doc has no route for public key")
	}
	url := "http://" + path.Join(doc.Domain, keyRoute)

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
	return WriteJSONToFile(getFullURI(uri), doc)
}
