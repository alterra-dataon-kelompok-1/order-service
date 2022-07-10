package fetcher

import "github.com/google/uuid"

type MockFetcher struct{}

func (f *MockFetcher) FetchMenuDetail(mid uuid.UUID) (*Menu, error) {
	return &Menu{
		ID:      mid,
		Price:   10_000,
		InStock: 1,
	}, nil
}
