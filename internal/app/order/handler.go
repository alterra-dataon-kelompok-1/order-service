package order

import (
	"net/http"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
	Create(c echo.Context) error
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

func (h *handler) Create(c echo.Context) error {
	payload := new(dto.CreateOrderRequest)
	c.Bind(&payload)

	newOrder, err := h.service.Create(c.Request().Context(), *payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, newOrder)
}
