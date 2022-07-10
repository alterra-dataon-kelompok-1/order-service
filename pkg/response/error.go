package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Response `json:"response"`
	HTTPCode int
}

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error const
type errorConst struct {
	Duplicate           Error
	ResourceNotFound    Error
	RouteNotFound       Error
	BadRequest          Error
	Unauthorized        Error
	Validation          Error
	InternalServerError Error

	NotEnoughStock    Error
	NoOrderItem       Error
	CannotCancelOrder Error
}

var ErrorConst errorConst = errorConst{
	Duplicate: Error{
		Response: Response{
			Code:    "E_DUPLICATE",
			Message: "duplicate",
		},
		HTTPCode: http.StatusConflict,
	},
	ResourceNotFound: Error{
		Response: Response{
			Code:    "E_NOT_FOUND",
			Message: "resource not found",
		},
		HTTPCode: http.StatusNotFound,
	},
	BadRequest: Error{
		Response: Response{
			Code:    "E_BAD_REQUEST",
			Message: "bad request",
		},
		HTTPCode: http.StatusBadRequest,
	},
	Unauthorized: Error{
		Response: Response{
			Code:    "E_UNAUTHORIZED",
			Message: "unauthorized",
		},
		HTTPCode: http.StatusUnauthorized,
	},
	Validation: Error{
		Response: Response{
			Code:    "E_VAL",
			Message: "request format is not supported",
		},
		HTTPCode: http.StatusBadRequest,
	},
	InternalServerError: Error{
		Response: Response{
			Code:    "E_SERVER",
			Message: "internal server error",
		},
		HTTPCode: http.StatusInternalServerError,
	},
	NotEnoughStock: Error{
		Response: Response{
			Code:    "E_NO_STOCK",
			Message: "required resource cannot satisfy request",
		},
		HTTPCode: http.StatusBadRequest,
	},
	NoOrderItem: Error{
		Response: Response{
			Code:    "E_NO_ITEM_PROVIDED",
			Message: "order item quantity must not be zero",
		},
		HTTPCode: http.StatusBadRequest,
	},
	CannotCancelOrder: Error{
		Response: Response{
			Code:    "E_CANNOT_CANCEL",
			Message: "prepared order cannot be canceled",
		},
		HTTPCode: http.StatusBadRequest,
	},
}

func NewErrorResponse(c echo.Context, error Error) error {
	return c.JSON(error.HTTPCode, error.Response)
}
