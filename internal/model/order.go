package model

type Order struct {
	// ID            string `gorm:"primaryKey"`
	ID            int    `gorm:"primaryKey;autoIncrement"`
	UserID        string `gorm:"not null"`
	MenuID        string `gorm:"not null"`
	OrderStatusID string `gorm:"not null"`
	TotalPrice    string `gorm:"not null"`
	TotalQuantity string `gorm:"not null"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`

	Model
}
