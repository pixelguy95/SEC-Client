package persistence

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pixelguy95/sec"
)

type RedisPersistenceLayer struct {
	client *redis.Client
	config *redis.Options
}

func NewRedisPersistenceLaye(config *redis.Options) *RedisPersistenceLayer {
	client := redis.NewClient(config)
	return &RedisPersistenceLayer{
		client: client, config: config,
	}
}

var ctx = context.Background()

func (persistenceLayer *RedisPersistenceLayer) SaveFacts(ticker sec.Ticker, facts *sec.CompanyFacts) error {

	var pcf = &PersistedCompanyFacts{Facts: *facts, Timestamp: time.Now().UnixMilli()}
	b := bytes.Buffer{}
	encoder := gob.NewEncoder(&b)
	encoder.Encode(pcf)

	err := persistenceLayer.client.Set(ctx, ticker.Symbol, b.Bytes(), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (persistenceLayer *RedisPersistenceLayer) LoadFacts(ticker sec.Ticker) (*sec.CompanyFacts, error) {

	pcf := &PersistedCompanyFacts{}

	start := time.Now()
	result := persistenceLayer.client.Get(ctx, ticker.Symbol)
	fmt.Printf("%-15s %10s\n", "GET", time.Since(start))
	content, err := result.Bytes()
	fmt.Printf("%-15s %10s %10d bytes\n", "BYTES", time.Since(start), len(content))
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	buffer := &bytes.Buffer{}
	buffer.Write(content)
	fmt.Printf("%-15s %10s\n", "BUFFER", time.Since(start))
	decoder := gob.NewDecoder(buffer)
	fmt.Printf("%-15s %10s\n", "DECODER", time.Since(start))

	err = decoder.Decode(pcf)
	fmt.Printf("%-15s %10s\n", "DECODED", time.Since(start))
	if err != nil {
		return nil, err
	}

	return &pcf.Facts, nil
}
