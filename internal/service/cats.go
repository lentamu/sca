package service

import (
	"context"
	"time"

	"sca/internal/models"
	"sca/internal/storage"
	"sca/pkg/cache"

	"github.com/google/uuid"
)

type CreateCatInput struct {
	Name              string
	YearsOfExperience int
	Breed             string
	Salary            float64
}

type CatService interface {
	Create(ctx context.Context, input CreateCatInput) (*models.Cat, error)
	ById(ctx context.Context, id uuid.UUID) (*models.Cat, error)
	All(ctx context.Context) ([]*models.Cat, error)
	Update(ctx context.Context, id uuid.UUID, salayry float64) (*models.Cat, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CatServiceImpl struct {
	store storage.CatStorage
	cache cache.Cache
}

func NewCatService(store storage.CatStorage, cache cache.Cache) *CatServiceImpl {
	return &CatServiceImpl{
		store: store,
		cache: cache,
	}
}

func (s *CatServiceImpl) Create(ctx context.Context, input CreateCatInput) (*models.Cat, error) {
	const cacheKey = "cats"

	cat := &models.Cat{
		ID:                uuid.New(),
		Name:              input.Name,
		YearsOfExperience: input.YearsOfExperience,
		Breed:             input.Breed,
		Salary:            input.Salary,
	}
	err := s.store.Create(ctx, cat)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return cat, nil
}

func (s *CatServiceImpl) ById(ctx context.Context, id uuid.UUID) (*models.Cat, error) {
	cat, err := s.store.ById(ctx, id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CatServiceImpl) All(ctx context.Context) ([]*models.Cat, error) {
	const cacheKey = "cats"

	if items, _ := s.cache.Get(ctx, cacheKey); items != nil {
		return items.([]*models.Cat), nil
	}

	cats, err := s.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, cats, time.Minute*10)

	return cats, nil
}

func (s *CatServiceImpl) Update(ctx context.Context, id uuid.UUID, salary float64) (*models.Cat, error) {
	const cacheKey = "cats"

	cat, err := s.ById(ctx, id)
	if err != nil {
		return nil, err
	}

	cat.Salary = salary

	err = s.store.Update(ctx, cat)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Del(ctx, cacheKey)

	return cat, nil
}

func (s *CatServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const cacheKey = "cats"

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
