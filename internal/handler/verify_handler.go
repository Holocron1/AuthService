package handler

import (
	"errors"
	"net/http"
)

type VerifyHandler struct {
	userService UserService
}

func NewVerifyHandler(userService UserService) *VerifyHandler {
	return &VerifyHandler{userService: userService}
}

func parseToken(w http.ResponseWriter, r *http.Request) (string, error) {
	hString := r.Header.Get("Authorization")
	if hString == "" || len(hString) <= 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("invalid token")
	}
	if hString[0:7] != "Bearer " {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("invalid token")
	}

	return hString[7:], nil
}

func (v *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hashedToken, err := parseToken(w, r)

	if err != nil {
		return
	}

	token, err := v.userService.RefreshToken(ctx, hashedToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
