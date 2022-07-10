package fetcher

import "github.com/google/uuid"

type httpResp struct {
	Data Menu        `json:"data"`
	Meta interface{} `json:"meta"`
}

type Menu struct {
	ID             uuid.UUID `json:"id" gorm:"primarykey;autoIncrement"`
	MenuKategoriID int       `json:"menu_kategori_id" gorm:"not null"`
	Name           string    `json:"name" gorm:"size:200;unique;not null"`
	Description    string    `json:"description" gorm:"not null"`
	ImageUrl       string    `json:"image_url" gorm:"not null"`
	Price          float64   `json:"price" gorm:"not null"`
	InStock        int64     `json:"in_stock" gorm:"not null"`
}
