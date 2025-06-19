package handler

import (
	"sca/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MissionHandler struct {
	service service.MissionService
}

func NewMissionHandler(service service.MissionService) *MissionHandler {
	return &MissionHandler{service: service}
}

func (h *MissionHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/missions", h.Create)
	router.Get("/missions/:id", h.ById)
	router.Get("/missions", h.List)
	router.Delete("/missions/:id", h.Delete)
	router.Post("/missions/assign-cat", h.AssignCat)
	router.Post("/missions/:id/complete", h.MarkComplete)
	router.Post("/missions/:id/targets", h.AddTarget)
}

func (h *MissionHandler) Create(c fiber.Ctx) error {
	var req struct {
		CatId   uuid.UUID `json:"cat_id" validate:"omitempty,uuid"`
		Targets []struct {
			Name    string `json:"name" validate:"required,min=3,max=32"`
			Country string `json:"country" validate:"required,min=3,max=32"`
			Notes   string `json:"notes" validate:"required,min=3,max=255"`
		} `json:"targets" validate:"omitempty,min=1,max=3,dive"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	inputTargets := make([]service.CreateTargetInput, len(req.Targets))
	for i, t := range req.Targets {
		inputTargets[i] = service.CreateTargetInput{
			Name:    t.Name,
			Country: t.Country,
			Notes:   t.Notes,
		}
	}

	mission, err := h.service.Create(c.Context(), service.CreateMissionInput{
		CatId:   req.CatId,
		Targets: inputTargets,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&mission)
}

func (h *MissionHandler) ById(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	mission, err := h.service.ById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&mission)
}

func (h *MissionHandler) List(c fiber.Ctx) error {
	missions, err := h.service.All(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&missions)
}

func (h *MissionHandler) Delete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Mission deleted successfully"})
}

func (h *MissionHandler) MarkComplete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	err = h.service.MarkComplete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Mission marked as complete"})
}

func (h *MissionHandler) AssignCat(c fiber.Ctx) error {
	var req struct {
		MissionId uuid.UUID `json:"mission_id" validate:"required,uuid"`
		CatId     uuid.UUID `json:"cat_id" validate:"required,uuid"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	err := h.service.AssignCat(c.Context(), service.AssignCatInput{
		MissionId: req.MissionId,
		CatId:     req.CatId,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Cat assigned to mission"})
}

func (h *MissionHandler) AddTarget(c fiber.Ctx) error {
	var req struct {
		MissionId uuid.UUID `json:"mission_id" validate:"required,uuid"`
		TargetId  uuid.UUID `json:"target_id" validate:"required,uuid"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	err := h.service.AddTarget(c.Context(), service.AddTargetInput{
		MissionId: req.MissionId,
		TargetId:  req.TargetId,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Target added to mission"})
}
