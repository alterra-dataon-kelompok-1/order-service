package model

// TODO: Change ID type from int to be string to implement UUID
// TODO: Implement model by ourselved to remove dependency to Gorm
type OrderItem struct {
	OrderID           *uint   `json:"order_id"`
	MenuID            int     `json:"menu_id" gorm:"primaryKey;not null"`
	OrderItemStatusId int     `json:"order_item_status_id" gorm:"size:2"`
	Quantity          int     `json:"quantity"`
	Price             float32 `json:"price"`

	Model
}
