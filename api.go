package twipla3as

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

type Environment string

const (
	EnvironmentDevelop    Environment = "dev"
	EnvironmentStage      Environment = "stage"
	EnvironmentProduction Environment = "production"
)

var (
	ErrNoPrivateKey            = errors.New("no private key provided")
	ErrInvalidSubscriptionType = errors.New("invalid subscription type")
	ErrInvalidAccessToken      = errors.New("invalid access token")
)

type TwiplaConfig struct {
	// IntpID is the ID issued by TWIPLA for an INTP, in exchange for `jwtRS256.key.pub`
	IntpID string
	// PrivateKey contains the plaintext contents of the PEM file containing the INTP's private key
	PrivateKey string

	// Environment sets which TWIPLA deployment to use. If not [EnvironmentDevelop] or [EnvironmentStage], its value is assumed to be [EnvironmentProduction]
	Environment Environment
}

type TwiplaSDK struct {
	signer  *tokenSigner
	client  *http.Client
	env     Environment
	apiBase *url.URL
}

func NewSDK(config *TwiplaConfig) (*TwiplaSDK, error) {
	if config.PrivateKey == "" {
		return nil, ErrNoPrivateKey
	}

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.PrivateKey))
	if err != nil {
		return nil, err
	}

	signer := &tokenSigner{
		privateKey: pkey,
		intpID:     config.IntpID,
	}

	var apiPrefix string
	switch config.Environment {
	case EnvironmentDevelop:
		apiPrefix = "https://api-gateway.va-endpoint.com"
	case EnvironmentStage:
		apiPrefix = "https://stage-api-gateway.va-endpoint.com"
	default:
		apiPrefix = "https://api-gateway.visitor-analytics.io"
		config.Environment = EnvironmentProduction
	}

	apiURL, err := url.Parse(apiPrefix)

	return &TwiplaSDK{
		signer:  signer,
		env:     config.Environment,
		apiBase: apiURL,
	}, nil
}
