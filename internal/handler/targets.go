package handler

import (
	"sca/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type TargetHandler struct {
	service service.TargetService
}

func NewTargetHandler(service service.TargetService) *TargetHandler {
	return &TargetHandler{service: service}
}

func (h *TargetHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/targets", h.Create)
	router.Get("/targets/:id", h.ById)
	router.Get("/targets", h.List)
	router.Patch("/targets/:id/notes", h.UpdateNotes)
	router.Delete("/targets/:id", h.Delete)
	router.Post("/targets/:id/complete", h.MarkComplete)
}

func (h *TargetHandler) Create(c fiber.Ctx) error {
	var req struct {
		Name    string `json:"name" validate:"required,min=3,max=32"`
		Country string `json:"country" validate:"required,min=3,max=32"`
		Notes   string `json:"notes" validate:"required,min=3,max=255"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	target, err := h.service.Create(c.Context(), service.CreateTargetInput{
		Name:    req.Name,
		Country: req.Country,
		Notes:   req.Notes,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&target)
}

func (h *TargetHandler) ById(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	target, err := h.service.ById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&target)
}

func (h *TargetHandler) List(c fiber.Ctx) error {
	targets, err := h.service.All(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&targets)
}

func (h *TargetHandler) Delete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Target deleted successfully"})
}

func (h *TargetHandler) MarkComplete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	err = h.service.MarkComplete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Target marked as complete"})
}

func (h *TargetHandler) UpdateNotes(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	var req struct {
		Notes string `json:"notes" validate:"required,min=3,max=255"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	err = h.service.UpdateNotes(c.Context(), service.UpdateNotesInput{
		ID:    id,
		Notes: req.Notes,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Target notes updated successfully"})
}
