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
	MenuID   uuid.UUID `json:"menu_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required"`
}

type UpdateOrderRequest struct {
	OrderStatus *model.OrderStatus        `json:"order_status,omitempty"`
	OrderItems  *[]UpdateOrderItemRequest `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

type UpdateOrderItemRequest struct {
	OrderID  uuid.UUID              `json:"order_id" validate:"omitempty,uuid" gorm:"type:uuid;primaryKey"`
	MenuID   uuid.UUID              `json:"menu_id" gorm:"primaryKey;not null"`
	Status   *model.OrderItemStatus `json:"order_item_status"`
	Quantity *int                   `json:"quantity"`
}

type ByIDRequest struct {
	ID uuid.UUID `json:"id" param:"id" validate:"required"`
}

func (c CreateOrderItemRequest) GetQuantity() int {
	return c.Quantity
}
