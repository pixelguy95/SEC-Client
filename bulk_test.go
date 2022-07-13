package sec

/*
func TestBulk(t *testing.T) {
	client := NewClient()
	factChannel, err := client.GetBulk()
	if err != nil {
		t.Fatal(err)
	}

	expectedTickers := []string{"TSLA", "AAPL", "MSFT", "NDAQ"}

	count := 0
	errors := 0
	factsCount := 0
	for tickerFact := range factChannel {
		if tickerFact.Error != nil {
			fmt.Printf("Error Detected %s \n", tickerFact.Error)
			errors++
			t.Fatal(tickerFact.Error)
		}

		for i, expectedTicker := range expectedTickers {
			if tickerFact.Ticker.Symbol == expectedTicker {
				expectedTickers = append(expectedTickers[:i], expectedTickers[i+1:]...)
				fmt.Printf("Found expected ticker %s\nRemaining %s\n", expectedTicker, expectedTickers)
			}
		}

		if len(expectedTickers) == 0 {
			break
		}

		fmt.Printf("%-58s %010d %8s %5d\n", tickerFact.CompanyFacts.EntityName, tickerFact.CompanyFacts.CIK,
			tickerFact.Ticker.Symbol, len(tickerFact.CompanyFacts.Facts.UsGAAP)+len(tickerFact.CompanyFacts.Facts.DEI))
		factsCount += len(tickerFact.CompanyFacts.Facts.UsGAAP) + len(tickerFact.CompanyFacts.Facts.DEI)
		count++
	}

	fmt.Printf("%-20s %10d\n%-20s %10d\n%-20s %10d\n", "Ticker Count:", count, "Facts Count:", factsCount, "Errors Count:", errors)
}

*/
