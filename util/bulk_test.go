package util

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pixelguy95/sec"
	"github.com/pixelguy95/sec/persistence"
)

func TestBulk(t *testing.T) {
	tempFile, err := os.Create("./full-facts-" + time.Now().Format(sec.SECDateFormat))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bolt, err := persistence.NewBoltPersistenceLayer(tempFile.Name(), time.Hour)
	client := sec.NewSecClientWithPersistence(bolt)

	err = ProcessBulkCompanyFactsData(client, bolt)
	if err != nil {
		t.Fatal(err)
	}
}
