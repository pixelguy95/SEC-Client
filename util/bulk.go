package util

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pixelguy95/sec"
)

func ProcessBulkCompanyFactsData(client *sec.SecClient, persistenceLayer sec.PersistenceLayer) error {

	fmt.Println("Preparing request")
	httpClient := &http.Client{}
	req, err := client.GetHttpGetRequestWithProperHeaders(sec.AllCompanyFactsEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("Taking token from bucket")
	client.WaitForToken()

	fmt.Println("Performing request...")
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("Reading response to buffer")
	defer response.Body.Close()
	buff := bytes.NewBuffer([]byte{})
	written, err := io.Copy(buff, response.Body)
	if err != nil {
		return err
	}

	fmt.Println("Creating zip reader")
	reader := bytes.NewReader(buff.Bytes())
	zipReader, err := zip.NewReader(reader, written)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		fileread, err := file.Open()
		if err != nil {
			msg := "Failed to open zip %s for reading: %s"
			return fmt.Errorf(msg, file.Name, err)
		}
		defer fileread.Close()

		body, err := io.ReadAll(fileread)
		if err != nil {
			return err
		}

		var companyFacts sec.CompanyFacts
		err = json.Unmarshal(body, &companyFacts)
		if err != nil || companyFacts.CIK == 0 || (len(companyFacts.Facts.UsGAAP)+len(companyFacts.Facts.DEI)) == 0 {
			continue
		}

		ticker, err := client.GetTickerForCIK(companyFacts.CIK)
		if err != nil {
			fmt.Printf("%-58s %010d %8s %5d NOT PERSISTED\n", companyFacts.EntityName, companyFacts.CIK, "N/A", (len(companyFacts.Facts.UsGAAP) + len(companyFacts.Facts.DEI)))
			continue
		}

		if persistenceLayer != nil {
			persistenceLayer.SaveFacts(ticker, &companyFacts)
		}

		fmt.Printf("%-58s %010d %8s %5d\n", companyFacts.EntityName, companyFacts.CIK, ticker.Symbol, (len(companyFacts.Facts.UsGAAP) + len(companyFacts.Facts.DEI)))
	}

	return nil
}
