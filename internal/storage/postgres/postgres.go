package postgres

import (
	"fmt"

	"github.com/Richtermnd/RichterAuth/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	connString := buildConnectionString()
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}
	return conn
}

func buildConnectionString() string {
	cfg := config.Config().Storage
	template := "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	return fmt.Sprintf(template, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
