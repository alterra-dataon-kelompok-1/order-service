package dto

import (
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
)

type GetRequest struct {
	Pagination Pagination
	AscField   []string `query:"asc_field"`
	DscField   []string `query:"dsc_field"`
}

type CreateOrderRequest struct {
	UserID     *uuid.UUID               `json:"user_id"`
	OrderItems []CreateOrderItemRequest `json:"order_items" validation:"required"`
}

type CreateOrderItemRequest struct {
	MenuID   uuid.UUID `json:"menu_id"`
	Quantity int       `json:"quantity"`
}

func (c CreateOrderItemRequest) GetQuantity() int {
	return c.Quantity
}

type UpdateOrderRequest struct {
	OrderStatus *model.OrderStatus        `json:"order_status"`
	OrderItems  *[]UpdateOrderItemRequest `json:"order_items"`
}

type UpdateOrderItemRequest struct {
	MenuID   uuid.UUID              `json:"menu_id" validation:"required"`
	Status   *model.OrderItemStatus `json:"order_item_status"`
	Quantity *int                   `json:"quantity"`
}

type ByIDRequest struct {
	ID uuid.UUID `json:"id" param:"id" validate:"required"`
}
