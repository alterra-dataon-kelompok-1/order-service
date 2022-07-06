package model

import "github.com/google/uuid"

// TODO: Change ID type from int to be string to implement UUID
// TODO: Implement model by ourselved to remove dependency to Gorm
type Order struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        int         `json:"user_id" gorm:"not null"`
	Status        OrderStatus `json:"order_status" gorm:"type:enum('pending','paid','prepared','completed','canceled');default:'pending'"`
	TotalPrice    float32     `json:"total_price" gorm:"not null"`
	TotalQuantity int         `json:"total_quantity" gorm:"not null"`
	DailyID       uint        `json:"daily_id" gorm:"not null"`

	// Day 		Datetime

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`

	Model
}
