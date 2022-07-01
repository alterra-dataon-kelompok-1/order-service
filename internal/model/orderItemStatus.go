package model

type OrderItemStatus struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	State string `json:"state"`

	Model
}
