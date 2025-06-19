package mysql

import (
	"context"
	"database/sql"
	stderrors "errors"

	"sca/internal/models"
	"sca/pkg/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrMissionNotFound = errors.ErrNotFound{Msg: "Mission not found"}

type MissionStorage struct {
	db *sqlx.DB
}

func NewMissionStorage(db *sqlx.DB) *MissionStorage {
	return &MissionStorage{db: db}
}

func (s *MissionStorage) Create(ctx context.Context, mission *models.Mission, targets []*models.Target) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	queryMission := `INSERT INTO missions (id, cat_id, complete) VALUES (:id, :cat_id, :complete)`
	_, err = tx.NamedExecContext(ctx, queryMission, mission)
	if err != nil {
		return err
	}

	queryTarget := `INSERT INTO targets (id, name, country, notes, complete, mission_id) VALUES (:id, :name, :country, :notes, :complete, :mission_id)`
	for _, t := range targets {
		_, err = tx.NamedExecContext(ctx, queryTarget, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MissionStorage) ById(ctx context.Context, id uuid.UUID) (*models.Mission, error) {
	query := `SELECT * FROM missions WHERE id = ?`
	var mission models.Mission
	err := s.db.GetContext(ctx, &mission, query, id)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return nil, ErrMissionNotFound
		}
		return nil, err
	}

	targetsQuery := `SELECT * FROM targets WHERE mission_id = ?`
	var targets []*models.Target
	err = s.db.SelectContext(ctx, &targets, targetsQuery, mission.ID)
	if err != nil {
		return nil, err
	}
	mission.Targets = targets

	return &mission, nil
}

func (s *MissionStorage) All(ctx context.Context) ([]*models.Mission, error) {
	query := `SELECT * FROM missions`
	missions := []*models.Mission{}
	err := s.db.SelectContext(ctx, &missions, query)
	if err != nil {
		return nil, err
	}

	for _, mission := range missions {
		targetsQuery := `SELECT * FROM targets WHERE mission_id = ?`
		var targets []*models.Target
		err = s.db.SelectContext(ctx, &targets, targetsQuery, mission.ID)
		if err != nil {
			return nil, err
		}
		mission.Targets = targets
	}

	return missions, nil
}

func (s *MissionStorage) Update(ctx context.Context, mission *models.Mission) error {
	query := `UPDATE missions SET complete = :complete WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, mission)
	if err != nil {
		return err
	}
	return nil
}

func (s *MissionStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM missions WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *MissionStorage) AssignCat(ctx context.Context, missionId, catId uuid.UUID) error {
	query := `UPDATE missions SET cat_id = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, catId, missionId)
	if err != nil {
		return err
	}
	return nil
}

func (s *MissionStorage) AddTarget(ctx context.Context, missionId uuid.UUID, target *models.Target) error {
	query := `UPDATE targets SET mission_id = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, missionId, target.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MissionStorage) MarkComplete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE missions SET complete = true WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
