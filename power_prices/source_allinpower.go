package power_prices

import (
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type AllInPowerConfig struct {}

type AllInPower struct {
	name string
	Config AllInPowerConfig
	Resolution uint
	Client *http.Client
}

func newAllInPower(config AllInPowerConfig) *AllInPower {
	return &AllInPower{
		name: "all-in-power",
		Config: config,
		Resolution: 60,
		Client: &http.Client{},
	}
}

func (a *AllInPower) GetName() string {
	return a.name
}

type spotMarketPriceResponse struct {
	Id int `json:"id"`
	Timestamps []string `json:"timestamps"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	ProductType string `json:"product_type"`
	Date string `json:"date"`
	Unit string `json:"unit"`
	Prices []float64 `json:"prices"`
}

func (a *AllInPower) GetPricesKwH(timestamp time.Time) (map[time.Time]float64, error) {
	request, err := http.NewRequest("GET", "https://api.allinpower.nl/troodon/api/p/spot_market/prices/?date="+timestamp.Format(time.DateOnly)+"&product_type=ELK", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Del("User-Agent")

	response, err := a.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ErrFailedToRetrieveData
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseBody spotMarketPriceResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	prices := make(map[time.Time]float64)

	for i, price := range responseBody.Prices {
		timestamp, err := time.Parse("2006-01-02T15:04:05.000000Z", responseBody.Timestamps[i])
		if err != nil {
			return nil, err
		}

		prices[timestamp] = price
	}

	return prices, nil
}
