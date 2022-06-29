package order

import "github.com/labstack/echo/v4"

func RegisterHandlers(e *echo.Echo, h Handler) {
	e.GET("/orders", h.Get)
	// e.POST("/books", h.Create)
	//
	// e.GET("/books/:id", h.GetByID)
	// e.DELETE("/books/:id", h.Delete)
	// e.PUT("/books/:id", h.Update)
}
