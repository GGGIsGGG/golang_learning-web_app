package routes

import (
	"net/http"
	"web_app/logger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK!!")
	})

	return r
}
