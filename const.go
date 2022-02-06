package sec

const Protocol string = "https://"
const SecUrl string = "www.sec.gov"
const SecReportUrl string = "www.sec.report"
const SecDataUrl string = "data.sec.gov"

const FilingsUrlEndpoint string = Protocol + SecReportUrl + "/Senate-Stock-Disclosures/Filings"
const TickerEndpoint string = Protocol + SecUrl + "/files/company_tickers.json"
const CompanyFactsEndpoint string = Protocol + SecDataUrl + "/api/xbrl/companyfacts"

func CompleteUrl(partial string) string {
	return Protocol + SecReportUrl + partial
}
