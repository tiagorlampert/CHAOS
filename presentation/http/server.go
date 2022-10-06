package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/utils/template"
	"net/http"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Static("/static", "web/static")
	router.HTMLRender = template.LoadTemplates("web")
	return router
}

func NewServer(router *gin.Engine, configuration *environment.Configuration) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", configuration.Server.Port),
		http.TimeoutHandler(router, internal.TimeoutDuration, internal.TimeoutExceeded))
}
