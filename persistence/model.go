package persistence

import "github.com/pixelguy95/sec"

type PersistedCompanyFacts struct {
	Timestamp int64
	Facts     sec.CompanyFacts
}
