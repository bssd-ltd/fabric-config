package certreader

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/youmark/pkcs8"
	"io/ioutil"
)

func ReadX509Cert(filepath string) *x509.Certificate {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	return cert
}

func ReadPrivateKey(filepath string) *ecdsa.PrivateKey {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		panic("failed to parse key PEM")
	}
	key, err := pkcs8.ParsePKCS8PrivateKeyECDSA(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}
