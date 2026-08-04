package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	httpServer *http.Server
	router     chi.Router
	logger     *slog.Logger
}

func New(addr string, appLogger *slog.Logger) *Server {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(requestLogger(appLogger))
	router.Use(middleware.Recoverer)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		router: router,
		logger: appLogger,
	}
}

// Router expõe o roteador pra quem monta o processo (main.go) registrar rotas.
func (s *Server) Router() chi.Router {
	return s.router
}

// Run sobe o servidor e bloqueia até ctx ser cancelado. Quando isso
// acontece, para de aceitar conexão nova e espera as em andamento
// terminarem (até 10s) antes de retornar — é o "graceful" do shutdown.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("servidor http iniciado", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.logger.Info("desligando servidor http")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	}
}

func requestLogger(appLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			appLogger.Info("requisição",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
