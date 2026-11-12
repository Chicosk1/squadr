package api

import (
	"net/http"
	"runtime/debug"
	"time"
)

// gravadorStatus guarda o status escrito pelo handler, que o http.ResponseWriter não expõe.
type gravadorStatus struct {
	http.ResponseWriter
	status int
}

func (g *gravadorStatus) WriteHeader(status int) {
	if g.status == 0 {
		g.status = status
	}
	g.ResponseWriter.WriteHeader(status)
}

func (g *gravadorStatus) Write(b []byte) (int, error) {
	if g.status == 0 {
		g.status = http.StatusOK // Write sem WriteHeader responde 200
	}
	return g.ResponseWriter.Write(b)
}

// registrarLog escreve uma linha por requisição com método, caminho, status e latência.
// A latência inclui banco e scoring; a Fase 7 mede os dois separados.
func (s *Servidor) registrarLog(prox http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &gravadorStatus{ResponseWriter: w}
		prox.ServeHTTP(g, r)
		if g.status == 0 {
			g.status = http.StatusOK
		}
		s.log.Printf("%s %s %d %s", r.Method, r.URL.RequestURI(), g.status, time.Since(inicio).Round(time.Microsecond))
	})
}

// recuperar transforma um panic em 500, registra a pilha e mantém o servidor de pé.
func (s *Servidor) recuperar(prox http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler { // aborto proposital: o net/http trata
					panic(v)
				}
				s.log.Printf("panic em %s %s: %v\n%s", r.Method, r.URL.Path, v, debug.Stack())
				responderErro(w, http.StatusInternalServerError, "erro interno")
			}
		}()
		prox.ServeHTTP(w, r)
	})
}
