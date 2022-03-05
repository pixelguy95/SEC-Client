package util

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/pixelguy95/sec"
)

type TickerWithFile struct {
	File   zip.File
	Ticker sec.Ticker
}

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

	fmt.Println("Setting up thread pool")
	fileChannel := make(chan *TickerWithFile, 10000)

	var waitGroup sync.WaitGroup
	for i := 0; i < 255; i++ {
		waitGroup.Add(1)
		go parseCompanyFacts(*client, persistenceLayer, fileChannel, &waitGroup)
	}

	fmt.Println("Sending files to thread pool for parsing")
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, ".json") {
			cik, err := strconv.Atoi(strings.TrimLeft(strings.TrimLeft(strings.Split(file.Name, ".")[0], "0"), "CIK"))
			if err != nil {
				fmt.Printf("Error while parsing cik from %s\n", file.Name)
				continue
			}

			ticker, err := client.GetTickerForCIK(uint64(cik))
			if err != nil {
				continue
			}

			fileChannel <- &TickerWithFile{Ticker: ticker, File: *file}
		}
	}

	fmt.Println("Waiting for goroutines to finnish")
	close(fileChannel)
	waitGroup.Wait()

	fmt.Println("Done")

	return nil
}

func parseCompanyFacts(client sec.SecClient, persistenceLayer sec.PersistenceLayer,
	filesChannel chan *TickerWithFile,
	waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	for combo := range filesChannel {
		companyFactsFile, err := combo.File.Open()
		if err != nil {
			fmt.Printf("failed reading file %s\n", combo.File.Name)
			companyFactsFile.Close()
			return fmt.Errorf("failed to open zip %s for reading: %v", combo.File.Name, err)
		}

		fileContent, err := io.ReadAll(companyFactsFile)
		if err != nil {
			panic("Could not read bytes of file %s")
		}

		companyFactsFile.Close()

		var companyFacts sec.CompanyFacts
		err = json.Unmarshal(fileContent, &companyFacts)
		if err != nil || companyFacts.CIK == 0 || (len(companyFacts.Facts.UsGAAP)+len(companyFacts.Facts.DEI)) == 0 {
			continue
		}

		if persistenceLayer != nil {
			persistenceLayer.SaveFacts(combo.Ticker, &companyFacts)
		}

		fmt.Printf("%-58s %010d %8s %5d\n", companyFacts.EntityName, companyFacts.CIK, combo.Ticker.Symbol, (len(companyFacts.Facts.UsGAAP) + len(companyFacts.Facts.DEI)))
	}

	return nil
}
