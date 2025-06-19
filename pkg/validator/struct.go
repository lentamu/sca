package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator(v *validator.Validate) *StructValidator {
	return &StructValidator{validate: v}
}

func (v *StructValidator) Validate(out any) error {
	if err := v.validate.Struct(out); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		for _, vErr := range err.(validator.ValidationErrors) {
			return fiber.NewError(fiber.StatusBadRequest, vErr.Error())
		}
	}
	return nil
}
