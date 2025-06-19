package mysql

import (
	"context"
	"database/sql"
	stderrors "errors"

	"sca/internal/models"
	"sca/pkg/database/mysql"
	"sca/pkg/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrCatAlreadyExists = errors.ErrConflict{Msg: "Cat is already exists"}
	ErrCatNotFound      = errors.ErrNotFound{Msg: "Cat not found"}
)

type CatStorage struct {
	db *sqlx.DB
}

func NewCatStorage(db *sqlx.DB) *CatStorage {
	return &CatStorage{db: db}
}

func (s *CatStorage) Create(ctx context.Context, cat *models.Cat) error {
	query := `INSERT INTO cats (id, name, years_of_experience, breed, salary) VALUES (:id, :name, :years_of_experience, :breed, :salary)`
	_, err := s.db.NamedExecContext(ctx, query, cat)
	if err != nil {
		if mysql.IsDuplicate(err) {
			return ErrCatAlreadyExists
		}
		return err
	}
	return nil
}

func (s *CatStorage) ById(ctx context.Context, id uuid.UUID) (*models.Cat, error) {
	query := `SELECT * FROM cats WHERE id = ?`
	var cat models.Cat
	err := s.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return nil, ErrCatNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (s *CatStorage) All(ctx context.Context) ([]*models.Cat, error) {
	query := `SELECT * FROM cats`
	cats := []*models.Cat{}
	err := s.db.SelectContext(ctx, &cats, query)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *CatStorage) Update(ctx context.Context, cat *models.Cat) error {
	query := `UPDATE cats SET salary = :salary WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, cat)
	if err != nil {
		return err
	}
	return nil
}

func (s *CatStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cats WHERE id = ?`
	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
