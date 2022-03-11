package publicRouter

import (
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/config"
	log "bridge/service-managers/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Config(router gin.IRouter) {
	gr := router.Group("/p")
	gr.POST("/login", loginHandler)
	gr.POST("/logout" /*,authenMW*/, logoutHandler)
	gr.POST("/passwd" /*,authenMW*/, passwdHandler)
}

type loginReq struct {
	Username string
	Password string
}

var logger = log.Get()
var serverCnf = config.Get().HttpConfig

func loginHandler(c *gin.Context) {
	var lReq = loginReq{}
	if err := c.ShouldBindJSON(&lReq); err != nil {
		logger.Err(err).Msgf("[login handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, sessionID, sessionSecret, dur, err := userLogic.Login(lReq.Username, lReq.Password)
	if err != nil {
		logger.Err(err).Msgf("[login handler] Unable to create session")
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to login with username=%s, error: %s", lReq.Username, err.Error()))
		return
	}
	logger.Debug().Msgf("[login handler] token %s, sessionID %s, sessionSecret %s, age: %d", token, sessionID, sessionSecret, dur)
	c.SetCookie(sessionID, sessionSecret, int(dur), "/", serverCnf.Host, serverCnf.Mode == "prod", true)

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func logoutHandler(c *gin.Context) {

}

func passwdHandler(c *gin.Context) {

}
