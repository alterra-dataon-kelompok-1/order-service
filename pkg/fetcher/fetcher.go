package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func NewFetcher(uri string) Fetcher {
	return &fetcher{ExternalServiceURL: uri}
}

type Fetcher interface {
	FetchMenuDetail(mid uuid.UUID) (*Menu, error)
}

type fetcher struct {
	ExternalServiceURL string
}

func (f *fetcher) FetchMenuDetail(mid uuid.UUID) (*Menu, error) {
	fetchMenuURL := fmt.Sprintf("http://%s/%s", f.ExternalServiceURL, mid)
	res, err := http.Get(fetchMenuURL)
	if err != nil {
		return nil, err
	}

	resBody := httpResp{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	return &resBody.Data, nil
}
