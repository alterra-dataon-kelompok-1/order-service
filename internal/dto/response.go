package dto

import (
	"time"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
)

type SearchGetResponse[T any] struct {
	Data           []T `json:"data"`
	PaginationInfo *PaginationInfo
}

type GetOrderItemResponse struct {
	MenuID          uuid.UUID `json:"menu_id"`
	OrderItemStatus string    `json:"item_status"`
	Quantity        int       `json:"quantity"`
	Price           float32   `json:"price"`
}

type CreateOrderResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	OrderStatus   string    `json:"order_status"`
	TotalPrice    int       `json:"total_price"`
	TotalQuantity int       `json:"total_quantity"`

	OrderItems []GetOrderItemResponse `json:"items"`
	CreatedAt  time.Time              `json:"created_at"`
}

type UpdateOrderResponse struct {
	UserID        *uuid.UUID                `json:"user_id;omitempty"`
	OrderStatus   *model.OrderStatus        `json:"order_status;omitempty"`
	TotalPrice    *float32                  `json:"total_price;omitempty"`
	TotalQuantity *float32                  `json:"total_quantity;omitempty"`
	OrderItems    *[]UpdateOrderItemRequest `json:"items;omitempty"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type UpdateOrderItemResponse struct {
	MenuID    uuid.UUID `json:"menu_id" validation:"required"`
	Quantity  int       `json:"quantity"`
	Price     float32   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}
