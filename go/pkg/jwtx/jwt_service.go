package jwtx

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Options is a struct for specifying configuration options for the JWT service.
type Options struct {
	JWKS string `json:"jwks"`
	Aud  string `json:"aud"`
	Iss  string `json:"iss"`
}

type JWTService struct {
	Options Options

	keyset jwk.Set
}

// NewService constructs a new JWTService instance with supplied options.
func NewService(options ...Options) (*JWTService, error) {

	var opts Options
	if len(options) == 0 {
		opts = Options{}
	} else {
		opts = options[0]
	}

	var keyset jwk.Set

	if opts.JWKS != "" {
		var err error
		keyset, err = fetchKeyset(opts.JWKS)
		if err != nil {
			return nil, err
		}
	}

	return &JWTService{
		Options: opts,
		keyset:  keyset,
	}, nil
}

func fetchKeyset(url string) (jwk.Set, error) {
	keyset, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWK: %w", err)
	}

	return keyset, nil
}

func (s *JWTService) Parse(payload []byte) (jwt.Token, error) {
	token, err := jwt.Parse(payload,
		jwt.WithKeySet(s.keyset),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	return token, nil
}

func (s *JWTService) Validate(token jwt.Token) error {
	options := make([]jwt.ValidateOption, 0)

	if s.Options.Aud != "" {
		// Verify 'aud' claim
		options = append(options, jwt.WithAudience(s.Options.Aud))
	}

	if s.Options.Iss != "" {
		// Verify 'iss' claim
		options = append(options, jwt.WithIssuer(s.Options.Iss))
	}

	err := jwt.Validate(token, options...)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	return nil
}
