package model

import "github.com/google/uuid"

type Order struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        *uuid.UUID  `json:"user_id"`
	OrderStatus   OrderStatus `json:"order_status" gorm:"type:enum('pending','paid','prepared','completed','canceled');default:'pending'"`
	TotalPrice    float32     `json:"total_price" gorm:"not null"`
	TotalQuantity int         `json:"total_quantity" gorm:"not null"`

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`

	Model
}
