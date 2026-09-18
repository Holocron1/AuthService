package handler

import "context"

type UserService interface {
	ValidateCredentials(ctx context.Context, username, password string) (bool, error)
	GenerateToken(ctx context.Context, username string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}
