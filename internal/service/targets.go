package service

import (
	"context"
	"time"

	"sca/internal/models"
	"sca/internal/storage"
	"sca/pkg/cache"
	"sca/pkg/errors"

	"github.com/google/uuid"
)

type CreateTargetInput struct {
	Name    string
	Country string
	Notes   string
}

type UpdateNotesInput struct {
	ID    uuid.UUID
	Notes string
}

type TargetService interface {
	Create(ctx context.Context, input CreateTargetInput) (*models.Target, error)
	ById(ctx context.Context, id uuid.UUID) (*models.Target, error)
	All(ctx context.Context) ([]*models.Target, error)
	Delete(ctx context.Context, id uuid.UUID) error
	MarkComplete(ctx context.Context, id uuid.UUID) error
	UpdateNotes(ctx context.Context, input UpdateNotesInput) error
}

type TargetServiceImpl struct {
	store        storage.TargetStorage
	missionStore storage.MissionStorage
	cache        cache.Cache
}

func NewTargetService(store storage.TargetStorage, missionStore storage.MissionStorage, cache cache.Cache) *TargetServiceImpl {
	return &TargetServiceImpl{
		store:        store,
		missionStore: missionStore,
		cache:        cache,
	}
}

func (s *TargetServiceImpl) Create(ctx context.Context, input CreateTargetInput) (*models.Target, error) {
	const cacheKey = "targets"

	target := &models.Target{
		ID:      uuid.New(),
		Name:    input.Name,
		Country: input.Country,
		Notes:   input.Notes,
	}

	err := s.store.Create(ctx, target)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return target, nil
}

func (s *TargetServiceImpl) ById(ctx context.Context, id uuid.UUID) (*models.Target, error) {
	target, err := s.store.ById(ctx, id)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (s *TargetServiceImpl) All(ctx context.Context) ([]*models.Target, error) {
	const cacheKey = "targets"

	if items, _ := s.cache.Get(ctx, cacheKey); items != nil {
		return items.([]*models.Target), nil
	}

	targets, err := s.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, targets, time.Minute*10)

	return targets, nil
}

func (s *TargetServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const cacheKey = "targets"

	_, err := s.ById(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}

func (s *TargetServiceImpl) MarkComplete(ctx context.Context, id uuid.UUID) error {
	const cacheKey = "targets"

	_, err := s.ById(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.MarkComplete(ctx, id)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}

func (s *TargetServiceImpl) UpdateNotes(ctx context.Context, input UpdateNotesInput) error {
	const cacheKey = "targets"

	target, err := s.ById(ctx, input.ID)
	if err != nil {
		return err
	}
	if target.Complete {
		return errors.ErrConflict{Msg: "Cannot update notes: target is completed"}
	}

	if *target.MissionID != uuid.Nil {
		mission, err := s.missionStore.ById(ctx, *target.MissionID)
		if err != nil {
			return err
		}
		if mission.Complete {
			return errors.ErrConflict{Msg: "Cannot update notes: mission is completed"}
		}
	}

	err = s.store.UpdateNotes(ctx, input.ID, input.Notes)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}
