// Package main. types defines variables, constants and types.
package main

import (
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"log"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var (
	casdoorsdkClient *casdoorsdk.Client // Casdoor SDK client.
	privateKey       *rsa.PrivateKey    // Private key for signing JWTs. Casdoor must have the public key of this private key stored in cert.
)

var (
	//go:embed private_key.pem
	privateKeyPEM string

	//go:embed public_key.pem
	publicKeyPEM string // Public key for verification. This must be stored in casdoor -> identity -> cert -> <new_cert_name_chosen_in_application> -> Certificate
)

const (
	host                   = "http://localhost"
	port                   = ":8765"
	baseOauth2ClientAppUrl = host + port
	baseCasdoorUrl         = "http://localhost:8000"
	redirectUrl            = baseOauth2ClientAppUrl + "/callback" // redirectUrl is the callback URL for OAuth
	clientId               = "2afa0ecb38496579a666"               // Auto generated for application in casdoor
	orgName                = "built-in"                           // orgName is the name of the organization in casdoor
	appName                = "oauth-private_key_jwt"              // appName is the name of the application in casdoor
	sessionKey             = "sub"                                // sessionKey is the name of the session cookie
)

func init() {
	keyBytes := []byte(privateKeyPEM)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		log.Fatal("Failed to decode PEM block")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}
	privateKey = priv.(*rsa.PrivateKey)

	// Initialize the SDK client (cert parameter is the public key for verification)
	casdoorsdkClient = casdoorsdk.NewClient(
		baseCasdoorUrl,
		clientId,
		"",           // Client secret is not used for `private_key_jwt`
		publicKeyPEM, // This is the public certificate string generated from private key
		orgName,
		appName,
	)
}
