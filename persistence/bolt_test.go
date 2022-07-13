package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/pixelguy95/sec"
)

var bolt *BoltPersistenceLayer

func TestMain(m *testing.M) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "test_bolt_")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bolt, err = NewBoltPersistenceLayer(BoltPersistenceLayerConfig{Path: tempFile.Name(), ExpiresAfter: time.Hour, ReduceSize: false})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.Run()
	os.Exit(0)
}

func TestBoltPersistenceLayer(t *testing.T) {
	client := sec.NewClientWithPersistence(bolt)

	ticker, _ := client.GetTickerForSymbol("AAPL")

	_, _ = client.GetAllFactsForTicker(ticker)
	_, _ = client.GetAllFactsForTicker(ticker)
	_, _ = client.GetAllFactsForTicker(ticker)
	_, _ = client.GetAllFactsForTicker(ticker)

	facts, err := bolt.LoadFacts(ticker)

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
