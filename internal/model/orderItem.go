package model

type OrderItem struct {
	OrderID           int `gorm:"primaryKey"`
	MenuID            int `gorm:"primaryKey;not null"`
	OrderItemStatusId int `gorm:"size:2"`
	Quantity          int
	Price             float32

	Model
}
