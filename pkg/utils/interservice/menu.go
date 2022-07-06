package interservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
)

type Interservice interface {
}

// TODO: add menu model
type MenuModel struct {
}

type interservice struct {
	config config.Config
}

func NewInterservice(cfg config.Config) Interservice {
	return interservice{cfg}
}

func (i *interservice) FetchMenuDetailsByID(menuID string) (*MenuModel, error) {
	data := new(MenuModel)
	resp, err := http.Get(fmt.Sprintf("%s/%s", i.config.Get("MENU_SERVICE"), menuID))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
