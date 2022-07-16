# SEC Client
<p align="center">
  <img width="200" src="https://upload.wikimedia.org/wikipedia/commons/5/54/United_States_Securities_and_Exchange_Commission.svg">
</p>

This is an **unofficial** client for accessing the filing APIs that the Securities and Exchange Commision provides. 

**NOTE:** This is a work in progress and is subject to change.

## Installation
```
go get github.com/pixelguy95/sec
```

Built using golang version `1.18`

## Usage
Create a new client instance using `sec.NewClient()`, then use the client to access the ticker and filing data. Below is an example where all `Non-current liabilities` (See [US-GAAP](https://www.investopedia.com/terms/g/gaap.asp), and [XBRL](https://www.xbrl.org/the-standard/what/an-introduction-to-xbrl/)) are fetched for `Apple Inc`. 

```go
secClient := sec.NewClient()

appleTicker, err := secClient.GetTickerForSymbol("AAPL")
if err != nil {
    panic(err)
}

appleFacts, err := secClient.GetAllFactsForTicker(appleTicker)
if err != nil {
    panic(err)
}

nonCurrentLiabilities := appleFacts.Facts.UsGAAP["LiabilitiesNoncurrent"].Units["USD"]
for _, unit := range nonCurrentLiabilities {
    fmt.Printf("Non-current liabilities @ %s: %12.2f$\n", unit.End, unit.Value)
}
```

```
...
Non-current liabilities @ 2020-09-26: 153157000000.00$
Non-current liabilities @ 2020-12-26: 155323000000.00$
Non-current liabilities @ 2021-03-27: 161595000000.00$
Non-current liabilities @ 2021-06-26: 157806000000.00$
Non-current liabilities @ 2021-09-25: 162431000000.00$
Non-current liabilities @ 2021-12-25: 161685000000.00$
Non-current liabilities @ 2022-03-26: 155755000000.00$
...
```
To get a list of all possible tickers use `secClient.GetAllTickers()`

### Persistence layer
This client intends to follow the compliance guidance set by SEC (See next section). One of the more significant limits set, is the limit to the number of data requests you can do per second. To mitigate this the client supports a persistence layer that will be hit first before going to the actual source. 

You can create your own persistence layer using the simple interface provided
```go
type PersistenceLayer interface {
    SaveFacts(Ticker, *CompanyFacts) error
    LoadFacts(Ticker) (*CompanyFacts, error)
}
```

The default behaviour of `sec.NewClient()` is to always go to the source (SEC). To create a client with persistence use `sec.NewClientWithPersistence()`

A basic persistence layer using [bbolt](https://github.com/etcd-io/bbolt) is provided out of the box. For example: 
```go
persistenceLayerConfig := BoltPersistenceLayerConfig{Path: tempFile.Name(), ExpiresAfter: time.Hour}
bolt, err = NewBoltPersistenceLayer(persistenceLayerConfig)
if err != nil {
    panic(err)
}

clientWithPersistenceLayer := sec.NewClientWithPersistence(bolt)
```

### Bulk data
Getting filing data from all companies individually would be very time consuming with the rate limit set by SEC. Becasue of that a bulk function is provided that gets fetches **ALL** company data in one go. `client.GetBulk()` returns a channel that will eventually contain all the data. 

The latency and throughput of the bulk function is highly bandwidth and processing speed dependent. You can expect that it will take ~1 minute or so before the first ticker fact pair appears on the channel. 
```go
client := NewClient()
factChannel, err := client.GetBulk()
if err != nil {
    panic(err)
}

for tickerFact := range factChannel {
    fmt.Println(tickerFact)
}
```

## Compliance
Coming soon
