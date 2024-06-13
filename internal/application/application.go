package application

import (
	"log/slog"

	"github.com/Richtermnd/RichterAuth/internal/server"
	authservice "github.com/Richtermnd/RichterAuth/internal/service/auth"
	"github.com/Richtermnd/RichterAuth/internal/storage/postgres"
	"github.com/Richtermnd/RichterAuth/internal/storage/postgres/userrepo"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	log    *slog.Logger
	server *server.Server
	db     *sqlx.DB
}

func New(log *slog.Logger) *Application {
	db := postgres.Connect()
	userRepo := userrepo.New(db)
	authService := authservice.New(log, userRepo)
	server := server.NewServer(log, authService)
	log.Info("Application initialized")
	return &Application{
		log:    log,
		server: server,
		db:     db,
	}
}

func (a *Application) Run() {
	go a.server.Run()
}

func (a *Application) Stop() {
	a.server.Stop()
}
