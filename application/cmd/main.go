package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/ExchangeDiary/exchange-diary/application/route"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/persistence"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func main() {
	logger := setLogger()
	server := bootstrap(logger)
	server.Run(":8080") // TODO: viper
	shutdown(logger)
}

func setLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	return logger
}

func bootstrap(logger *zap.Logger) *gin.Engine {
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

	// init db
	db := infrastructure.ConnectDatabase(configName)
	infrastructure.Migrate(db)

	// set DI
	roomRepository := persistence.NewRoomRepository(db)
	roomMemberRepository := persistence.NewRoomMemberRepository(db)
	roomMemberService := service.NewRoomMemberService(roomMemberRepository)
	roomService := service.NewRoomService(roomRepository, roomMemberService)
	roomController := controller.NewRoomController(roomService)

	// init server
	server := gin.New()

	// zap middlewares
	server.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(logger, true)) // log all panic

	// init routes
	route.RoomRoutes(server, roomController)
	return server
}

func shutdown(logger *zap.Logger) {
	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
