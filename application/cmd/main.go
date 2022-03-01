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
	"github.com/ExchangeDiary/exchange-diary/application/middleware"
	"github.com/ExchangeDiary/exchange-diary/application/route"
	"github.com/ExchangeDiary/exchange-diary/docs"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/persistence"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// refs: https://github.com/swaggo/swag/blob/master/example/celler/main.go
// @title           Voice Of Diary API (voda)
// @version         1.0
// @description     This is a voda api server.

// @contact.name   API Support
// @contact.url    https://minkj1992.github.io
// @contact.email  minkj1992@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      exchange-diary-b4mzhzbzcq-du.a.run.app
// 로컬 테스트용 host      localhost:8080
// @BasePath  /v1

const (
	//
	versionPrefix = "/v1"
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
	db := infrastructure.ConnectDatabase(configName)
	infrastructure.Migrate(db)

	// set DI
	roomRepository := persistence.NewRoomRepository(db)
	roomMemberRepository := persistence.NewRoomMemberRepository(db)
	memberRepository := persistence.NewMemberRepository(db)

	roomMemberService := service.NewRoomMemberService(roomMemberRepository)
	memberService := service.NewMemberService(memberRepository)
	roomService := service.NewRoomService(roomRepository, memberRepository, roomMemberService)
	authCodeVerifier := service.NewTokenVerifier(service.AuthCodeSecretKey)
	refreshTokenVerifier := service.NewTokenVerifier(service.AccessTokenSecretKey)
	tokenService := service.NewTokenService(memberService, authCodeVerifier, refreshTokenVerifier)

	authController := controller.NewAuthController(conf.Client, memberService, tokenService)
	tokenController := controller.NewTokenController(tokenService)
	roomController := controller.NewRoomController(roomService)

	authenticationFilter := middleware.NewAuthenticationFilter(authCodeVerifier)

	// init server
	server := gin.New()

	// set swagger
	swagger(server)

	// zap middlewares
	server.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(logger, true)) // log all panic

	// init routes
	v1 := server.Group(versionPrefix)
	route.AuthRoutes(v1, authController)
	route.TokenRoutes(v1, tokenController)

	v1.Use(authenticationFilter.Authenticate())
	route.RoomRoutes(v1, roomController)

	return server
}

func swagger(server *gin.Engine) {
	docs.SwaggerInfo.BasePath = versionPrefix
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func shutdown(logger *zap.Logger) {
	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
