package apple

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func verifyIdToken(aud string, idToken string) error {

	if idToken == "" {
		return errors.New("empty token")
	}

	//split and decode token
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return errors.New("invalid format token")
	}
	jsonHeader, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	var jwtHeader JWTTokenHeader
	err = json.Unmarshal(jsonHeader, &jwtHeader)
	if err != nil {
		return err
	}
	jsonBody, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	var jwtBody JWTTokenBody
	err = json.Unmarshal(jsonBody, &jwtBody)
	if err != nil {
		return err
	}

	// Verify that the iss field contains https://appleid.apple.com
	if jwtBody.Iss != "https://appleid.apple.com" {
		return errors.New("invalid iss field")
	}

	//Verify that the aud field is the developer’s client_id
	if jwtBody.Aud != aud {
		return errors.New("invalid aud field")
	}

	// Verify that the time is earlier than the exp value of the token
	if jwtBody.Exp < time.Now().Unix() {
		return errors.New("the token is expired")
	}

	//Verify the JWS E256 signature using the server’s public key
	var decodedSignature []byte
	decodedSignature, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	} else if !verifyAppleRSA256(parts[0]+"."+parts[1], decodedSignature, jwtHeader.Kid) {
		return errors.New("signature verification failed")
	}

	return nil
}

func verifyAppleRSA256(message string, signature []byte, kid string) bool {
	var rsaPublicKey *rsa.PublicKey
	var err error
	var hashed [32]byte

	rsaPublicKey = getApplePublicKeyObject(kid, "RS256")

	if rsaPublicKey != nil {
		bytesToHash := []byte(message)

		hashed = sha256.Sum256(bytesToHash)
		err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signature)
		if err != nil {
			return false
		}
	}
	return true
}
