package sec

import "testing"

func TestUnitUtils(t *testing.T) {
	unit := &Unit{Value: 0, Start: "2021-01-01", End: "2022-01-01", FiscalYear: 2022,
		FiscalPeriod: "FY", Account: "A", FiledOn: "2022-01-02", Form: "10-K"}

	if !unit.IsPeriod() {
		t.Fatal("unit was not considered period")
	}

	if unit.IsQuarterRange() {
		t.Fatal("unit was considered to be quarter, even though it is not")
	}

	if !unit.IsYearRange() {
		t.Fatal("unit was not considered to be year range")
	}

	if unit.IsInstant() {
		t.Fatal("unit was considered to be instant, even though it is not")
	}

	if !unit.EndAsTime().Equal(StandardSecDateFormatParseSwallowError("2022-01-01")) {
		t.Fatalf("unit end time was not what was expected")
	}

	if !unit.StartAsTime().Equal(StandardSecDateFormatParseSwallowError("2021-01-01")) {
		t.Fatal("unit start time was not what was expected")
	}

	if !unit.FiledOnAsTime().Equal(StandardSecDateFormatParseSwallowError("2022-01-02")) {
		t.Fatal("unit filed on time was not what was expected")
	}
}

func TestUnitUtilsWithQuarterly(t *testing.T) {
	unit := &Unit{Value: 0, Start: "2021-01-01", End: "2021-03-31", FiscalYear: 2022,
		FiscalPeriod: "FY", Account: "A", FiledOn: "2022-01-02", Form: "10-K"}

	if !unit.IsPeriod() {
		t.Fatal("unit was not considered period")
	}

	if !unit.IsQuarterRange() {
		t.Fatal("unit was not considered to be quarter")
	}

	if unit.IsYearRange() {
		t.Fatal("unit was  considered to be year range, even though it is not")
	}

}
