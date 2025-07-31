package foundation

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeys(t *testing.T) {
	privateKeyFile, err := os.CreateTemp("", "test_private_*.pem")
	assert.NoError(t, err, "failed to create temp private key file")
	defer os.Remove(privateKeyFile.Name())
	privateKeyFile.Close()

	publicKeyFile, err := os.CreateTemp("", "test_public_*.pem")
	assert.NoError(t, err, "failed to create temp public key file")
	defer os.Remove(publicKeyFile.Name())
	publicKeyFile.Close()

	err = GenerateKeys(privateKeyFile.Name(), publicKeyFile.Name())
	assert.NoError(t, err, "GenerateKeys failed")

	// Check private key file
	privData, err := os.ReadFile(privateKeyFile.Name())
	assert.NoError(t, err, "failed to read private key file")
	privBlock, _ := pem.Decode(privData)
	if assert.NotNil(t, privBlock, "invalid private key PEM block") {
		assert.Equal(t, PrivateKey, privBlock.Type, "invalid private key PEM type")
		_, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		assert.NoError(t, err, "invalid private key")
	}

	// Check public key file
	pubData, err := os.ReadFile(publicKeyFile.Name())
	assert.NoError(t, err, "failed to read public key file")
	pubBlock, _ := pem.Decode(pubData)
	if assert.NotNil(t, pubBlock, "invalid public key PEM block") {
		assert.Equal(t, PublicKey, pubBlock.Type, "invalid public key PEM type")
		_, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
		assert.NoError(t, err, "invalid public key")
	}
}
