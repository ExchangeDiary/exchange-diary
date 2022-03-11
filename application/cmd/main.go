package main

import (
	"flag"
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
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/cloudstorage"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/task"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/persistence"
	"github.com/spf13/viper"

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
// @BasePath  /v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

const (
	//
	versionPrefix = "/v1"
	defaultPhase  = "dev"
	configPath    = "./infrastructure/configs"
)

var (
	phase string
	conf  configs.Config
)

func main() {
	var err error
	flag.StringVar(&phase, "phase", defaultPhase, "name of configuration file with no extension")
	flag.Parse()
	viper.SetDefault("PHASE", phase)

	conf, err = configs.Load(configPath)
	if err != nil {
		panic("Failed to load config file: " + err.Error())
	}
	logger.Info("cold start google cloud storage client")
	storageClient := cloudstorage.GetClient()
	defer storageClient.Close()

	logger.Info("cold start google cloud tasks client")
	taskClient := task.GetClient()
	defer taskClient.Close()

	logger.Info("cold start application")
	server := bootstrap()
	server.Run(":8080") // TODO: viper
	shutdown()
}

func bootstrap() *gin.Engine {
	// init db
	db := infrastructure.ConnectDatabase(phase)
	infrastructure.Migrate(db)

	// set DI
	roomRepository := persistence.NewRoomRepository(db)
	roomMemberRepository := persistence.NewRoomMemberRepository(db)
	memberRepository := persistence.NewMemberRepository(db)

	roomMemberService := service.NewRoomMemberService(roomMemberRepository, memberRepository)
	memberService := service.NewMemberService(memberRepository)
	roomService := service.NewRoomService(roomRepository, roomMemberService)
	authCodeVerifier := service.NewTokenVerifier(service.AuthCodeSecretKey)
	refreshTokenVerifier := service.NewTokenVerifier(service.AccessTokenSecretKey)
	tokenService := service.NewTokenService(memberService, authCodeVerifier, refreshTokenVerifier)
	fileService := service.NewFileService()
	alarmService := service.NewAlarmService()
	taskService := service.NewTaskService(alarmService, roomRepository)

	memberController := controller.NewMemberController(memberService)
	authController := controller.NewAuthController(conf.Client, memberService, tokenService)
	tokenController := controller.NewTokenController(tokenService)
	roomController := controller.NewRoomController(roomService)
	fileController := controller.NewFileController(fileService)
	taskController := controller.NewTaskController(taskService, memberService)

	authenticationFilter := middleware.NewAuthenticationFilter(authCodeVerifier)

	// init server
	server := gin.New()

	// set swagger
	swagger(server)

	// zap middlewares
	server.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	// server.Use(ginzap.RecoveryWithZap(logger.Log, true)) // log all panic

	// init routes
	v1 := server.Group(versionPrefix)
	route.AuthRoutes(v1, authController)
	route.TokenRoutes(v1, tokenController)

	v1.Use(authenticationFilter.Authenticate())

	route.RoomRoutes(v1, roomController)
	route.MemberRoutes(v1, memberController)
	route.FileRoutes(v1, fileController)
	route.TaskRoutes(v1, taskController)

	return server
}

func swagger(server *gin.Engine) {
	docs.SwaggerInfo.BasePath = versionPrefix
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func shutdown() {
	// Wait for termination signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-c
	logger.Info("Application terminates", zap.Any("Signal", osSignal))
}
