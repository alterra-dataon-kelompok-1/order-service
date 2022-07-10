package order

import (
	m "github.com/alterra-dataon-kelompok-1/order-service/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo, h Handler) {
	e.GET("/v1/orders", h.Get, m.NewAuthMiddleware(m.StaffAndAdmin, false).Authenticate)
	e.POST("/v1/orders", h.Create)

	e.GET("/v1/orders/:id", h.GetOrderByID)
	e.DELETE("/v1/orders/:id", h.DeleteOrderByID, m.NewAuthMiddleware(m.StaffAndAdmin, false).Authenticate)
	e.PUT("/v1/orders/:id", h.UpdateOrderByID, m.NewAuthMiddleware(m.StaffAndAdmin, false).Authenticate)
}
