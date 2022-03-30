package admRouter

import (
	ethRouter "bridge/micros/core/http/admRouter/eth-router"
	"bridge/micros/core/http/admRouter/manageUserRouter"
	welRouter "bridge/micros/core/http/admRouter/wel-router"
	"bridge/micros/core/http/bridgeRouter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	gr := router.Group("/a", mw... /*,middlewares.Author*/)

	//gr.GET("/someRoute", someHandler)
	gr.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	// subRouters
	manageUserRouter.Config(gr)
	ethRouter.Config(gr)
	welRouter.Config(gr)
	bridgeRouter.Config(gr)
}
