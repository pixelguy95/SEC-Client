package persistence

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pixelguy95/sec"
)

func TestRedisPersistenceLayer(t *testing.T) {

	redisClient := NewRedisPersistenceLaye(&redis.Options{Addr: "192.168.0.102:6379", Password: "", DB: 0})

	client := sec.NewSecClientWithPersistence(redisClient)

	ticker, _ := client.GetTickerForSymbol("AAPL")

	start := time.Now()
	client.GetAllFactsForTicker(ticker)
	fmt.Println(time.Since(start))
	start = time.Now()
	client.GetAllFactsForTicker(ticker)
	fmt.Println(time.Since(start))

	start = time.Now()
	facts, err := redisClient.LoadFacts(ticker)
	fmt.Println(time.Since(start))

	if err != nil {
		t.Fatal(err)
	}

	if facts.EntityName != "Apple Inc." {
		t.Fatal("Company name was not as expected, expected Apple INC., was " + facts.EntityName)
	}

	if len(facts.Facts.UsGAAP) < 486 {
		t.Fatalf("Expected at least 486 facts, was %d", len(facts.Facts.UsGAAP))
	}
}
