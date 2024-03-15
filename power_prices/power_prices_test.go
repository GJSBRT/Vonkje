package power_prices

// import (
// 	"time"
// 	"testing"
// )

// func Test_Entsoe_GetPricesKwH(t *testing.T) {
// 	entsoeClient := newEntsoe(EntsoeConfig{
// 		Domain: "",
// 		SecurityToken: "",
// 	})

// 	prices, err := entsoeClient.GetPricesKwH(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC))
// 	if err != nil {
// 		t.Fatalf("Error getting prices: %s", err)
// 	}

// 	t.Log(prices)
// }