package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testIssuer   = "https://projeto-de-teste.supabase.co/auth/v1"
	testAudience = "authenticated"
	testKeyID    = "chave-de-teste"
)

func jwksServer(t *testing.T, keyID string) (*ecdsa.PrivateKey, string) {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("erro ao gerar chave de teste: %v", err)
	}

	coordinateSize := 32
	x := base64.RawURLEncoding.EncodeToString(privateKey.X.FillBytes(make([]byte, coordinateSize)))
	y := base64.RawURLEncoding.EncodeToString(privateKey.Y.FillBytes(make([]byte, coordinateSize)))

	jwks := map[string]any{
		"keys": []map[string]any{{
			"kty": "EC",
			"crv": "P-256",
			"alg": "ES256",
			"use": "sig",
			"kid": keyID,
			"x":   x,
			"y":   y,
		}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			t.Errorf("erro ao servir o JWKS de teste: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return privateKey, server.URL
}

func signToken(t *testing.T, key *ecdsa.PrivateKey, keyID string, claims jwt.RegisteredClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyID

	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("erro ao assinar token de teste: %v", err)
	}
	return signed
}

func validClaims(subject string) jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    testIssuer,
		Audience:  jwt.ClaimStrings{testAudience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	}
}

func newTestVerifier(t *testing.T, jwksURL string) *Verifier {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	verifier, err := NewVerifier(ctx, jwksURL, testIssuer, testAudience)
	if err != nil {
		t.Fatalf("NewVerifier() falhou: %v", err)
	}
	return verifier
}

func TestVerify_TokenValido(t *testing.T) {
	key, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	const wantSubject = "e1f8a0d2-0000-4000-8000-000000000001"
	token := signToken(t, key, testKeyID, validClaims(wantSubject))

	subject, err := verifier.Verify(token)
	if err != nil {
		t.Fatalf("Verify() recusou um token válido: %v", err)
	}
	if subject != wantSubject {
		t.Errorf("Verify() = %q, esperado o sub %q", subject, wantSubject)
	}
}

func TestVerify_TokensRecusados(t *testing.T) {
	key, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	outraChave, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("erro ao gerar chave impostora: %v", err)
	}

	expiradas := validClaims("usuario")
	expiradas.IssuedAt = jwt.NewNumericDate(time.Now().Add(-2 * time.Hour))
	expiradas.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))

	emissorErrado := validClaims("usuario")
	emissorErrado.Issuer = "https://projeto-do-atacante.supabase.co/auth/v1"

	audienciaErrada := validClaims("usuario")
	audienciaErrada.Audience = jwt.ClaimStrings{"outra-audiencia"}

	semExpiracao := validClaims("usuario")
	semExpiracao.ExpiresAt = nil

	semSub := validClaims("")

	tests := []struct {
		name  string
		token string
	}{
		{"assinado por outra chave", signToken(t, outraChave, testKeyID, validClaims("usuario"))},
		{"expirado", signToken(t, key, testKeyID, expiradas)},
		{"emissor errado", signToken(t, key, testKeyID, emissorErrado)},
		{"audiência errada", signToken(t, key, testKeyID, audienciaErrada)},
		{"sem claim exp", signToken(t, key, testKeyID, semExpiracao)},
		{"sem claim sub", signToken(t, key, testKeyID, semSub)},
		{"kid desconhecido", signToken(t, key, "kid-que-nao-existe", validClaims("usuario"))},
		{"texto que não é JWT", "isto-nao-e-um-token"},
		{"string vazia", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject, err := verifier.Verify(tt.token)
			if err == nil {
				t.Fatalf("Verify() aceitou um token que deveria recusar (sub devolvido: %q)", subject)
			}
			if !errors.Is(err, ErrInvalidToken) {
				t.Errorf("erro deveria envolver ErrInvalidToken, veio: %v", err)
			}
			if subject != "" {
				t.Errorf("Verify() devolveu sub %q junto com erro — deveria ser vazio", subject)
			}
		})
	}
}

func TestVerify_RecusaConfusaoDeAlgoritmo(t *testing.T) {
	_, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims("atacante"))
	token.Header["kid"] = testKeyID

	signed, err := token.SignedString([]byte("chave-publica-usada-como-segredo"))
	if err != nil {
		t.Fatalf("erro ao assinar token HS256 de teste: %v", err)
	}

	if _, err := verifier.Verify(signed); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Verify() deveria recusar HS256, veio: %v", err)
	}
}

func TestTokenFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{"bearer normal", "Bearer abc.def.ghi", "abc.def.ghi", true},
		{"prefixo em caixa diferente", "bearer abc.def.ghi", "abc.def.ghi", true},
		{"espaços em volta do token", "Bearer   abc.def.ghi  ", "abc.def.ghi", true},
		{"header vazio", "", "", false},
		{"sem o prefixo Bearer", "abc.def.ghi", "", false},
		{"esquema errado", "Basic abc.def.ghi", "", false},
		{"prefixo sem token", "Bearer ", "", false},
		{"só espaços depois do prefixo", "Bearer    ", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, ok := TokenFromHeader(tt.header)
			if ok != tt.wantOK {
				t.Errorf("TokenFromHeader(%q) ok = %v, esperado %v", tt.header, ok, tt.wantOK)
			}
			if token != tt.wantToken {
				t.Errorf("TokenFromHeader(%q) = %q, esperado %q", tt.header, token, tt.wantToken)
			}
		})
	}
}

func TestVerify_ErroPreservaMotivoParaLog(t *testing.T) {
	key, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	expiradas := validClaims("usuario")
	expiradas.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))

	_, err := verifier.Verify(signToken(t, key, testKeyID, expiradas))
	if err == nil {
		t.Fatal("esperava erro para token expirado")
	}
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("erro deveria preservar jwt.ErrTokenExpired para o log, veio: %v", err)
	}
}
