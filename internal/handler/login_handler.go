package handler

import (
	"net/http"
)

type LoginHandler struct {
	userService UserService
}

func NewLoginHandler(userService UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ok, err := h.userService.ValidateCredentials(ctx, username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := h.userService.GenerateToken(ctx, username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
