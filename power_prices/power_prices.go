package power_prices

import (
	"fmt"
	"time"
	"context"

	"gijs.eu/vonkje/packages/victoria_metrics"

	"github.com/sirupsen/logrus"
)

type PowerPriceSourceConfigs struct {
	AllInPowerConfig AllInPowerConfig `json:"all-in-power"`
}

type Config struct {
	Sources PowerPriceSourceConfigs `json:"sources"`
}

type PowerPrices struct {
	Config Config
	errChannel chan error
	ctx context.Context
	logger *logrus.Logger
	VictoriaMetrics *victoria_metrics.VictoriaMetrics
}

type PowerPriceSource interface {
	GetName() string
	GetPricesKwH(time.Time) (map[time.Time]float64, error)
}

var sources = []PowerPriceSource{}

var (
	ErrFailedToAuthenticate = fmt.Errorf("Failed to authenticate")
	ErrFailedToRetrieveData = fmt.Errorf("Failed to retrieve data")
)

func New(
	config Config, 
	errChannel chan error,
	ctx context.Context,
	logger *logrus.Logger,
	victoriaMetrics *victoria_metrics.VictoriaMetrics,
) *PowerPrices {
	sources = append(sources, newAllInPower(config.Sources.AllInPowerConfig))

	return &PowerPrices{
		Config: config,
		errChannel: errChannel,
		ctx: ctx,
		logger: logger,
		VictoriaMetrics: victoriaMetrics,
	}
}

func (pp *PowerPrices) addPricesOfSource(source PowerPriceSource) error {
	prices, err := source.GetPricesKwH(time.Now())
	if err != nil {
		return err
	}

	metrics := []victoria_metrics.VictoriaMetricsRequest{}
	for timestamp, price := range prices {
		metrics = append(metrics, victoria_metrics.VictoriaMetricsRequest{
			Metric: map[string]string{
				"__name__": "power_price",
				"source": source.GetName(),
			},
			Values: []float64{price},
			Timestamps: []int64{int64(timestamp.Unix() * 1000)},
		})
	}

	err = pp.VictoriaMetrics.SendMetrics(metrics)
	if err != nil {
		return err
	}

	return nil
}

func (pp *PowerPrices) updateMetrics() error {
	pp.logger.Info("Updating power prices")

	for _, source := range sources {
		err := pp.addPricesOfSource(source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pp *PowerPrices) Start() {
	pp.logger.Info("Starting power price collector")

	now := time.Now()

    nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 0, 0, 0, time.Local)
    durationUntilMidnight := nextMidnight.Sub(now)

    time.Sleep(durationUntilMidnight)

	err := pp.updateMetrics()
	if err != nil {
		pp.errChannel <- err
	}

    go pp.Start()
}
