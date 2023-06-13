package service

import "sync"

// RatingStore is an interface to store laptop ratings
type RatingStore interface {
	// Add adds a new laptop score to the store and returns its rating
	Add(laptopId string, score float64) (*Rating, error)
}

// Rating contains the rating information of a laptop
type Rating struct {
	Count uint32
	Sum   float64
}

// InMemoryRatingStore stores laptop rating in memory
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

// NewInMemoryRatingStore returns a new in memory rating store
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

// Add adds a new laptop score to the store and return its rating
func (store *InMemoryRatingStore) Add(laptopId string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopId]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}
	store.rating[laptopId] = rating
	return rating, nil
}
