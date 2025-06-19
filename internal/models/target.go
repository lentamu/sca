package models

import (
	"github.com/google/uuid"
)

type Target struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Country   string     `json:"country"`
	Notes     string     `json:"notes"`
	Complete  bool       `json:"complete"`
	MissionID *uuid.UUID `json:"mission_id" db:"mission_id"`
}
