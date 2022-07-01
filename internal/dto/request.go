package dto

type CreateOrderRequest struct {
	UserID     int                      `json:"user_id"`
	OrderItems []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	MenuID   int `json:"menu_id"`
	Quantity int `json:"quantity"`
}
