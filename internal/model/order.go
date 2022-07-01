package model

// TODO: Change ID type from int to be string to implement UUID
// TODO: Implement model by ourselved to remove dependency to Gorm
type Order struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        int         `json:"user_id" gorm:"not null"`
	Status        OrderStatus `json:"order_status" gorm:"type:enum('pending','paid','prepared','completed','canceled');default:'pending'"`
	TotalPrice    float32     `json:"total_price" gorm:"not null"`
	TotalQuantity int         `json:"total_quantity" gorm:"not null"`

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`

	Model
}
