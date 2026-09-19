package store

import (
	"errors"
	"github.com/Holocron1/authservice/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryStore struct {
	Users map[string]domain.User
}

func NewInMemoryStore() *InMemoryStore {
	userMap := make(map[string]domain.User)
	store := &InMemoryStore{userMap}
	hash, err := bcrypt.GenerateFromPassword([]byte("secret password"), 10)
	if err != nil {
		panic(err)
	}

	store.Users["testguy1"] = domain.User{Username: "testGuy1", Password: string(hash)}
	store.Users["testguy2"] = domain.User{Username: "testGuy2", Password: string(hash)}
	return store
}

func (s *InMemoryStore) Get(username string) (domain.User, error) {
	user, ok := s.Users[username]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}
