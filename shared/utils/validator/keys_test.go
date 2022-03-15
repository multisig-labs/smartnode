package validator

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetValidatorPrivateKey(t *testing.T) {
	// read file at path, put into raw bytes
	raw, err := ioutil.ReadFile("/home/chandler/Downloads/staker1.key")
	if err != nil {
		t.Error(err)
	}

	// decode raw bytes into private key
	block, _ := pem.Decode(raw)

	fmt.Println(block.Type)

	if block.Type != "RSA PRIVATE KEY" {
		t.Error("private key not found")
	}

	// decode into x509 private key
	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Error(err)
	}
}
