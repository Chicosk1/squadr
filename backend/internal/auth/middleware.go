package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type contextKey struct{}

var userIDKey contextKey

func Middleware(verifier *Verifier, appLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := TokenFromHeader(r.Header.Get("Authorization"))
			if !ok {
				unauthorized(w, appLogger, "cabeçalho Authorization ausente ou malformado", nil)
				return
			}

			userID, err := verifier.Verify(token)
			if err != nil {
				unauthorized(w, appLogger, "token recusado", err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func unauthorized(w http.ResponseWriter, appLogger *slog.Logger, reason string, err error) {
	appLogger.Info("requisição não autenticada", "motivo", reason, "erro", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	body := map[string]string{
		"code":    "unauthorized",
		"message": "Sessão inválida ou expirada. Entre novamente.",
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		appLogger.Error("erro ao escrever resposta 401", "erro", err)
	}
}
