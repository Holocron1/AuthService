package handler

import (
	"AuthService/internal/service"
	"net/http"
)

type LoginHandler struct {
	service.UserService
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, err := h.ValidateCredentials(username, password)
	if !ok || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := h.GenerateToken(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

type VerifyHandler struct {
	service.UserService
}

func (v *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	hString := r.Header.Get("Authorization")
	if hString == "" || len(hString) <= 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if hString[0:7] == "Bearer " {
		hString = hString[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := v.RefreshToken(hString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
