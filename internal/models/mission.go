package models

import (
	"github.com/google/uuid"
)

type Mission struct {
	ID       uuid.UUID  `json:"id"`
	Complete bool       `json:"complete"`
	CatId    *uuid.UUID `json:"cat_id" db:"cat_id"`
	Targets  []*Target  `json:"targets"`
}
