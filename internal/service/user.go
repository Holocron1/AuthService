package service

import (
	"context"
	"errors"
	"github.com/Holocron1/authservice/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserServiceImpl struct {
	UserStore UserStore
	JWTSecret string
	JWTTTL    time.Duration
}

func NewUserServiceImpl(store UserStore, JWTSecret string, JWTTTL time.Duration) *UserServiceImpl {
	return &UserServiceImpl{UserStore: store, JWTSecret: JWTSecret, JWTTTL: JWTTTL}
}

type UserStore interface {
	Get(ctx context.Context, username string) (domain.User, error)
}

var ErrUserNotFound = errors.New("user not found")

func (u *UserServiceImpl) ValidateCredentials(ctx context.Context, username, password string) (bool, error) {
	user, err := u.UserStore.Get(ctx, username)
	if err != nil {
		return false, ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, errors.New("Wrong password")
	}
	return true, nil
}

const issuer = "auth-service"

func (u *UserServiceImpl) GenerateToken(ctx context.Context, username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   username,
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.JWTTTL)),
		NotBefore: nil,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.JWTSecret))
}

func (u *UserServiceImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	newToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(u.JWTSecret), nil })
	if newToken == nil || err != nil {
		return "", errors.New("Invalid token")
	}

	username, _ := newToken.Claims.GetSubject()
	return u.GenerateToken(ctx, username)
}
