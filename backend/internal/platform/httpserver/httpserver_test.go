package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := requestLogger(testLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/qualquer-coisa", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log não é JSON válido: %v\nsaída: %s", err, buf.String())
	}
	if entry["method"] != http.MethodGet {
		t.Errorf("method = %v, esperado %q", entry["method"], http.MethodGet)
	}
	if entry["path"] != "/qualquer-coisa" {
		t.Errorf("path = %v, esperado %q", entry["path"], "/qualquer-coisa")
	}
	if entry["status"] != float64(http.StatusTeapot) {
		t.Errorf("status = %v, esperado %d", entry["status"], http.StatusTeapot)
	}
}

func TestServer_Run_GracefulShutdown(t *testing.T) {
	testLogger := slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))

	// ":0" pede pro sistema operacional escolher uma porta livre — evita
	// conflito se outra coisa já estiver usando uma porta fixa na sua máquina.
	srv := New(":0", testLogger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- srv.Run(ctx)
	}()

	time.Sleep(50 * time.Millisecond) // dá tempo do ListenAndServe subir antes de cancelar
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() retornou erro inesperado no shutdown: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run() não retornou depois do shutdown")
	}
}
