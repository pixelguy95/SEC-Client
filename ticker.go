package sec

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// Represents a ticker in the given SEC format
type Ticker struct {
	// Central Index Key (CIK) is used on the SEC's computer systems to identify corporations
	// and individual people who have filed disclosure with the SEC
	CIK uint64 `json:"cik_str"`

	// The symbol of the company for exampel NDAQ
	Symbol string `json:"ticker"`

	// The name of the company
	Name string `json:"title"`
}

// Fetches all current tickers from the sec.gov api
func (client *SecClient) GetAllTickers() ([]Ticker, error) {

	// If tickers are cached and younger than 5 minutes, return cached
	if client.cachedTickers != nil {
		if (client.cachedTickersTimeStamp + (5 * 1000 * 60)) > time.Now().UnixMilli() {
			return *client.cachedTickers, nil
		}
	}

	httpClient := &http.Client{}
	req, err := client.GetHttpGetRequestWithProperHeaders(TickerEndpoint)
	if err != nil {
		return nil, err
	}

	client.bucket.Take()
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var tickerMap map[string]Ticker
	err = json.Unmarshal(body, &tickerMap)
	if err != nil {
		return nil, err
	}

	tickers := make([]Ticker, 0, len(tickerMap))
	for _, v := range tickerMap {
		tickers = append(tickers, v)
	}

	client.cachedTickers = &tickers
	client.cachedTickersTimeStamp = time.Now().UnixMilli()

	return tickers, nil
}

func (client *SecClient) GetTickerForSymbol(symbol string) (Ticker, error) {
	tickers, err := client.GetAllTickers()
	if err != nil {
		return Ticker{}, err
	}

	for _, ticker := range tickers {
		if strings.EqualFold(ticker.Symbol, symbol) {
			return ticker, nil
		}
	}

	return Ticker{}, errors.New("ticker with symbol " + symbol + " not found")
}
