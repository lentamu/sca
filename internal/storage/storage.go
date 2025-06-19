package storage

import (
	"context"

	"sca/internal/models"
	"sca/internal/storage/mysql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CatStorage interface {
	Create(ctx context.Context, cat *models.Cat) error
	ById(ctx context.Context, id uuid.UUID) (*models.Cat, error)
	All(ctx context.Context) ([]*models.Cat, error)
	Update(ctx context.Context, cat *models.Cat) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MissionStorage interface {
	Create(ctx context.Context, mission *models.Mission, targets []*models.Target) error
	ById(ctx context.Context, id uuid.UUID) (*models.Mission, error)
	All(ctx context.Context) ([]*models.Mission, error)
	Update(ctx context.Context, mission *models.Mission) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignCat(ctx context.Context, missionId, catId uuid.UUID) error
	AddTarget(ctx context.Context, missionId uuid.UUID, target *models.Target) error
	MarkComplete(ctx context.Context, id uuid.UUID) error
}

type TargetStorage interface {
	Create(ctx context.Context, target *models.Target) error
	ById(ctx context.Context, id uuid.UUID) (*models.Target, error)
	All(ctx context.Context) ([]*models.Target, error)
	Delete(ctx context.Context, id uuid.UUID) error
	MarkComplete(ctx context.Context, id uuid.UUID) error
	UpdateNotes(ctx context.Context, id uuid.UUID, notes string) error
}

type Storage struct {
	CatStorage     CatStorage
	TargetStorage  TargetStorage
	MissionStorage MissionStorage
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		CatStorage:     mysql.NewCatStorage(db),
		TargetStorage:  mysql.NewTargetStorage(db),
		MissionStorage: mysql.NewMissionStorage(db),
	}
}
