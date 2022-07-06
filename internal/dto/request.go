package dto

import "github.com/google/uuid"

type GetRequest struct {
	Pagination Pagination
	AscField   []string `query:"asc_field"`
	DscField   []string `query:"dsc_field"`
}

type CreateOrderRequest struct {
	UserID     int                      `json:"user_id"`
	OrderItems []CreateOrderItemRequest `json:"items" validation:"required"`
}

type CreateOrderItemRequest struct {
	MenuID   int `json:"menu_id"`
	Quantity int `json:"quantity"`
}

func (c CreateOrderItemRequest) GetQuantity() int {
	return c.Quantity
}

type ByIDRequest struct {
	ID uuid.UUID `json:"id" param:"id" validate:"required"`
}
