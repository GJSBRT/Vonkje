package power_prices

import (
	"fmt"
	"time"
	"context"

	"gijs.eu/vonkje/packages/victoria_metrics"

	"github.com/sirupsen/logrus"
)

type PowerPriceSourceConfigs struct {
	AllInPower AllInPowerConfig `mapstructure:"all-in-power"`
	Entsoe EntsoeConfig `mapstructure:"entsoe"`
}

type Config struct {
	Enabled bool `mapstructure:"enabled"`
	Sources PowerPriceSourceConfigs `mapstructure:"sources"`
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
	if !config.Sources.AllInPower.Enable {
		sources = append(sources, newAllInPower(config.Sources.AllInPower))
	}

	if !config.Sources.Entsoe.Enable {
		sources = append(sources, newEntsoe(config.Sources.Entsoe))
	}

	return &PowerPrices{
		Config: config,
		errChannel: errChannel,
		ctx: ctx,
		logger: logger,
		VictoriaMetrics: victoriaMetrics,
	}
}

func (pp *PowerPrices) addPricesOfSource(source PowerPriceSource) error {
	prices, err := source.GetPricesKwH(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC))
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
	if !pp.Config.Enabled {
		pp.logger.Warn("Power price collector is disabled")
		return
	}

	pp.logger.Info("Starting power price collector")

	err := pp.updateMetrics()
	if err != nil {
		pp.errChannel <- err
	}
	now := time.Now()

    nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 0, 0, 0, time.Local)
    durationUntilMidnight := nextMidnight.Sub(now)

    time.Sleep(durationUntilMidnight)

    go pp.Start()
}
