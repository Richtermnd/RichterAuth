package auth

import (
	"context"

	"github.com/Richtermnd/RichterAuth/internal/domain/requests"
	"github.com/Richtermnd/RichterAuth/internal/domain/responses"
	"github.com/go-pkgz/routegroup"
)

type UserService interface {
	Register(ctx context.Context, user requests.Register) error
	Login(ctx context.Context, login requests.Login) (responses.Token, error)
	Confirm(ctx context.Context, confirm requests.Confirm) error
}

type Router struct {
	mux         *routegroup.Bundle
	userService UserService
}

func Register(mux *routegroup.Bundle, userService UserService) {
	mux = mux.Mount("/auth")
	r := &Router{mux: mux, userService: userService}
	r.mux.HandleFunc("POST /register", r.Register)
	r.mux.HandleFunc("POST /login", r.Login)
	r.mux.HandleFunc("GET /{id}/confirm/{key}", r.Confirm)
}
