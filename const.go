package sec

const protocol string = "https://"
const secUrl string = "www.sec.gov"
const secReportUrl string = "www.sec.report"
const secDataUrl string = "data.sec.gov"
const companyFactsZip string = "/Archives/edgar/daily-index/xbrl/companyfacts.zip"

const filingsUrlEndpoint = protocol + secReportUrl + "/Senate-Stock-Disclosures/Filings"
const tickerEndpoint = protocol + secUrl + "/files/company_tickers.json"
const companyFactsEndpoint = protocol + secDataUrl + "/api/xbrl/companyfacts"
const allCompanyFactsEndpoint = protocol + secUrl + companyFactsZip

func CompleteUrl(partial string) string {
	return protocol + secReportUrl + partial
}
