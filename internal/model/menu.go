package model

type Menu struct {
	Model

	Name        string `json:"name" gorm:"size:200;unique;not null"`
	Description string `json:"description" gorm:"not null"`
	Image_url   string `json:"image_url" gorm:"not null"`
	Price       int    `json:"price" gorm:"not null"`
	In_stock    int    `json:"in_stock" gorm:"not null"`
	Created_at  int
	Updated_at  int
	Deleted_at  int
}
