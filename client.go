package sec

import (
	"sync"

	"go.uber.org/ratelimit"
)

type Client struct {
	// The SEC only allows 10 requests per second, so this client has built in rate limiter.
	bucket ratelimit.Limiter

	// Cached tickers with lock and timestamp for refreshing cache
	cachedTickers          *[]Ticker
	cachedTickersTimeStamp int64
	cachedTickersLock      sync.RWMutex

	persistenceLayer PersistenceLayer
}

// NewClient creates a new client pointer without any persistence layer, meaning that all requests will go directly
// to sec.gov and will have to be rate limited.
func NewClient() *Client {
	bucket := ratelimit.New(10)
	return &Client{bucket: bucket, persistenceLayer: nil}
}

// NewClientWithPersistence creates a new client pointer with given persistence layer implementation
func NewClientWithPersistence(persistenceLayer PersistenceLayer) *Client {
	bucket := ratelimit.New(10)
	return &Client{bucket: bucket, persistenceLayer: persistenceLayer}
}

// PersistenceLayer is an interface that you can use to create custom persistence layers for the sec client
// The client makes no assumption about thread safety, that has to be handled by the implementation if required
type PersistenceLayer interface {
	SaveFacts(Ticker, *CompanyFacts) error
	LoadFacts(Ticker) (*CompanyFacts, error)
}
