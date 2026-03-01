package repositorysql

import (
	"context"
	"sync"

	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
)

type ActivityLogRepository interface {
	Append(ctx context.Context, log entity.ActivityLog) error
}

type activityLogRepository struct {
	mu   sync.Mutex
	logs []entity.ActivityLog
}

func NewActivityLogRepository() ActivityLogRepository {
	return &activityLogRepository{logs: []entity.ActivityLog{}}
}

func (s *activityLogRepository) Append(ctx context.Context, log entity.ActivityLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = append(s.logs, log)
	return nil
}
