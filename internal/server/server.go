package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Richtermnd/RichterAuth/internal/config"
	"github.com/Richtermnd/RichterAuth/internal/server/auth"
	"github.com/Richtermnd/RichterAuth/internal/server/middlewares"
	"github.com/go-pkgz/routegroup"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
	mux    *routegroup.Bundle
}

func NewServer(
	log *slog.Logger,
	userService auth.UserService,
) *Server {
	// create router and middlewares
	httpMux := http.NewServeMux()
	mux := routegroup.New(httpMux)
	mux.Use(middlewares.LoggingMiddleware(log))

	// register routes
	auth.Register(mux, userService)

	// Create server
	return &Server{
		log:    log,
		server: createHttpServer(httpMux),
		mux:    mux,
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func createHttpServer(mux http.Handler) *http.Server {
	cfg := config.Config().Server
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		Handler:      mux,
	}
}
