package bridgeRouter

import (
	"bridge/micros/core/http/bridgeRouter/welethRouter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	gr := router.Group("/b", mw...)

	//gr.GET("/someRoute", someHandler)
	gr.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	// subRouters
	welethRouter.Config(gr)
}
