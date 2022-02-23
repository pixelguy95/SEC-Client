package sec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
