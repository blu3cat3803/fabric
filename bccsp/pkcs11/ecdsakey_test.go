// +build pkcs11

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package pkcs11

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/stretchr/testify/assert"
)

func TestECDSAPKIXPublicKeyImportOptsKeyImporter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestECDSAPKIXPublicKeyImportOptsKeyImporter")
	}
	ki := currentBCCSP

	_, err := ki.KeyImport("Hello World", &bccsp.ECDSAPKIXPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSAPKIXPublicKeyImportOpts] Invalid raw material. Expected byte array")

	_, err = ki.KeyImport("Hello World", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Opts parameter. It must not be nil")

	_, err = ki.KeyImport(nil, &bccsp.ECDSAPKIXPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid raw. Cannot be nil")

	_, err = ki.KeyImport([]byte(nil), &bccsp.ECDSAPKIXPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSAPKIXPublicKeyImportOpts] Invalid raw. It must not be nil")

	_, err = ki.KeyImport([]byte{0}, &bccsp.ECDSAPKIXPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed converting PKIX to ECDSA public key")

	k, err := rsa.GenerateKey(rand.Reader, 512)
	assert.NoError(t, err)
	raw, err := utils.PublicKeyToDER(&k.PublicKey)
	assert.NoError(t, err)
	_, err = ki.KeyImport(raw, &bccsp.ECDSAPKIXPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed casting to ECDSA public key. Invalid raw material")
}

func TestECDSAPrivateKeyImportOptsKeyImporter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestECDSAPrivateKeyImportOptsKeyImporter")
	}

	ki := currentBCCSP

	// set privImport to false
	ki.(*impl).privImport = false
	_, err := ki.KeyImport([]byte{0}, &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSADERPrivateKeyImportOpts] PKCS11 options 'sensitivekeys' is set to true. Cannot import")

	// set privImport to true
	ki.(*impl).privImport = true

	_, err = ki.KeyImport("Hello World", &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSADERPrivateKeyImportOpts] Invalid raw material. Expected byte array")

	_, err = ki.KeyImport(nil, &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid raw. Cannot be nil")

	_, err = ki.KeyImport([]byte(nil), &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSADERPrivateKeyImportOpts] Invalid raw. It must not be nil")

	_, err = ki.KeyImport([]byte{0}, &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed converting PKIX to ECDSA public key")

	k, err := rsa.GenerateKey(rand.Reader, 512)
	assert.NoError(t, err)
	raw := x509.MarshalPKCS1PrivateKey(k)
	_, err = ki.KeyImport(raw, &bccsp.ECDSAPrivateKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed casting to ECDSA public key. Invalid raw material")

}

func TestECDSAGoPublicKeyImportOptsKeyImporter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestECDSAGoPublicKeyImportOptsKeyImporter")
	}
	ki := currentBCCSP

	_, err := ki.KeyImport("Hello World", &bccsp.ECDSAGoPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ECDSAGoPublicKeyImportOpts] Invalid raw material. Expected *ecdsa.PublicKey")

	_, err = ki.KeyImport(nil, &bccsp.ECDSAGoPublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid raw. Cannot be nil")
}

func TestX509PublicKeyImportOptsKeyImporter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestX509PublicKeyImportOptsKeyImporter")
	}
	ki := currentBCCSP

	_, err := ki.KeyImport("Hello World", &bccsp.X509PublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[X509PublicKeyImportOpts] Invalid raw material. Expected *x509.Certificate")

	_, err = ki.KeyImport(nil, &bccsp.X509PublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid raw. Cannot be nil")

	cert := &x509.Certificate{}
	cert.PublicKey = "Hello world"
	_, err = ki.KeyImport(cert, &bccsp.X509PublicKeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Certificate's public key type not recognized. Supported keys: [ECDSA, RSA]")
}
