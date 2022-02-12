package sec

const Protocol string = "https://"
const SecUrl string = "www.sec.gov"
const SecReportUrl string = "www.sec.report"
const SecDataUrl string = "data.sec.gov"
const CompanyFactsZip string = "/Archives/edgar/daily-index/xbrl/companyfacts.zip"

const FilingsUrlEndpoint string = Protocol + SecReportUrl + "/Senate-Stock-Disclosures/Filings"
const TickerEndpoint string = Protocol + SecUrl + "/files/company_tickers.json"
const CompanyFactsEndpoint string = Protocol + SecDataUrl + "/api/xbrl/companyfacts"
const AllCompanyFactsEndpoint string = Protocol + SecUrl + CompanyFactsZip

func CompleteUrl(partial string) string {
	return Protocol + SecReportUrl + partial
}
