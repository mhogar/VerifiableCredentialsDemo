package main

import (
	"flag"
	"log"
	"vcd/common"
)

func Run(keyURI string, did string, docURI string) error {
	doc, err := common.LoadDIDDocumentFromURI(docURI)
	if err != nil {
		return common.ChainError("error loading DID doc from URI", err)
	}

	sigs := doc.Signatures
	doc.Signatures = nil

	sig := common.Signature{}
	err = common.SignStruct(keyURI, &sig, &doc)
	if err != nil {
		return common.ChainError("error signing DID doc", err)
	}

	sigs[did] = sig.Signature
	doc.Signatures = sigs

	err = common.SaveDIDDocument(docURI, doc)
	if err != nil {
		return common.ChainError("error saving DID doc", err)
	}

	return nil
}

func main() {
	keyURI := flag.String("key", "", "URI of the private key")
	did := flag.String("did", "", "DID of the signature")
	docURI := flag.String("doc", "", "URI of the DID document to sign")
	flag.Parse()

	err := Run(*keyURI, *did, *docURI)
	if err != nil {
		log.Fatal(err)
	}
}
