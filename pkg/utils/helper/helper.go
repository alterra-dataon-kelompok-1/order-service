package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
)

// TODO: Implement get item price (API call to other table)
func GetItemPrice(id uuid.UUID) float32 {
	return float32(3000)
}

func GetItemPriceImplement(uid uuid.UUID) float32 {
	resp, err := http.Get(fmt.Sprintf("localhost/menu/%s", uid))
	if err != nil {
		log.Println(err)
	}

	menu := model.Menu{}

	err = json.NewDecoder(resp.Body).Decode(&menu)
	// menu.Price, menu.InStock

	return float32(3000)
}
