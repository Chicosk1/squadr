package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL         string
	APIPort             int
	LogLevel            string
	SupabaseJWKSURL     string
	SupabaseJWTIssuer   string
	SupabaseJWTAudience string
}

func Load() (Config, error) {
	ld := &loader{}

	cfg := Config{
		DatabaseURL:         ld.require("DATABASE_URL"),
		APIPort:             ld.requireInt("API_PORT"),
		LogLevel:            ld.optional("LOG_LEVEL", "info"),
		SupabaseJWKSURL:     ld.require("SUPABASE_JWKS_URL"),
		SupabaseJWTIssuer:   ld.require("SUPABASE_JWT_ISSUER"),
		SupabaseJWTAudience: ld.require("SUPABASE_JWT_AUDIENCE"),
	}

	if err := ld.validationError(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

type loader struct {
	errs []error
}

func (ld *loader) require(key string) string {
	value := os.Getenv(key)
	if value == "" {
		ld.errs = append(ld.errs, fmt.Errorf("%s: variável de ambiente obrigatória não definida", key))
	}
	return value
}

func (ld *loader) requireInt(key string) int {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		ld.errs = append(ld.errs, fmt.Errorf("%s: variável de ambiente obrigatória não definida", key))
		return 0
	}

	parsed, err := strconv.Atoi(rawValue)
	if err != nil {
		ld.errs = append(ld.errs, fmt.Errorf("%s: valor %q não é um número inteiro válido", key, rawValue))
		return 0
	}
	return parsed
}

func (ld *loader) optional(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (ld *loader) validationError() error {
	if len(ld.errs) == 0 {
		return nil
	}
	return fmt.Errorf("config inválida:\n%w", errors.Join(ld.errs...))
}
