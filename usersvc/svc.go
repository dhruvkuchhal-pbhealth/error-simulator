package usersvc

import "github.com/your-org/error-simulator/models"

// Fetcher is the interface for fetching users (impl lives in another package — genre: interface boundary).
type Fetcher interface {
	FetchUser(id string) (*models.User, error)
}

// Service uses a Fetcher to get users. Panic happens inside the implementation (userfetcher package).
type Service struct {
	Fetcher Fetcher
}

// GetUser delegates to the injected Fetcher. Stack crosses into userfetcher impl.
func (s *Service) GetUser(id string) (*models.User, error) {
	return s.Fetcher.FetchUser(id)
}
