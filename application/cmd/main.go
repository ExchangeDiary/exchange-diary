package main

import (
	"flag"
	"fmt"
	"github.com/exchange-diary/domain/service"
	"github.com/exchange-diary/infrastructure"
	"github.com/exchange-diary/infrastructure/persistence"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exchange-diary/application/controller"
	"github.com/exchange-diary/application/route"
	"github.com/exchange-diary/infrastructure/configs"

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

	conf, err := configs.Load(configPath, configName)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config file: %s", err.Error()))
	}

	// init db
	db := infrastructure.ConnectDatabase()
	infrastructure.Migrate(db)

	// set DI
	roomRepository := persistence.NewRoomRepository(db)
	roomMemberRepository := persistence.NewRoomMemberRepository(db)
	memberRepository := persistence.NewMemberRepository(db)

	roomMemberService := service.NewRoomMemberService(roomMemberRepository)
	roomService := service.NewRoomService(roomRepository, roomMemberService)
	memberService := service.NewMemberService(memberRepository)
	authCodeVerifier := service.NewTokenVerifier(service.AUTH_CODE_SECRET_KEY)
	refreshTokenVerifier := service.NewTokenVerifier(service.ACCESS_TOKEN_SECRET_KEY)
	tokenService := service.NewTokenService(memberService, authCodeVerifier, refreshTokenVerifier)

	authController := controller.NewAuthController(conf.Client, memberService, tokenService)
	tokenController := controller.NewTokenController(tokenService)
	roomController := controller.NewRoomController(roomService)

	// init server
	server := gin.New()

	// zap middlewares
	server.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(logger, true)) // log all panic

	// init routes
	v1 := server.Group("/api/v1")
	route.RoomRoutes(server, roomController)
	route.AuthRoutes(v1, authController)
	route.TokenRoutes(v1, tokenController)
	return server
}

func shutdown(logger *zap.Logger) {
	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
