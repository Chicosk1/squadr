package auth

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))
}

func TestMiddleware_TokenValidoSegueParaOHandler(t *testing.T) {
	key, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	const wantSubject = "e1f8a0d2-0000-4000-8000-000000000002"
	token := signToken(t, key, testKeyID, validClaims(wantSubject))

	var gotSubject string
	var handlerRodou bool

	handler := Middleware(verifier, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerRodou = true
		userID, ok := UserID(r.Context())
		if !ok {
			t.Error("UserID() não encontrou o ID no contexto da requisição")
		}
		gotSubject = userID
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !handlerRodou {
		t.Fatal("o handler protegido não foi executado")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, esperado %d", rec.Code, http.StatusOK)
	}
	if gotSubject != wantSubject {
		t.Errorf("UserID no contexto = %q, esperado %q", gotSubject, wantSubject)
	}
}

func TestMiddleware_RequisicoesRecusadas(t *testing.T) {
	_, jwksURL := jwksServer(t, testKeyID)
	verifier := newTestVerifier(t, jwksURL)

	outraChave, _ := jwksServer(t, testKeyID)

	tests := []struct {
		name       string
		authHeader string
	}{
		{"sem cabeçalho Authorization", ""},
		{"cabeçalho sem o prefixo Bearer", "abc.def.ghi"},
		{"esquema errado", "Basic dXN1YXJpbzpzZW5oYQ=="},
		{"token que não é JWT", "Bearer isto-nao-e-um-token"},
		{"assinado por outra chave", "Bearer " + signToken(t, outraChave, testKeyID, validClaims("atacante"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Middleware(verifier, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("o handler protegido não deveria ter sido executado")
			}))

			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, esperado %d", rec.Code, http.StatusUnauthorized)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("corpo do 401 não é JSON válido: %v\ncorpo: %s", err, rec.Body.String())
			}
			if body["code"] != "unauthorized" {
				t.Errorf("code = %v, esperado %q", body["code"], "unauthorized")
			}
			if body["message"] == "" || body["message"] == nil {
				t.Error("o 401 deveria trazer uma message exibível ao usuário")
			}
		})
	}
}

func TestUserID_ContextoSemMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	if userID, ok := UserID(req.Context()); ok {
		t.Errorf("UserID() deveria devolver ok=false fora do middleware, veio %q", userID)
	}
}
