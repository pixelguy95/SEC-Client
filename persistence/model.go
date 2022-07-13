package persistence

import "github.com/pixelguy95/sec"

// PersistedCompanyFacts A representation of company facts with appended Timestamp
type PersistedCompanyFacts struct {
	Timestamp int64
	Facts     sec.CompanyFacts
}
