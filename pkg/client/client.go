package client

import (
	"encoding/json"
	"fmt"
	"sca/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
)

func FetchCatBreeds(url string) ([]models.Breed, error) {
	cc := client.New()
	resp, err := cc.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != fiber.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status())
	}
	var breeds []models.Breed
	if err := json.Unmarshal(resp.Body(), &breeds); err != nil {
		return nil, err
	}
	result := make([]models.Breed, 0, len(breeds))
	for _, b := range breeds {
		result = append(result, models.Breed{ID: b.ID, Name: b.Name})
	}
	return result, nil
}
