package sec

import (
	"strings"
	"testing"
)

func TestGetAllTickers(t *testing.T) {
	client := NewClient()
	tickers, err := client.GetAllTickers()

	if err != nil {
		t.Error("Error while getting tickers", err)
	}

	if len(tickers) < 1 {
		t.Error("no tickers in response")
	}

	foundAppleTicker := false
	for _, ticker := range tickers {
		if strings.EqualFold(ticker.Symbol, "AAPL") {
			foundAppleTicker = true
		}

		if len(ticker.Symbol) == 0 {
			t.Error("Found ticker without symbol")
		}

		if ticker.CIK == 0 {
			t.Error("Found ticker with no CIK")
		}
	}

	if !foundAppleTicker {
		t.Error("Expected to find ticker for AAPL in ticker results")
	}
}

func TestGetAllTickersCached(t *testing.T) {
	client := NewClient()
	_, _ = client.GetAllTickers()
	_, _ = client.GetAllTickers()
	_, _ = client.GetAllTickers()
	_, _ = client.GetAllTickers()
	_, _ = client.GetAllTickers()
	tickers, err := client.GetAllTickers()

	if err != nil {
		t.Error("Error while getting tickers", err)
	}

	if len(tickers) < 1 {
		t.Error("no tickers in response")
	}

	foundAppleTicker := false
	for _, ticker := range tickers {
		if strings.EqualFold(ticker.Symbol, "AAPL") {
			foundAppleTicker = true
		}

		if len(ticker.Symbol) == 0 {
			t.Error("Found ticker without symbol")
		}

		if ticker.CIK == 0 {
			t.Error("Found ticker with no CIK")
		}
	}

	if !foundAppleTicker {
		t.Error("Expected to find ticker for AAPL in ticker results")
	}
}
