package handler

import (
	"sca/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CatHandler struct {
	service service.CatService
}

func NewCatHandler(service service.CatService) *CatHandler {
	return &CatHandler{service: service}
}

func (h *CatHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/cats", h.Create)
	router.Get("/cats/:id", h.ById)
	router.Get("/cats", h.List)
	router.Patch("/cats/:id", h.Update)
	router.Delete("/cats/:id", h.Delete)
}

func (h *CatHandler) Create(c fiber.Ctx) error {
	var req struct {
		Name              string  `json:"name" validate:"required,min=3,max=32"`
		YearsOfExperience int     `json:"years_of_experience" validate:"required,gte=0,lte=10"`
		Breed             string  `json:"breed" validate:"required,breed"`
		Salary            float64 `json:"salary" validate:"required,gt=0,lte=10000"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	cat, err := h.service.Create(c.Context(), service.CreateCatInput{
		Name:              req.Name,
		YearsOfExperience: req.YearsOfExperience,
		Breed:             req.Breed,
		Salary:            req.Salary,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&cat)
}

func (h *CatHandler) ById(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	cat, err := h.service.ById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&cat)
}

func (h *CatHandler) List(c fiber.Ctx) error {
	cats, err := h.service.All(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&cats)
}

func (h *CatHandler) Update(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	var req struct {
		Salary float64 `json:"salary" validate:"required,gt=0,lte=10000"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	cat, err := h.service.Update(c.Context(), id, req.Salary)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&cat)
}

func (h *CatHandler) Delete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return err
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Cat deleted successfully"})
}
