package store

import (
	"AuthService/configs"
	"AuthService/internal/domain"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PostgresStore struct {
	Pool *pgxpool.Pool
}

func NewPostgresStore(c *configs.Config) *PostgresStore {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, c.DATABASE_URL)
	if err != nil {
		log.Println("database error", err)
	}
	return &PostgresStore{Pool: pool}
}

func (s *PostgresStore) Get(username string) (domain.User, error) {
	user := domain.User{}
	result := s.Pool.QueryRow(context.Background(), "SELECT username, password FROM users WHERE username = $1", username)
	err := result.Scan(&user.Username, &user.Password)

	if err != nil {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}
