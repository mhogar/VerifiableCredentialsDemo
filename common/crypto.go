package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func SignStruct(keyURI string, sig *Signature, v interface{}) error {
	key, err := loadPrivateKeyFromFile(keyURI)
	if err != nil {
		return ChainError("error loading private key from file", err)
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return ChainError("error marshaling json", err)
	}

	hash := sha256.Sum256(bytes)
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return ChainError("error signing hash", err)
	}

	sig.Signature = base64.RawStdEncoding.EncodeToString(sigBytes)
	return nil
}

func VerifyStructSignature(DID []byte, sig *Signature, v interface{}) error {
	sigStr := sig.Signature
	sig.Signature = ""

	sigBytes, err := base64.RawStdEncoding.DecodeString(sigStr)
	if err != nil {
		return ChainError("error decoding signature", err)
	}

	key, err := loadPublicKeyFromBytes(DID)
	if err != nil {
		return ChainError("error loading public key from DID", err)
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return ChainError("error marshaling json", err)
	}

	hash := sha256.Sum256(bytes)
	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sigBytes)
	if err != nil {
		return ChainError("error verifying signature", err)
	}

	return nil
}

func LoadKeyFromFile(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ChainError("error reading file", err)
	}

	return bytes, nil
}

func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	//read file
	bytes, err := LoadKeyFromFile(filename)
	if err != nil {
		return nil, ChainError("error reading private key file", err)
	}

	//parse PEM block
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("error parsing PEM block")
	}

	//parse private key
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, ChainError("error parsing private key from PEM bytes", err)
	}

	//verify key is RSA private key
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not an RSA private key")
	}

	return rsaKey, nil
}

func loadPublicKeyFromBytes(bytes []byte) (*rsa.PublicKey, error) {
	//parse PEM block
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("error parsing PEM block")
	}

	//parse public cert
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, ChainError("error parsing certificate", err)
	}

	//verify key is RSA public key
	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an RSA public key")
	}

	return key, nil
}
