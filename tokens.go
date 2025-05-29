package twipla3as

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (sdk *TwiplaSDK) IntpAccessToken() (string, error) {
	return sdk.signer.IntpToken()
}

func (sdk *TwiplaSDK) IntpcAccessToken(intpcID string) (string, error) {
	return sdk.signer.IntpcToken(intpcID)
}

type tokenSigner struct {
	privateKey *rsa.PrivateKey
	intpID     string
}

func (t *tokenSigner) IntpToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":     "twipla-3as-go-sdk",
		"roles":   []string{"intp"},
		"intp_id": t.intpID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 4).Unix(),
	})

	token.Header["kid"] = t.intpID
	return token.SignedString(t.privateKey)
}

func (t *tokenSigner) IntpcToken(intpcID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      "twipla-3as-go-sdk",
		"roles":    []string{"intpc"},
		"intp_id":  t.intpID,
		"intpc_id": intpcID,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 4).Unix(),
	})

	token.Header["kid"] = t.intpID
	return token.SignedString(t.privateKey)
}
