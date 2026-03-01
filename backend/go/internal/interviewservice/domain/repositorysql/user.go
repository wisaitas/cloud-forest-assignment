package repositorysql

import (
	"context"
	"errors"
	"sync"

	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
	"github.com/wisaitas/cloud-forest-assignment/pkg/bcryptx"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type userRepository struct {
	mu    sync.RWMutex
	users []entity.User
}

func NewUserRepository(bcrypt bcryptx.Bcrypt) (UserRepository, error) {
	hashed, err := bcrypt.GenerateFromPassword("not-so-secure-password", 10)
	if err != nil {
		return nil, err
	}
	return &userRepository{users: []entity.User{
		{
			ID:       "123123123",
			Email:    "john.smith@gmail.com",
			Password: string(hashed),
		},
	}}, nil
}

func (s *userRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return entity.User{}, ErrNotFound
}

func (s *userRepository) CreateUser(ctx context.Context, user entity.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)
	return nil
}
