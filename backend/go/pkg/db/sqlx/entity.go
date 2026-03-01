package sqlx

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Version   int             `gorm:"type:integer;not null;default:0"`
	CreatedAt time.Time       `gorm:"type:timestamp;not null;default:now()"`
	CreatedBy *uuid.UUID      `gorm:"type:uuid"`
	UpdatedAt time.Time       `gorm:"type:timestamp;not null;default:now()"`
	UpdatedBy *uuid.UUID      `gorm:"type:uuid"`
	DeletedAt *gorm.DeletedAt `gorm:"type:timestamp"`
	DeletedBy *uuid.UUID      `gorm:"type:uuid"`
}
