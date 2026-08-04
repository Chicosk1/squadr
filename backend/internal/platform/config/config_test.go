package config

import (
	"strings"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://user:pass@host:5432/db")
	t.Setenv("API_PORT", "8080")
	t.Setenv("LOG_LEVEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() retornou erro inesperado: %v", err)
	}

	if cfg.DatabaseURL != "postgresql://user:pass@host:5432/db" {
		t.Errorf("DatabaseURL = %q, esperado o valor do env", cfg.DatabaseURL)
	}
	if cfg.APIPort != 8080 {
		t.Errorf("APIPort = %d, esperado 8080", cfg.APIPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, esperado o default %q", cfg.LogLevel, "info")
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("API_PORT", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() deveria falhar sem DATABASE_URL/API_PORT definidas")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("erro deveria mencionar DATABASE_URL, veio: %v", err)
	}
	if !strings.Contains(err.Error(), "API_PORT") {
		t.Errorf("erro deveria mencionar API_PORT, veio: %v", err)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://user:pass@host:5432/db")
	t.Setenv("API_PORT", "nao-e-um-numero")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() deveria falhar com API_PORT inválido")
	}
	if !strings.Contains(err.Error(), "API_PORT") {
		t.Errorf("erro deveria mencionar API_PORT, veio: %v", err)
	}
}
