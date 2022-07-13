package sec

import (
	"testing"
)

type RecorderPersistenceLayer struct {
	shouldLoad bool
	loads      map[string]int
	saves      map[string]int
}

func NewRecorderPersistenceLayer() (*RecorderPersistenceLayer, error) {
	return &RecorderPersistenceLayer{loads: make(map[string]int), saves: make(map[string]int), shouldLoad: true}, nil
}

func (persistenceLayer *RecorderPersistenceLayer) SaveFacts(ticker Ticker, _ *CompanyFacts) error {
	persistenceLayer.saves[ticker.Symbol]++
	return nil
}

func (persistenceLayer *RecorderPersistenceLayer) LoadFacts(ticker Ticker) (*CompanyFacts, error) {
	persistenceLayer.loads[ticker.Symbol]++

	if persistenceLayer.shouldLoad {
		return nil, nil
	}

	return &CompanyFacts{}, nil
}

func TestPersistenceLayer(t *testing.T) {
	recorder, _ := NewRecorderPersistenceLayer()
	client := NewClientWithPersistence(recorder)

	ticker, _ := client.GetTickerForSymbol("AAPL")

	_, _ = client.GetAllFactsForTicker(ticker)

	recorder.shouldLoad = false

	if recorder.loads["AAPL"] != 1 {
		t.Fatal("Expected exactly 1 load call to persistence layer")
	}

	if recorder.saves["AAPL"] != 1 {
		t.Fatal("Expected exactly 1 save call to persistence layer")
	}

	_, _ = client.GetAllFactsForTicker(ticker)
	_, _ = client.GetAllFactsForTicker(ticker)
	_, _ = client.GetAllFactsForTicker(ticker)

	if recorder.loads["AAPL"] != 4 {
		t.Fatal("Expected exactly 4 load call to persistence layer")
	}

	if recorder.saves["AAPL"] != 1 {
		t.Fatal("Expected exactly 1 save call to persistence layer")
	}

	recorder.shouldLoad = true

	_, _ = client.GetAllFactsForTicker(ticker)

	if recorder.loads["AAPL"] != 5 {
		t.Fatal("Expected exactly 5 load call to persistence layer")
	}

	if recorder.saves["AAPL"] != 2 {
		t.Fatal("Expected exactly 2 save call to persistence layer")
	}
}
