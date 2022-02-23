package sec

import "time"

type CompanyFacts struct {
	CIK        uint64 `json:"cik"`
	EntityName string `json:"entityName"`
	Facts      Facts  `json:"facts"`
}

// A representation of a
type Facts struct {
	DEI    map[string]Fact `json:"dei"`
	UsGAAP map[string]Fact `json:"us-gaap"`
}

type Fact struct {
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Units       map[string][]Unit `json:"units"`
}

type Unit struct {
	Value float64 `json:"val"`

	Start string `json:"start"`
	End   string `json:"end"`

	FiscalYear   uint16 `json:"fy"`
	FiscalPeriod string `json:"fp"`

	Account string `json:"accn"`

	Form    string `json:"form"`
	FiledOn string `json:"filed"`
}

func (unit *Unit) IsInstant() bool {
	return unit.End != "" && unit.Start == ""
}

func (unit *Unit) IsPeriod() bool {
	return unit.End != "" && unit.Start != ""
}

func (unit *Unit) StartAsTime() time.Time {
	return StandardSecDateFormatParseSwallowError(unit.Start)
}

func (unit *Unit) EndAsTime() time.Time {
	return StandardSecDateFormatParseSwallowError(unit.End)
}

func (unit *Unit) FiledOnAsTime() time.Time {
	return StandardSecDateFormatParseSwallowError(unit.FiledOn)
}

func (unit *Unit) IsQuarterRange() bool {
	return unit.EndAsTime().Sub(unit.StartAsTime()) > (time.Hour*24*85) && unit.EndAsTime().Sub(unit.StartAsTime()) < (time.Hour*24*95)
}

func (unit *Unit) IsYearRange() bool {
	return unit.EndAsTime().Sub(unit.StartAsTime()) > (time.Hour*24*360) && unit.EndAsTime().Sub(unit.StartAsTime()) < (time.Hour*24*370)
}
