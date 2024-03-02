package main

import (
	"os"
	"flag"
	"syscall"
	"context"
	"os/signal"

	"gijs.eu/huawei-modbus/http"
	"gijs.eu/huawei-modbus/modbus"

	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel 	string `mapstructure:"log-level"`
	HTTP 		http.Config `mapstructure:"http"`
	Modbus 		modbus.Config `mapstructure:"modbus"`
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

	<-stopCtx.Done()

	modbusClient.Close()

	logger.Info("Exited")
}
