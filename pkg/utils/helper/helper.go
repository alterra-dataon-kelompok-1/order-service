package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

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

// TODO: Implement get item price (API call to other table)
func GetItemPrice(id uuid.UUID) float64 {
	return float64(3000)
}

func GetItemPriceImplement(uid uuid.UUID) float64 {
	resp, err := http.Get(fmt.Sprintf("localhost/menu/%s", uid))
	if err != nil {
		log.Println(err)
	}

	menu := Menu{}

	err = json.NewDecoder(resp.Body).Decode(&menu)
	// menu.Price, menu.InStock

	return float64(3000)
}

func MockFetchMenuDetail(mid uuid.UUID) (*Menu, error) {
	return &Menu{
		ID:      mid,
		Price:   9999,
		InStock: 1,
	}, nil
}

func FetchMenuDetail(mid uuid.UUID) (*Menu, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/%s", os.Getenv("MENU_SERVICE_ENDPOINT"), mid))
	if err != nil {
		return nil, err
	}

	resBody := httpResp{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	return &resBody.Data, nil
}
