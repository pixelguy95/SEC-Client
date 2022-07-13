package sec

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const workerCount = 32

// tickerWithFile used internally to pair a file (in-memory) with it's ticker.
type tickerWithFile struct {
	File   zip.File
	Ticker Ticker
}

type TickerWithFacts struct {
	Ticker       Ticker
	CompanyFacts CompanyFacts
	Error        error
}

// GetBulk fetches ALL company facts, and provides them on the returned channel.
// This is highly dependent on your connection speed and your computational resources.
// Typically, it takes a minute or so before the first facts start appearing on the channel.
// Errors related to parsing of individual facts will be sent in the channel and will not be returned
// in this function call
func (client *Client) GetBulk() (chan *TickerWithFacts, error) {

	// Preparing request
	httpClient := &http.Client{}
	req, err := getHttpGetRequestWithProperHeaders(AllCompanyFactsEndpoint)
	if err != nil {
		return nil, err
	}

	// Taking token from bucket
	client.bucket.Take()

	// Performing the request
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	// Copy response to buffer
	buffer := bytes.NewBuffer([]byte{})
	written, err := io.Copy(buffer, response.Body)
	if err != nil {
		return nil, err
	}

	// Create zip reader form buffer
	reader := bytes.NewReader(buffer.Bytes())
	zipReader, err := zip.NewReader(reader, written)
	if err != nil {
		return nil, err
	}

	// Create channels
	fileChannel := make(chan *tickerWithFile, 10000)
	factChannel := make(chan *TickerWithFacts, 10000)

	go readZipFile(client, zipReader, fileChannel, factChannel)

	return factChannel, nil
}

func readZipFile(client *Client, zipReader *zip.Reader,
	fileChannel chan *tickerWithFile, factChannel chan *TickerWithFacts) {

	// Create thread-pool reading from fileChannel
	var waitGroup sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		waitGroup.Add(1)
		go parseCompanyFacts(factChannel, fileChannel, &waitGroup)
	}

	// Send each individual file to the file channel
	for _, file := range zipReader.File {

		// Use the file name to extract the CIK, use the cik to get the proper ticker for the company in question
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

			// Emit the unparsed file together with identifying ticker on file channel
			fileChannel <- &tickerWithFile{Ticker: ticker, File: *file}
		}
	}

	// Close file channel when all fields are sent
	close(fileChannel)

	// Wait for the go routines to process all files and then close the facts channel
	waitGroup.Wait()
	close(factChannel)
}

func parseCompanyFacts(factsChannel chan *TickerWithFacts, filesChannel chan *tickerWithFile, waitGroup *sync.WaitGroup) {

	// Tell wait group that you are finished when file channel has been emptied and closed
	defer waitGroup.Done()

	// Process each file in file channel
	for combo := range filesChannel {

		// Open file for reading, report any errors on the channel
		companyFactsFile, err := combo.File.Open()
		if err != nil {
			err := companyFactsFile.Close()
			if err != nil {
				factsChannel <- &TickerWithFacts{
					Ticker:       Ticker{},
					CompanyFacts: CompanyFacts{},
					Error:        fmt.Errorf("failed to close zip %s for reading: %v", combo.File.Name, err),
				}
				continue
			}

			factsChannel <- &TickerWithFacts{
				Ticker:       Ticker{},
				CompanyFacts: CompanyFacts{},
				Error:        fmt.Errorf("failed to open zip %s for reading: %v", combo.File.Name, err),
			}
			continue
		}

		// Read file contents to a bytes variable
		fileContent, err := io.ReadAll(companyFactsFile)
		if err != nil {
			factsChannel <- &TickerWithFacts{
				Ticker:       Ticker{},
				CompanyFacts: CompanyFacts{},
				Error:        fmt.Errorf("could not read bytes of file %s, %v", combo.File.Name, err),
			}
		}

		// Close file
		err = companyFactsFile.Close()
		if err != nil {
			factsChannel <- &TickerWithFacts{
				Ticker:       Ticker{},
				CompanyFacts: CompanyFacts{},
				Error:        fmt.Errorf("failed to close zip %s for reading: %v", combo.File.Name, err),
			}
			continue
		}

		// Parse file contents into proper facts' struct, ignore companies that have no  data
		var companyFacts CompanyFacts
		err = easyjson.Unmarshal(fileContent, &companyFacts)
		if err != nil || companyFacts.CIK == 0 || (len(companyFacts.Facts.UsGAAP)+len(companyFacts.Facts.DEI)) == 0 {
			continue
		}

		// Emit newly parsed facts on the facts channel
		factsChannel <- &TickerWithFacts{Ticker: combo.Ticker, CompanyFacts: companyFacts, Error: nil}
	}
}
