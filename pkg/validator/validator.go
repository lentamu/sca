package validator

import (
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(v *validator.Validate) {
	_ = v.RegisterValidation("breed", breedValidator)
}
