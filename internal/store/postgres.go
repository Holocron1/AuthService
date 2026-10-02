package store

import (
	"context"
	"errors"
	"github.com/Holocron1/authservice/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	Pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{Pool: pool}
}

func (s *PostgresStore) Get(ctx context.Context, username string) (domain.User, error) {
	user := domain.User{}
	result := s.Pool.QueryRow(ctx, "SELECT username, password FROM users WHERE username = $1", username)
	err := result.Scan(&user.Username, &user.Password)

	if err != nil {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}
