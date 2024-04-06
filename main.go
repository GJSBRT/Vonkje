package main

import (
	"os"
	"flag"
	"syscall"
	"context"
	"os/signal"

	"gijs.eu/vonkje/http"
	"gijs.eu/vonkje/modbus"
	"gijs.eu/vonkje/control"
	"gijs.eu/vonkje/power_prices"
	"gijs.eu/vonkje/packages/victoria_metrics"

	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel 			string `mapstructure:"log-level"`
	HTTP 				http.Config `mapstructure:"http"`
	Modbus 				modbus.Config `mapstructure:"modbus"`
	VictoriaMetrics 	victoria_metrics.Config `mapstructure:"victoria-metrics"`
	PowerPrices 		power_prices.Config `mapstructure:"power-prices"`
	Control 			control.Config `mapstructure:"control"`
}

var (
	config Config
	logger	= logrus.New()
	errChannel = make(chan error)
)

func init() {
	logger.Info("Starting...")
	configPath := flag.String("config", "config.yaml", "Full path to your config file")
	flag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	
	logger.Info("Read config file successfully")

	switch config.LogLevel {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	config = Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		logger.WithError(err).Panic("Failed to unmarshal config")
	}
}

func main() {
	go func () {
		for {
			select {
			case err := <-errChannel:
				logger.WithError(err).Error("An error occurred")
			}
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	modbusClient, err := modbus.New(config.Modbus, errChannel, stopCtx, logger)
	if err != nil {
		logger.WithError(err).Panic("Failed to create modbus client")
	}
	go modbusClient.Start()

	httpServer := http.New(config.HTTP, errChannel, stopCtx, logger)
	go httpServer.Start()

	victoriaMetricsClient := victoria_metrics.New(config.VictoriaMetrics)

	powerPricesClient := power_prices.New(config.PowerPrices, errChannel, stopCtx, logger, victoriaMetricsClient)
	go powerPricesClient.Start()

	controlClient := control.New(config.Control, errChannel, stopCtx, logger, victoriaMetricsClient, modbusClient)
	go controlClient.Start()

	<-stopCtx.Done()

	modbusClient.Close()

	logger.Info("Exited")
}
