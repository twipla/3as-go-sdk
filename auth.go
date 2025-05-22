package twipla3as

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func (sdk *TwiplaSDK) apiCall(ctx context.Context, method string, path string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, method, sdk.apiBase.JoinPath(path).String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	token, err := sdk.signer.IntpToken()
	if err != nil {
		return nil, fmt.Errorf("can't sign bearer intp token: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+token)
	if !(method == http.MethodGet || method == http.MethodDelete) {
		r.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(r)
}
