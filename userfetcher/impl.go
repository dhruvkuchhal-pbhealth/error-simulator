package userfetcher

import "github.com/your-org/error-simulator/models"

// Impl implements usersvc.Fetcher (interface defined there). Panic happens in this package = interface-boundary genre.
type Impl struct {
	cache map[string]*models.User
}

// NewImpl returns an impl with empty cache (cache miss → nil → panic on deref).
func NewImpl() *Impl {
	return &Impl{cache: make(map[string]*models.User)}
}

// FetchUser implements Fetcher. BUG: on cache miss u is nil; we "normalize" and deref → panic in this file.
func (i *Impl) FetchUser(id string) (*models.User, error) {
	var u *models.User
	if v, ok := i.cache[id]; ok {
		u = v
	}
	// BUG: no nil check; panic here (genre: interface boundary — impl in different pkg)
	_ = u.ID
	return u, nil
}
