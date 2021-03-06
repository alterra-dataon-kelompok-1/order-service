package model

import "github.com/google/uuid"

type OrderItem struct {
	OrderID         uuid.UUID       `json:"order_id" gorm:"type:char(36);primaryKey"`
	MenuID          uuid.UUID       `json:"menu_id" gorm:"primaryKey;not null"`
	OrderItemStatus OrderItemStatus `json:"order_item_status" gorm:"type:enum('pending', 'paid', 'prepared', 'delivered');default:'pending'"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price" gorm:"not null;type:decimal(25,2)"`

	Model
}

// Implemented to apply generic
func (o OrderItem) GetQuantity() int {
	return o.Quantity
}
