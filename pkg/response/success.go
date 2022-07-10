package response

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Data     interface{} `json:"data"`
	Meta     interface{} `json:"meta,omitempty"`
	Location string      `json:"_location,omitempty"`
}

type Success struct {
	Response SuccessResponse `json:"response"`
	HTTPCode int
}

type successConst struct {
	ResourceCreated Success
}

// NewSuccessResponse create a new Success instance
func NewSuccessBuilder() *Success {
	return &Success{}
}

// Status Builder for Success struct
func (s *Success) Status(httpCode int) *Success {
	s.HTTPCode = httpCode
	return s
}

// Data builder for Success struct
func (s *Success) WithData(data interface{}) *Success {
	s.Response.Data = data
	return s
}

// Resource Location builder for Success struct
func (s *Success) WithResourceLocation(appURL, resourceID string) *Success {
	s.Response.Location = fmt.Sprintf("http://%s/v1/orders/%s", appURL, resourceID)
	return s
}

// SendJSON: Final function to in NewSuccessBuilder to send JSON response
func (s *Success) SendJSON(c echo.Context) error {
	return c.JSON(s.HTTPCode, s.Response)
}
