package cachesvc

import "github.com/your-org/error-simulator/models"

// CacheService gets users with cache fallback to repo (layer 2). Panic in repo (layer 3).
type CacheService struct {
	repo *Repo
}

// NewCacheService returns a cache service using the given repo.
func NewCacheService(repo *Repo) *CacheService {
	return &CacheService{repo: repo}
}

// GetUserByID returns user from cache or repo. On miss, calls repo.FindByID → panic if repo.db is nil.
func (c *CacheService) GetUserByID(id string) (*models.User, error) {
	// No in-memory cache for this test; always miss and hit repo.
	return c.repo.FindByID(id)
}
