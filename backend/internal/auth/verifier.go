package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("auth: token inválido")

var allowedAlgorithms = []string{"ES256", "RS256"}

type Verifier struct {
	keyfunc  keyfunc.Keyfunc
	issuer   string
	audience string
}

func NewVerifier(ctx context.Context, jwksURL, issuer, audience string) (*Verifier, error) {
	kf, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("auth: erro ao carregar o JWKS de %s: %w", jwksURL, err)
	}

	return &Verifier{
		keyfunc:  kf,
		issuer:   issuer,
		audience: audience,
	}, nil
}

func (v *Verifier) Verify(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		v.keyfunc.Keyfunc,
		jwt.WithValidMethods(allowedAlgorithms),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("%w: claims em formato inesperado", ErrInvalidToken)
	}

	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return "", fmt.Errorf("%w: token sem claim sub", ErrInvalidToken)
	}

	return subject, nil
}

func TokenFromHeader(header string) (string, bool) {
	const prefix = "Bearer "

	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}

	token := strings.TrimSpace(header[len(prefix):])
	if token == "" {
		return "", false
	}
	return token, true
}
