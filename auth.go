package twipla3as

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
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

type APIError struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Code       int    `json:"code"`
	OtherError string `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API error: %d %s (Code: %d)", e.Status, e.Message, e.Code)
}

func (sdk *TwiplaSDK) apiCall(ctx context.Context, method string, path string, body any) (*http.Response, error) {
	var r *http.Request
	query, ok := body.(url.Values)
	if !ok {
		query = nil
	} else {
		body = nil
	}

	finalPath := sdk.apiBase.JoinPath(path)
	if query != nil {
		finalPath.RawQuery = query.Encode()
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r, err = http.NewRequestWithContext(ctx, method, finalPath.String(), bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		r, err = http.NewRequestWithContext(ctx, method, finalPath.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	token, err := sdk.signer.IntpToken()
	if err != nil {
		return nil, fmt.Errorf("can't sign bearer intp token: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+token)
	if !(method == http.MethodGet || method == http.MethodDelete) {
		r.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		if !strings.Contains(resp.Header.Get("Content-Type"), "json") {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("non-json error response: %q", string(data))
		}
		var apiError APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, err
		}
		if apiError.OtherError == "invalid access token" {
			return nil, ErrInvalidAccessToken
		}
		return nil, apiError
	}

	return resp, nil
}

// PaginationMetadata is returned for paginated responses and shows metadata about the number of items that can be returned by the endpoint, as well as information about the current page.
type PaginationMetadata struct {
	// Page is the current page being queried
	Page int `json:"page"`
	// PageSize is the current page size
	PageSize int `json:"pageSize"`
	// PageTotal is the total number of pages available, based on the number of results.
	PageTotal int `json:"pageTotal"`
	// Total is the total number of results available.
	Total int `json:"total"`
}

type twiplaResponse[T any] struct {
	Payload  T                  `json:"payload"`
	Metadata PaginationMetadata `json:"meta"`
}

func parseResponse[T any](w *http.Response, err error) (twiplaResponse[T], error) {
	var t twiplaResponse[T]
	if err != nil {
		return t, err
	}
	defer w.Body.Close()

	if w.StatusCode == http.StatusNoContent {
		return t, nil
	}

	contentType, _, err := mime.ParseMediaType(w.Header.Get("Content-Type"))
	if err != nil {
		return t, fmt.Errorf("error parsing content type: %w", err)
	}
	if contentType != "application/json" {
		data, err := io.ReadAll(w.Body)
		defer w.Body.Close()
		if err != nil {
			return t, err
		}
		return t, fmt.Errorf("unexpected content type: %s, body: %q", w.Header.Get("Content-Type"), string(data))
	}
	if err := json.NewDecoder(w.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}
