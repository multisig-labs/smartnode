package validator

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func GetValidatorPrivateKey(path string) (*rsa.PrivateKey, error) {
	// read file at path, put into raw bytes
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// decode raw bytes into private key
	block, _ := pem.Decode(raw)

	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("private key not found")
	}

	// decode into x509 private key
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsaKey, nil
}
