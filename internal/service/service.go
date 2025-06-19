package service

import (
	"sca/internal/storage"
	"sca/pkg/cache"
)

type Depends struct {
	Storage *storage.Storage
	Cache   cache.Cache
}

type Service struct {
	Cats     CatService
	Missions MissionService
	Targets  TargetService
}

func NewService(depends *Depends) *Service {
	return &Service{
		Cats:     NewCatService(depends.Storage.CatStorage, depends.Cache),
		Missions: NewMissionService(depends.Storage.MissionStorage, depends.Storage.CatStorage, depends.Storage.TargetStorage, depends.Cache),
		Targets:  NewTargetService(depends.Storage.TargetStorage, depends.Storage.MissionStorage, depends.Cache),
	}
}
