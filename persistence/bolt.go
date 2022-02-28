package persistence

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/pixelguy95/sec"
	"go.etcd.io/bbolt"
)

type BoltPersistenceLayer struct {
	Db     *bbolt.DB
	config BoltPersistenceLayerConfig
}

type BoltPersistenceLayerConfig struct {
	Path         string
	ExpiresAfter time.Duration
}

func NewBoltPersistenceLayer(config BoltPersistenceLayerConfig) (*BoltPersistenceLayer, error) {
	db, err := bbolt.Open(config.Path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("companyFacts"))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BoltPersistenceLayer{Db: db, config: config}, nil
}

func (persistenceLayer *BoltPersistenceLayer) SaveFacts(ticker sec.Ticker, facts *sec.CompanyFacts) error {
	err := persistenceLayer.Db.Update(func(tx *bbolt.Tx) error {
		var pcf = &PersistedCompanyFacts{Facts: *facts, Timestamp: time.Now().UnixMilli()}
		b := bytes.Buffer{}
		encoder := gob.NewEncoder(&b)
		encoder.Encode(pcf)

		bucket := tx.Bucket([]byte("companyFacts"))
		bucket.Put([]byte(ticker.Symbol), b.Bytes())

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (persistenceLayer *BoltPersistenceLayer) LoadFacts(ticker sec.Ticker) (*sec.CompanyFacts, error) {

	pcf := &PersistedCompanyFacts{}

	err := persistenceLayer.Db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte("companyFacts"))
		content := bucket.Get([]byte(ticker.Symbol))

		if content == nil {
			pcf = nil
			return nil
		}

		buffer := &bytes.Buffer{}
		buffer.Write(content)
		decoder := gob.NewDecoder(buffer)

		err := decoder.Decode(pcf)
		if err != nil {
			return err
		}

		if (pcf.Timestamp + persistenceLayer.config.ExpiresAfter.Milliseconds()) < time.Now().UnixMilli() {
			bucket.Delete([]byte(ticker.Symbol))
			pcf = nil
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if pcf == nil {
		return nil, nil
	}

	return &pcf.Facts, nil
}
