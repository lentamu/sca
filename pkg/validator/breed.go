package validator

import (
	"context"
	"sca/internal/models"
	"sca/pkg/cache"
	"sca/pkg/client"
	"time"

	"github.com/go-playground/validator/v10"
)

var getBreedsFunc func() ([]models.Breed, error)

func InitBreedValidator(
	cache cache.Cache,
	url string,
	cacheKey string,
	ttl time.Duration,
) {
	getBreedsFunc = func() ([]models.Breed, error) {
		ctx := context.Background()
		if val, err := cache.Get(ctx, cacheKey); err == nil && val != nil {
			if breeds, ok := val.([]models.Breed); ok {
				return breeds, nil
			}
		}
		breeds, err := client.FetchCatBreeds(url)
		if err != nil {
			return nil, err
		}
		_ = cache.Set(ctx, cacheKey, breeds, ttl)
		return breeds, nil
	}
}

func breedValidator(fl validator.FieldLevel) bool {
	breed := fl.Field().String()
	breeds, err := getBreedsFunc()
	if err != nil {
		return false
	}
	for _, b := range breeds {
		if b.Name == breed {
			return true
		}
	}
	return false
}
