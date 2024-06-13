package auth

import (
	"net/http"

	"github.com/Richtermnd/RichterAuth/internal/domain/requests"
	"github.com/Richtermnd/RichterAuth/internal/server/utils"
	"github.com/Richtermnd/goreq"
)

func (m *Router) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest requests.Register
	goreq.Decode(r, &registerRequest)

	err := m.userService.Register(r.Context(), registerRequest)
	if err != nil {
		utils.Encode(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Router) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest requests.Login
	goreq.Decode(r, &loginRequest)

	token, err := m.userService.Login(r.Context(), loginRequest)
	if err != nil {
		utils.Encode(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, r, token)
}

func (m *Router) Confirm(w http.ResponseWriter, r *http.Request) {
	var confirmRequest requests.Confirm
	goreq.Decode(r, &confirmRequest)
	err := m.userService.Confirm(r.Context(), confirmRequest)
	if err != nil {
		utils.Encode(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
