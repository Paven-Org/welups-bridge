package versionRouters

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MkVRouter(version string, parentRouter gin.IRouter) gin.IRouter {
	gr := parentRouter.Group("/" + version /*,middlewares.Author*/)

	//gr.GET("/someRoute", someHandler)
	gr.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	// subRouters
	//eth2wel.Config(gr)

	return gr
}
