package model

// TODO: Change ID type from int to be string to implement UUID
// TODO: Implement model by ourselved to remove dependency to Gorm
type OrderItem struct {
	OrderID  *uint           `json:"order_id"`
	MenuID   int             `json:"menu_id" gorm:"primaryKey;not null"`
	Status   OrderItemStatus `json:"order_item_status" gorm:"type:enum('pending', 'paid', 'prepared', 'delivered');default:'pending'"`
	Quantity int             `json:"quantity"`
	Price    float32         `json:"price"`

	Model
}
