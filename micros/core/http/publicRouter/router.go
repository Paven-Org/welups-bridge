package publicRouter

import (
	"bridge/common"
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/config"
	log "bridge/service-managers/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Config(router gin.IRouter) {
	initialize()
	gr := router.Group("/p")
	gr.POST("/login", loginHandler)
	gr.POST("/logout" /*,authenMW*/, logoutHandler)
	gr.POST("/passwd" /*,authenMW*/, passwdHandler)
}

type loginReq struct {
	Username string
	Password string
}

var logger *zerolog.Logger
var serverCnf common.HttpConf

func initialize() {
	logger = log.Get()
	serverCnf = config.Get().HttpConfig
	logger.Info().Msg("public handlers initialized")
}

func loginHandler(c *gin.Context) {
	// request
	var lReq = loginReq{}
	if err := c.ShouldBindJSON(&lReq); err != nil {
		logger.Err(err).Msgf("[login handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// process
	token, sessionID, sessionSecret, dur, err := userLogic.Login(lReq.Username, lReq.Password)
	if err != nil {
		logger.Err(err).Msgf("[login handler] Unable to create session")
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to login with username=%s, error: %s", lReq.Username, err.Error()))
		return
	}
	logger.Debug().Msgf("[login handler] token %s, sessionID %s, sessionSecret %s, age: %d", token, sessionID, sessionSecret, dur)

	// response
	c.SetCookie(sessionID, sessionSecret, int(dur), "/", "", serverCnf.Mode == "prod", true)
	c.JSON(http.StatusOK, gin.H{"token": token})
	return
}

func logoutHandler(c *gin.Context) {
	// request
	tokenS := c.GetHeader("Authorization")
	tokenS = strings.TrimPrefix(tokenS, "Bearer ")
	logger.Debug().Msg("Token string: " + tokenS)

	claims, err := userLogic.ParseTokenToClaims(tokenS)
	if err != nil {
		logger.Err(err).Msgf("Unable to parse JWT")
		c.JSON(http.StatusBadRequest, "Unable to parse JWT")
		return
	}
	sessionID := claims.Session

	cookie, err := c.Cookie(sessionID)
	if err != nil {
		logger.Err(err).Msgf("No cookie in request")
		c.JSON(http.StatusBadRequest, "No cookie in request")
		return
	}

	// process
	if err = userLogic.Logout(tokenS, cookie); err != nil {
		logger.Err(err).Msg("Unable to logout user " + claims.Username)
		c.JSON(http.StatusBadRequest, "Unable to logout user "+claims.Username)
		return
	}

	c.SetCookie(sessionID, "", -1, "/", serverCnf.Host, serverCnf.Mode == "prod", true)

	// response
	c.JSON(http.StatusOK, fmt.Sprintf("User %s logout successfully", claims.Username))
	return
}

func passwdHandler(c *gin.Context) {

}
