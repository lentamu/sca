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

var ErrTargetNotFound = errors.ErrNotFound{Msg: "Target not found"}

type TargetStorage struct {
	db *sqlx.DB
}

func NewTargetStorage(db *sqlx.DB) *TargetStorage {
	return &TargetStorage{db: db}
}

func (s *TargetStorage) Create(ctx context.Context, target *models.Target) error {
	query := `INSERT INTO targets (id, name, country, notes, complete) VALUES (:id, :name, :country, :notes, :complete)`
	_, err := s.db.NamedExecContext(ctx, query, target)
	if err != nil {
		return err
	}
	return nil
}

func (s *TargetStorage) ById(ctx context.Context, id uuid.UUID) (*models.Target, error) {
	query := `SELECT * FROM targets WHERE id = ?`
	var target models.Target
	err := s.db.GetContext(ctx, &target, query, id)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return &target, nil
}

func (s *TargetStorage) All(ctx context.Context) ([]*models.Target, error) {
	query := `SELECT * FROM targets`
	targets := []*models.Target{}
	err := s.db.SelectContext(ctx, &targets, query)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (s *TargetStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM targets WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TargetStorage) MarkComplete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE targets SET complete = true WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TargetStorage) UpdateNotes(ctx context.Context, id uuid.UUID, notes string) error {
	query := `UPDATE targets SET notes = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, notes, id)
	if err != nil {
		return err
	}
	return nil
}
