package sec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CompanyFacts struct {
	CIK        uint64 `json:"cik"`
	EntityName string `json:"entityName"`
	Facts      Facts  `json:"facts"`
}

// A representation of a
type Facts struct {
	DEI    map[string]Fact `json:"dei"`
	UsGAAP map[string]Fact `json:"us-gaap"`
}

type Fact struct {
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Units       map[string][]Unit `json:"units"`
}

type Unit struct {
	Value float64 `json:"val"`

	Start string `json:"start"`
	End   string `json:"end"`

	FiscalYear   uint16 `json:"fy"`
	FiscalPeriod string `json:"fp"`

	Account string `json:"accn"`

	Form    string `json:"form"`
	FiledOn string `json:"filed"`
}

func (client *SecClient) GetAllFactsForTicker(ticker Ticker) (CompanyFacts, error) {

	if client.persistenceLayer != nil {
		persistedFacts, err := client.persistenceLayer.LoadFacts(ticker)
		if err != nil {
			return CompanyFacts{}, err
		}

		if persistedFacts != nil {
			return *persistedFacts, nil
		}
	}

	httpClient := &http.Client{}
	req, err := client.GetHttpGetRequestWithProperHeaders(factEndpointUrl(ticker))

	client.bucket.Take()
	response, err := httpClient.Do(req)

	if err != nil {
		return CompanyFacts{}, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return CompanyFacts{}, err
	}

	var companyFacts CompanyFacts
	err = json.Unmarshal(body, &companyFacts)
	if err != nil {
		return CompanyFacts{}, err
	}

	if client.persistenceLayer != nil {
		err := client.persistenceLayer.SaveFacts(ticker, &companyFacts)
		if err != nil {
			return CompanyFacts{}, nil
		}
	}

	return companyFacts, nil
}

func factEndpointUrl(ticker Ticker) string {
	return CompanyFactsEndpoint + "/" + "CIK" + fmt.Sprintf("%010d", ticker.CIK) + ".json"
}
