package dto

import (
	"time"
)

type SearchGetResponse[T any] struct {
	Data           []T `json:"data"`
	PaginationInfo *PaginationInfo
}

type GetOrderItemResponse struct {
	MenuID          int     `json:"menu_id"`
	OrderItemStatus string  `json:"item_status"`
	Quantity        int     `json:"quantity"`
	Price           float32 `json:"price"`
}

type CreateOrderResponse struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	OrderStatus   string `json:"order_status"`
	TotalPrice    int    `json:"total_price"`
	TotalQuantity int    `json:"total_quantity"`

	OrderItems []GetOrderItemResponse `json:"items"`
	CreatedAt  time.Time              `json:"created_at"`
}
