package handler

import (
	"sca/internal/service"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	cats     *CatHandler
	missions *MissionHandler
	targets  *TargetHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		cats:     NewCatHandler(service.Cats),
		missions: NewMissionHandler(service.Missions),
		targets:  NewTargetHandler(service.Targets),
	}
}

func (s *Handler) RegisterRoutes(router fiber.Router) {
	s.cats.RegisterRoutes(router)
	s.missions.RegisterRoutes(router)
	s.targets.RegisterRoutes(router)
}
