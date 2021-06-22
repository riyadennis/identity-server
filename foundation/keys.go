package foundation

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const (
	PrivateKey = "RSA PRIVATE KEY"
	PublicKey  = "PUBLIC KEY"
)

func GenerateKeys(privateKeyPath, publicKeyPath string) error {
	privateKey, keyErr := rsa.GenerateKey(rand.Reader, 2048)
	if keyErr != nil {
		return keyErr
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  PrivateKey,
		Bytes: privateKeyBytes,
	}

	privatePem, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		return err
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicKeyBlock := &pem.Block{
		Type:  PublicKey,
		Bytes: publicKeyBytes,
	}

	publicPem, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}

	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		return err
	}

	return nil
}
