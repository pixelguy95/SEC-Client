package sec

const Protocol string = "https://"
const SecUrl string = "www.sec.gov"
const SecReportUrl string = "www.sec.report"
const SecDataUrl string = "data.sec.gov"
const CompanyFactsZip string = "/Archives/edgar/daily-index/xbrl/companyfacts.zip"

const FilingsUrlEndpoint = Protocol + SecReportUrl + "/Senate-Stock-Disclosures/Filings"
const TickerEndpoint = Protocol + SecUrl + "/files/company_tickers.json"
const CompanyFactsEndpoint = Protocol + SecDataUrl + "/api/xbrl/companyfacts"
const AllCompanyFactsEndpoint = Protocol + SecUrl + CompanyFactsZip

func CompleteUrl(partial string) string {
	return Protocol + SecReportUrl + partial
}
