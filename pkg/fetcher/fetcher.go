package fetcher

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
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
	log.Println("fetching menu: ", fetchMenuURL)
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	res, err := retryClient.Get(fetchMenuURL)
	if err != nil {
		return nil, err
	}

	resBody := httpResp{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	log.Println("received menu data: ", resBody.Data)
	return &resBody.Data, nil
}
