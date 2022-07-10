package dto

import (
	"time"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
)

type SearchGetResponse[T any] struct {
	Data           []T             `json:"data"`
	PaginationInfo *PaginationInfo `json:"meta"`
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
	TotalPrice    *float64                  `json:"total_price;omitempty"`
	TotalQuantity *int                      `json:"total_quantity;omitempty"`
	OrderItems    *[]UpdateOrderItemRequest `json:"items;omitempty"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type GetOrderResponse struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        *uuid.UUID        `json:"user_id"`
	OrderStatus   model.OrderStatus `json:"order_status"`
	TotalPrice    float64           `json:"total_price"`
	TotalQuantity int               `json:"total_quantity"`

	OrderItems []GetOrderItemResponse `json:"order_items" gorm:"foreignKey:OrderID"`
}

type GetOrderItemResponse struct {
	MenuID          uuid.UUID `json:"menu_id"`
	OrderItemStatus string    `json:"item_status"`
	Quantity        int       `json:"quantity"`
	Price           float64   `json:"price"`
}

type UpdateOrderItemResponse struct {
	MenuID    uuid.UUID `json:"menu_id" validation:"required"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}
