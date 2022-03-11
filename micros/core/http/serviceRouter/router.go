package admRouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	gr := router.Group("/s", mw... /*,middlewares.Author*/)

	//gr.GET("/someRoute", someHandler)
	gr.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	gr.GET("/pkey/:chain", getPkeyHandler)
}

func getPkeyHandler(c *gin.Context) {

}
