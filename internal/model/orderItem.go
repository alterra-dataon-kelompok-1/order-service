package model

import "github.com/google/uuid"

// TODO: Change ID type from int to be string to implement UUID
type OrderItem struct {
	OrderID  uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	MenuID   int             `json:"menu_id" gorm:"primaryKey;not null"`
	Status   OrderItemStatus `json:"order_item_status" gorm:"type:enum('pending', 'paid', 'prepared', 'delivered');default:'pending'"`
	Quantity int             `json:"quantity"`
	Price    float32         `json:"price"`

	Model
}

func (o OrderItem) GetQuantity() int {
	return o.Quantity
}
