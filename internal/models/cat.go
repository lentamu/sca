package models

import (
	"github.com/google/uuid"
)

type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Cat struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	YearsOfExperience int       `json:"years_of_experience" db:"years_of_experience"`
	Breed             string    `json:"breed"`
	Salary            float64   `json:"salary"`
}
