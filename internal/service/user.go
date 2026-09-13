package service

import (
	"AuthService/configs"
	"AuthService/internal/domain"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	ValidateCredentials(username, password string) (bool, error)
	GenerateToken(username string) (string, error)
	RefreshToken(token string) (string, error)
}

type UserServiceImpl struct {
	UserStore UserStore
}

type UserStore interface {
	Get(username string) (domain.User, error)
}

func (u *UserServiceImpl) ValidateCredentials(username, password string) (bool, error) {
	user, err := u.UserStore.Get(username)
	if err != nil {
		return false, errors.New("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, errors.New("Wrong password")
	}
	return true, nil
}

func (u *UserServiceImpl) GenerateToken(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "auth-service",
		Subject:   username,
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(configs.LoadConfig().JWT_TTL)),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.LoadConfig().JWT_SECRET))
}

func (u *UserServiceImpl) RefreshToken(token string) (string, error) {
	newToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(configs.LoadConfig().JWT_SECRET), nil })
	if newToken == nil || err != nil {
		return "", errors.New("Invalid token")
	}

	username, _ := newToken.Claims.GetSubject()
	return u.GenerateToken(username)
}
