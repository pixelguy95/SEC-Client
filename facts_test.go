package sec

import (
	"fmt"
	"testing"
)

func TestGetCompanyFacts(t *testing.T) {

	client := NewClient()

	companyFacts, err := client.GetAllFactsForTicker(Ticker{CIK: 320193})
	if err != nil {
		t.Error("Error while getting facts")
		t.Fatal(err)
	}

	if companyFacts.EntityName != "Apple Inc." {
		t.Fatal("Company name was not expected Apple Inc.")
	}

	liabilitiesNonCurrent, ok := companyFacts.Facts.UsGAAP["LiabilitiesNoncurrent"]
	if ok == false {
		t.Fatal("Long term liabilities was nil")
	}

	if len(liabilitiesNonCurrent.Units["USD"]) <= 0 {
		t.Fatal("Expected to find non-current liabilities (LiabilitiesNoncurrent)")
	}

	fmt.Println("Liabilities non-current")
	for _, unit := range liabilitiesNonCurrent.Units["USD"] {
		fmt.Println(unit)
		if !unit.IsInstant() || unit.IsPeriod() {
			t.Fatal("Liabilities non-current should be instant and not period")
		}
	}

	liabilitiesCurrent, ok := companyFacts.Facts.UsGAAP["LiabilitiesCurrent"]
	if ok == false {
		t.Fatal("Short term liabilities was nil")
	}

	if len(liabilitiesCurrent.Units["USD"]) <= 0 {
		t.Fatal("Expected to find current liabilities (LiabilitiesCurrent)")
	}

	fmt.Println()
	fmt.Println("Liabilities current")
	for _, unit := range liabilitiesCurrent.Units["USD"] {
		fmt.Println(unit)
	}

	revenues, ok := companyFacts.Facts.UsGAAP["Revenues"]
	if ok == false {
		t.Fatal("Revenues was nil")
	}

	if len(revenues.Units["USD"]) <= 0 {
		t.Fatal("Expected to find revenues (revenues)")
	}

	fmt.Println()
	fmt.Println("Revenues")
	for _, unit := range revenues.Units["USD"] {
		fmt.Println(unit)
		if !unit.IsPeriod() || unit.IsInstant() {
			t.Fatal("Revenues should be period and not instant")
		}
	}
}
