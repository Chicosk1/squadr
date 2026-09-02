package config

import (
	"strings"
	"testing"
)

const (
	testDatabaseURL = "postgresql://user:pass@host:5432/db"
	testJWKSURL     = "https://projeto.supabase.co/auth/v1/.well-known/jwks.json"
	testJWTIssuer   = "https://projeto.supabase.co/auth/v1"
	testJWTAudience = "authenticated"
)

func setValidEnv(t *testing.T) {
	t.Helper()

	t.Setenv("DATABASE_URL", testDatabaseURL)
	t.Setenv("API_PORT", "8080")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("SUPABASE_JWKS_URL", testJWKSURL)
	t.Setenv("SUPABASE_JWT_ISSUER", testJWTIssuer)
	t.Setenv("SUPABASE_JWT_AUDIENCE", testJWTAudience)
}

func TestLoad_Success(t *testing.T) {
	setValidEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() retornou erro inesperado: %v", err)
	}

	if cfg.DatabaseURL != testDatabaseURL {
		t.Errorf("DatabaseURL = %q, esperado %q", cfg.DatabaseURL, testDatabaseURL)
	}
	if cfg.APIPort != 8080 {
		t.Errorf("APIPort = %d, esperado 8080", cfg.APIPort)
	}
	if cfg.SupabaseJWKSURL != testJWKSURL {
		t.Errorf("SupabaseJWKSURL = %q, esperado %q", cfg.SupabaseJWKSURL, testJWKSURL)
	}
	if cfg.SupabaseJWTIssuer != testJWTIssuer {
		t.Errorf("SupabaseJWTIssuer = %q, esperado %q", cfg.SupabaseJWTIssuer, testJWTIssuer)
	}
	if cfg.SupabaseJWTAudience != testJWTAudience {
		t.Errorf("SupabaseJWTAudience = %q, esperado %q", cfg.SupabaseJWTAudience, testJWTAudience)
	}
}

func TestLoad_LogLevelTemDefault(t *testing.T) {
	setValidEnv(t)
	t.Setenv("LOG_LEVEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() retornou erro inesperado: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, esperado o default %q", cfg.LogLevel, "info")
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	obrigatorias := []string{
		"DATABASE_URL",
		"API_PORT",
		"SUPABASE_JWKS_URL",
		"SUPABASE_JWT_ISSUER",
		"SUPABASE_JWT_AUDIENCE",
	}

	setValidEnv(t)
	for _, key := range obrigatorias {
		t.Setenv(key, "")
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() deveria falhar sem as variáveis obrigatórias")
	}
	for _, key := range obrigatorias {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("erro deveria mencionar %s, veio: %v", key, err)
		}
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("API_PORT", "nao-e-um-numero")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() deveria falhar com API_PORT inválido")
	}
	if !strings.Contains(err.Error(), "API_PORT") {
		t.Errorf("erro deveria mencionar API_PORT, veio: %v", err)
	}
}
