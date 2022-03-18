package admRouter

import (
	"bridge/micros/core/http/admRouter/manageUserRouter"
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
}
