package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service}
}

func (h *handler) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, "success")
}
