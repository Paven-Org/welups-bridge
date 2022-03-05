package admRouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Config(router gin.IRouter) {
	gr := router.Group("/api/v1" /*,middlewares.Author*/)

	//gr.GET("/someRoute", someHandler)
	gr.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	// subRouters
	//eth2wel.Config(gr)
}
