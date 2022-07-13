package sec

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"io"
	"net/http"
	"strings"
	"time"
)

// Ticker Represents a ticker in the given SEC format
type Ticker struct {
	// Central Index Key (CIK) is used on the SEC's computer systems to identify corporations
	// and individual people who have filed disclosure with the SEC
	CIK uint64 `json:"cik_str"`

	// The symbol of the company for example NDAQ
	Symbol string `json:"ticker"`

	// The name of the company
	Name string `json:"title"`
}

// GetAllTickers Fetches all current tickers from the sec.gov api, uses caching for speed
func (client *Client) GetAllTickers() ([]Ticker, error) {

	// If tickers are cached and younger than 5 minutes, return copy from cached
	client.cachedTickersLock.RLock()
	if client.cachedTickers != nil {
		if (client.cachedTickersTimeStamp + (5 * 1000 * 60)) > time.Now().UnixMilli() {
			copyOfAllTickers, err := getCopyOfTickers(*client.cachedTickers)
			if err != nil {
				return nil, err
			}

			defer client.cachedTickersLock.RUnlock()
			return copyOfAllTickers, nil
		}
	}

	client.cachedTickersLock.RUnlock()
	client.cachedTickersLock.Lock()
	httpClient := &http.Client{}
	req, err := getHttpGetRequestWithProperHeaders(TickerEndpoint)
	if err != nil {
		return nil, err
	}

	client.bucket.Take()
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Fatal error occurred while closing body: %s\n", err)
		}
	}(response.Body)

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

	copyOfAllTickers, err := getCopyOfTickers(tickers)
	if err != nil {
		return nil, err
	}

	defer client.cachedTickersLock.Unlock()
	return copyOfAllTickers, nil
}

func (client *Client) GetTickerForSymbol(symbol string) (Ticker, error) {
	tickers, err := client.GetAllTickers()
	if err != nil {
		return Ticker{}, err
	}

	for _, ticker := range tickers {
		if strings.EqualFold(ticker.Symbol, symbol) {
			copyOfTicker, err := getCopyOfTicker(ticker)
			if err != nil {
				return Ticker{}, err
			}

			return copyOfTicker, nil
		}
	}

	return Ticker{}, errors.New("ticker with symbol " + symbol + " not found")
}

func (client *Client) GetTickerForCIK(cik uint64) (Ticker, error) {
	tickers, err := client.GetAllTickers()
	if err != nil {
		return Ticker{}, err
	}

	for _, ticker := range tickers {
		if ticker.CIK == cik {
			copyOfTicker, err := getCopyOfTicker(ticker)
			if err != nil {
				return Ticker{}, err
			}

			return copyOfTicker, nil
		}
	}

	return Ticker{}, errors.New("ticker with cik " + fmt.Sprintf("%d", cik) + " not found")
}

// getCopyOfTickers returns an exact copy of the tickers given, used so that when a client requests all tickers they
// don't use the exact same memory as the client.
// Might be room for optimization here
func getCopyOfTickers(tickers []Ticker) ([]Ticker, error) {
	var copyOfTickers []Ticker
	for _, ticker := range tickers {
		copiedTicker, err := getCopyOfTicker(ticker)
		if err != nil {
			return nil, err
		}

		copyOfTickers = append(copyOfTickers, copiedTicker)
	}

	return copyOfTickers, nil
}

// getCopyOfTickers copies a single ticker
func getCopyOfTicker(t Ticker) (Ticker, error) {
	copyOfTicker := Ticker{}
	err := copier.Copy(&copyOfTicker, &t)
	if err != nil {
		return Ticker{}, err
	}

	return copyOfTicker, nil
}
