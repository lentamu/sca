package errors

import "github.com/gofiber/fiber/v3"

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}

	switch e := err.(type) {
	case ErrNotFound:
		code = fiber.StatusNotFound
		msg = e.Msg
	case ErrConflict:
		code = fiber.StatusConflict
		msg = e.Msg
	}

	return c.Status(code).JSON(&ErrorResponse{
		Error: true,
		Response: ErrorDetail{
			Message:  msg,
			Instance: c.Path(),
		},
	})
}

type ErrorResponse struct {
	Error    bool        `json:"error"`
	Response ErrorDetail `json:"response"`
}

type ErrorDetail struct {
	Message  string `json:"message"`
	Instance string `json:"instance"`
}
