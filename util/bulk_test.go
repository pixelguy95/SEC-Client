package util

import (
	"os"
	"testing"
	"time"

	"github.com/pixelguy95/sec"
	"github.com/pixelguy95/sec/persistence"
)

func TestBulk(t *testing.T) {
	tempFile, err := os.Create("L:/full-facts-" + time.Now().Format(sec.SECDateFormat))
	if err != nil {
		t.Fatal(err)
	}

	bolt, err := persistence.NewBoltPersistenceLayer(
		persistence.BoltPersistenceLayerConfig{Path: tempFile.Name(), ExpiresAfter: time.Hour, ReduceSize: true})
	if err != nil {
		t.Fatal(err)
	}

	client := sec.NewSecClientWithPersistence(bolt)

	err = ProcessBulkCompanyFactsData(client, bolt)
	if err != nil {
		t.Fatal(err)
	}
}
