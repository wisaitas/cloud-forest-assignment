package repositorysql

import (
	"context"
	"errors"
	"sync"

	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
)

type ServerRepository interface {
	Create(ctx context.Context, server entity.Server) error
	GetByID(ctx context.Context, id string) (entity.Server, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.Server, error)
	UpdatePowerStatus(ctx context.Context, id, powerStatus string) error
}

type serverRepository struct {
	mu      sync.RWMutex
	servers []entity.Server
}

func NewServerRepository() ServerRepository {
	return &serverRepository{servers: []entity.Server{}}
}

func (s *serverRepository) Create(ctx context.Context, server entity.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers = append(s.servers, server)
	return nil
}

func (s *serverRepository) GetByID(ctx context.Context, id string) (entity.Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sv := range s.servers {
		if sv.ID == id {
			return sv, nil
		}
	}
	return entity.Server{}, errors.New("server not found")
}

func (s *serverRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []entity.Server
	for _, server := range s.servers {
		if server.UserID == userID {
			out = append(out, server)
		}
	}
	return out, nil
}

func (s *serverRepository) UpdatePowerStatus(ctx context.Context, id, powerStatus string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.servers {
		if s.servers[i].ID == id {
			s.servers[i].PowerStatus = powerStatus
			return nil
		}
	}
	return nil
}
