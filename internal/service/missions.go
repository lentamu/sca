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

type CreateMissionInput struct {
	CatId   uuid.UUID
	Targets []CreateTargetInput
}

type AssignCatInput struct {
	MissionId uuid.UUID
	CatId     uuid.UUID
}

type AddTargetInput struct {
	MissionId uuid.UUID
	TargetId  uuid.UUID
}

type MissionService interface {
	Create(ctx context.Context, input CreateMissionInput) (*models.Mission, error)
	ById(ctx context.Context, id uuid.UUID) (*models.Mission, error)
	All(ctx context.Context) ([]*models.Mission, error)
	Delete(ctx context.Context, id uuid.UUID) error
	MarkComplete(ctx context.Context, id uuid.UUID) error
	AssignCat(ctx context.Context, input AssignCatInput) error
	AddTarget(ctx context.Context, input AddTargetInput) error
}

type MissionServiceImpl struct {
	store       storage.MissionStorage
	catStore    storage.CatStorage
	targetStore storage.TargetStorage
	cache       cache.Cache
}

func NewMissionService(store storage.MissionStorage, catStore storage.CatStorage, targetStore storage.TargetStorage, cache cache.Cache) *MissionServiceImpl {
	return &MissionServiceImpl{
		store:       store,
		catStore:    catStore,
		targetStore: targetStore,
		cache:       cache,
	}
}

func (s *MissionServiceImpl) Create(ctx context.Context, input CreateMissionInput) (*models.Mission, error) {
	const cacheKey = "missions"

	var cat *models.Cat
	var err error
	var catIdPtr *uuid.UUID
	if input.CatId != uuid.Nil {
		cat, err = s.catStore.ById(ctx, input.CatId)
		if err != nil {
			return nil, err
		}
		catIdPtr = &cat.ID
	}

	mission := &models.Mission{
		ID:       uuid.New(),
		CatId:    catIdPtr,
		Complete: false,
	}

	targets := make([]*models.Target, len(input.Targets))
	for i, t := range input.Targets {
		targets[i] = &models.Target{
			ID:        uuid.New(),
			Name:      t.Name,
			Country:   t.Country,
			Notes:     t.Notes,
			Complete:  false,
			MissionID: &mission.ID,
		}
	}
	mission.Targets = targets

	err = s.store.Create(ctx, mission, targets)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return mission, nil
}

func (s *MissionServiceImpl) ById(ctx context.Context, id uuid.UUID) (*models.Mission, error) {
	mission, err := s.store.ById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mission, nil
}

func (s *MissionServiceImpl) All(ctx context.Context) ([]*models.Mission, error) {
	const cacheKey = "missions"

	if items, _ := s.cache.Get(ctx, cacheKey); items != nil {
		return items.([]*models.Mission), nil
	}

	missions, err := s.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, missions, time.Minute*10)

	return missions, nil
}

func (s *MissionServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const cacheKey = "missions"

	mission, err := s.ById(ctx, id)
	if err != nil {
		return err
	}
	if mission.CatId != nil && *mission.CatId != uuid.Nil {
		return errors.ErrConflict{Msg: "Cannot delete mission: cat is assigned"}
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}

func (s *MissionServiceImpl) MarkComplete(ctx context.Context, id uuid.UUID) error {
	const cacheKey = "missions"

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

func (s *MissionServiceImpl) AssignCat(ctx context.Context, input AssignCatInput) error {
	const cacheKey = "missions"

	mission, err := s.ById(ctx, input.MissionId)
	if err != nil {
		return err
	}
	if mission.Complete {
		return errors.ErrConflict{Msg: "Cannot assign cat: mission is completed"}
	}

	_, err = s.catStore.ById(ctx, input.CatId)
	if err != nil {
		return err
	}

	err = s.store.AssignCat(ctx, input.MissionId, input.CatId)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}

func (s *MissionServiceImpl) AddTarget(ctx context.Context, input AddTargetInput) error {
	const cacheKey = "missions"

	mission, err := s.ById(ctx, input.MissionId)
	if err != nil {
		return err
	}
	if mission.Complete {
		return errors.ErrConflict{Msg: "Cannot add target: mission is completed"}
	}
	if len(mission.Targets) >= 3 {
		return errors.ErrConflict{Msg: "Cannot add target: mission already has 3 targets"}
	}

	target, err := s.targetStore.ById(ctx, input.TargetId)
	if err != nil {
		return err
	}

	err = s.store.AddTarget(ctx, input.MissionId, target)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return nil
}
