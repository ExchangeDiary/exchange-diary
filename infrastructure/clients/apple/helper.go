package apple

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type AppleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

const APPLE_BASE_URL = "https://appleid.apple.com"

func getApplePublicKeyObject(kid string, alg string) *rsa.PublicKey {

	client := NewClient(APPLE_BASE_URL)

	applePublicKey := client.getApplePublicKey(kid, alg)

	if applePublicKey == nil {
		return nil
	}

	key := getPublicKeyObject(applePublicKey.N, applePublicKey.E)
	return key
}

func getPublicKeyObject(encodedN string, encodedE string) *rsa.PublicKey {

	var publicKey rsa.PublicKey
	var eInt int

	decodedN, err := base64.RawURLEncoding.DecodeString(encodedN)
	if err != nil {
		return nil
	}
	publicKey.N = new(big.Int)
	publicKey.N.SetBytes(decodedN)

	decodedE, err := base64.RawURLEncoding.DecodeString(encodedE)
	if err != nil {
		return nil
	}

	for _, v := range decodedE {
		eInt = eInt << 8
		eInt = eInt | int(v)
	}
	publicKey.E = eInt

	return &publicKey
}
