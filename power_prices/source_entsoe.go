package power_prices

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/xml"
)

const entsoeTimestampFormat = "200601021504"

type EntsoeConfig struct {
	SecurityToken string `mapstructure:"security-token"`
	Domain string `mapstructure:"domain"`
}

type Entsoe struct {
	name string
	Client *http.Client
	Config EntsoeConfig
}

func newEntsoe(config EntsoeConfig) *Entsoe {
	return &Entsoe{
		name: "entsoe",
		Client: &http.Client{},
		Config: config,
	}
}

func (e *Entsoe) GetName() string {
	return e.name
}

type entsoeDayAheadResponse struct {
	MRID                        string   `xml:"mRID"`
	RevisionNumber              string   `xml:"revisionNumber"`
	Type                        string   `xml:"type"`
	SenderMarketParticipantMRID struct {
		Text         string `xml:",chardata"`
		CodingScheme string `xml:"codingScheme,attr"`
	} `xml:"sender_MarketParticipant.mRID"`
	SenderMarketParticipantMarketRoleType string `xml:"sender_MarketParticipant.marketRole.type"`
	ReceiverMarketParticipantMRID         struct {
		Text         string `xml:",chardata"`
		CodingScheme string `xml:"codingScheme,attr"`
	} `xml:"receiver_MarketParticipant.mRID"`
	ReceiverMarketParticipantMarketRoleType string `xml:"receiver_MarketParticipant.marketRole.type"`
	CreatedDateTime                         string `xml:"createdDateTime"`
	PeriodTimeInterval                      struct {
		Text  string `xml:",chardata"`
		Start string `xml:"start"`
		End   string `xml:"end"`
	} `xml:"period.timeInterval"`
	TimeSeries []struct {
		Text         string `xml:",chardata"`
		MRID         string `xml:"mRID"`
		BusinessType string `xml:"businessType"`
		InDomainMRID struct {
			Text         string `xml:",chardata"`
			CodingScheme string `xml:"codingScheme,attr"`
		} `xml:"in_Domain.mRID"`
		OutDomainMRID struct {
			Text         string `xml:",chardata"`
			CodingScheme string `xml:"codingScheme,attr"`
		} `xml:"out_Domain.mRID"`
		CurrencyUnitName     string `xml:"currency_Unit.name"`
		PriceMeasureUnitName string `xml:"price_Measure_Unit.name"`
		CurveType            string `xml:"curveType"`
		Period               struct {
			Text         string `xml:",chardata"`
			TimeInterval struct {
				Text  string `xml:",chardata"`
				Start string `xml:"start"`
				End   string `xml:"end"`
			} `xml:"timeInterval"`
			Resolution string `xml:"resolution"`
			Point      []struct {
				Text        string `xml:",chardata"`
				Position    string `xml:"position"`
				PriceAmount string `xml:"price.amount"`
			} `xml:"Point"`
		} `xml:"Period"`
	} `xml:"TimeSeries"`
} 

func (e *Entsoe) GetPricesKwH(timestamp time.Time) (map[time.Time]float64, error) {
	startTimestamp := timestamp.Format(entsoeTimestampFormat)
	endTimestamp := timestamp.Add(24 * time.Hour).Format(entsoeTimestampFormat)

	requestUrl := fmt.Sprintf("https://web-api.tp.entsoe.eu/api?securityToken=%s&documentType=A44&in_Domain=%s&out_Domain=%s&periodStart=%s&periodEnd=%s", e.Config.SecurityToken, e.Config.Domain, e.Config.Domain, startTimestamp, endTimestamp)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Del("User-Agent")

	response, err := e.Client.Do(request)
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

	var responseBody entsoeDayAheadResponse
	err = xml.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	prices := make(map[time.Time]float64)
	startingTimestamp, err := time.Parse("2006-01-02T15:04Z", responseBody.TimeSeries[1].Period.TimeInterval.Start)
	if err != nil {
		return nil, err
	}

	for i, entry := range responseBody.TimeSeries[1].Period.Point {
		timestamp := startingTimestamp.Add(time.Duration(i) * time.Hour)

		price, err := strconv.ParseFloat(entry.PriceAmount, 2)
		if err != nil {
			return nil, err
		}

		prices[timestamp] = price / 1000
	}

	return prices, nil
}
