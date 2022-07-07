package middleware

import (
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
}
