package sec

import (
	"net/http"

	"go.uber.org/ratelimit"
)

type SecClient struct {
	// The SEC only allows 10 requests per second, so this client has built in
	// rate limiter.
	bucket ratelimit.Limiter

	cachedTickersTimeStamp int64
	cachedTickers          *[]Ticker

	persistenceLayer PersistenceLayer
}

func NewSecClient() *SecClient {
	bucket := ratelimit.New(10)
	return &SecClient{bucket: bucket, persistenceLayer: nil}
}

func NewSecClientWithPersistence(persistenceLayer PersistenceLayer) *SecClient {
	bucket := ratelimit.New(10)
	return &SecClient{bucket: bucket, persistenceLayer: persistenceLayer}
}

func (client *SecClient) getHttpGetRequestWithProperHeaders(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `sec-client`)

	return req, err
}

type PersistenceLayer interface {
	SaveFacts(Ticker, *CompanyFacts) error
	LoadFacts(Ticker) (*CompanyFacts, error)
}
