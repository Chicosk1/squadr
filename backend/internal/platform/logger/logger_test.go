package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"info explícito", "info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"maiúsculas", "DEBUG", slog.LevelDebug},
		{"vazio cai no default", "", slog.LevelInfo},
		{"desconhecido cai no default", "verbose", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLevel(tt.input)
			if got != tt.want {
				t.Errorf("parseLevel(%q) = %v, esperado %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewLogger_RespeitaONivel(t *testing.T) {
	var buf bytes.Buffer
	log := newLogger("warn", &buf)

	log.Info("não deveria aparecer")
	log.Warn("deveria aparecer")

	output := buf.String()
	if strings.Contains(output, "não deveria aparecer") {
		t.Error("mensagem de nível Info apareceu mesmo com LOG_LEVEL=warn")
	}
	if !strings.Contains(output, "deveria aparecer") {
		t.Error("mensagem de nível Warn deveria ter aparecido")
	}
}

func TestNewLogger_FormatoJSON(t *testing.T) {
	var buf bytes.Buffer
	log := newLogger("info", &buf)

	log.Info("mensagem de teste", "chave", "valor")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("saída não é JSON válido: %v\nsaída: %s", err, buf.String())
	}
	if entry["msg"] != "mensagem de teste" {
		t.Errorf("campo msg = %v, esperado %q", entry["msg"], "mensagem de teste")
	}
	if entry["chave"] != "valor" {
		t.Errorf("campo chave = %v, esperado %q", entry["chave"], "valor")
	}
}
