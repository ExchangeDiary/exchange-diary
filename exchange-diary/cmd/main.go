package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ExchangeDiary_Server/exchange-diary/configs"

	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	defaultConfig := "dev"

	var configName string
	flag.StringVar(&configName, "phase", defaultConfig, "name of configuration file with no extension")
	flag.Parse()

	_, err := config.Load("./configs", configName)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config file: %s", err.Error()))
	}
	// logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
