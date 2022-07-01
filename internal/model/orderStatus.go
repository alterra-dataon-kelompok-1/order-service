package model

// TODO: Implement model by ourselved to remove dependency to Gorm
type OrderStatus struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	State string `json:"state"`

	Model
}
