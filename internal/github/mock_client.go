package github

import (
	"time"
)

// MockClient implements the Client interface (or matches your client methods for testing).
type MockClient struct {
	Contributors []Contributor
	Err          error
}

func NewMockClient(contributors []Contributor, err error) *MockClient {
	return &MockClient{
		Contributors: contributors,
		Err:          err,
	}
}

// Correct receiver name (*MockClient) and matching return type ([]Contributor)
func (m *MockClient) GetTopContributors(since time.Time) ([]Contributor, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Contributors, nil
}
