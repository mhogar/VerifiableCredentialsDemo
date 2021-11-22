package common

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

type DIDDocument struct {
	Domain string `json:"domain"`
	Route  string `json:"route"`
}

type DIDLoader interface {
	LoadDIDDocument(uri string) (*DIDDocument, error)
	LoadPublicKey(uri string) ([]byte, error)
}

type DIDFileLoader struct{}

func (DIDFileLoader) LoadDIDDocumentFromURI(uri string) (*DIDDocument, error) {
	//open file
	f, err := os.Open(path.Join("..", "blockchain", uri))
	if err != nil {
		return nil, ChainError("error opening DID document file", err)
	}
	defer f.Close()

	//decode json
	doc := DIDDocument{}
	err = DecodeJSON(f, &doc)
	if err != nil {
		return nil, ChainError("error decoding JSON", err)
	}

	return &doc, nil
}

func (DIDFileLoader) LoadPublicKeyFromDocument(doc *DIDDocument) ([]byte, error) {
	url := "http://" + path.Join(doc.Domain, doc.Route)

	//load issuer DID
	res, err := http.Get(url)
	if err != nil {
		return nil, ChainError("error sending DID Get request", err)
	}
	defer res.Body.Close()

	//check request was successful
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error getting DID from url: " + url)
	}

	//read body
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ChainError("error reading request body", err)
	}

	return bytes, nil
}

func (l DIDFileLoader) LoadPublicKeyFromURI(uri string) ([]byte, error) {
	doc, err := l.LoadDIDDocumentFromURI(uri)
	if err != nil {
		return nil, ChainError("error loading DID document", err)
	}

	return l.LoadPublicKeyFromDocument(doc)
}
