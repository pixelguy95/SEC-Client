package sec

import (
	"net/http"
	"time"
)

// getHttpGetRequestWithProperHeaders takes in an endpoint and gives back a http request struct with all headers
// configured as recommended by sec.gov
func getHttpGetRequestWithProperHeaders(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `sec-client`)

	return req, err
}

const SECDateFormat = "2006-01-02"

func StandardSecDateFormatParse(date string) (time.Time, error) {
	return time.Parse(SECDateFormat, date)
}

func StandardSecDateFormatParseSwallowError(date string) time.Time {
	parsedTime, _ := time.Parse(SECDateFormat, date)
	return parsedTime
}
