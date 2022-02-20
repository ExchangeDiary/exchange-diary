package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ExchangeDiary_Server/exchange-diary/application/route"
	"github.com/ExchangeDiary_Server/exchange-diary/infrastructure/configs"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	defaultConfig := "dev"
	configPath := "./infrastructure/configs" // TODO: Dockerfile path

	var configName string
	flag.StringVar(&configName, "phase", defaultConfig, "name of configuration file with no extension")
	flag.Parse()
	
	_, err := configs.Load(configPath, configName)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config file: %s", err.Error()))
	}
	
	run()
}

func run() {
	server := gin.New()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	
	// zap middlewares
	server.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(logger, true)) // log all panic

	// TODO: need DI

	// routes
	route.RoomRoutes(server)

	server.Run(":8080") // TODO: viper
	
	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
