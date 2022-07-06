package userRouter

import (
	"bridge/common"
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/config"
	"bridge/micros/core/model"
	log "bridge/service-managers/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// router for users internal to the bridge system, e.g. admin, service manager etc...
func Config(router gin.IRouter, authMW gin.HandlerFunc) {
	initialize()
	gr := router.Group("/u")
	gr.POST("/login", loginHandler)
	gr.POST("/logout", logoutHandler)
	gr.POST("/passwd", authMW, passwdHandler)
	gr.POST("/update", authMW, userUpdateHandler)
	gr.GET("/:username", getUserHandler)
	gr.GET("/myroles", authMW, getCurrentUserRoles)
}

var logger *zerolog.Logger
var serverCnf common.HttpConf

func initialize() {
	logger = log.Get()
	serverCnf = config.Get().HttpConfig
	logger.Info().Msg("user handlers initialized")
}

func loginHandler(c *gin.Context) {
	// request
	type loginReq struct {
		Username string
		Password string
	}

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
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(sessionID, sessionSecret, int(dur), "/", "", true, true)
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

	c.SetCookie(sessionID, "", -1, "/", serverCnf.Host, true, true)

	// response
	c.JSON(http.StatusOK, fmt.Sprintf("User %s logout successfully", claims.Username))
	return
}

func passwdHandler(c *gin.Context) {
	// request

	type passwdRequest struct {
		OldPasswd string `json:"old_passwd"`
		NewPasswd string `json:"new_passwd"`
	}

	var req passwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[passwd handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// process
	username := c.GetString("username")
	if err := userLogic.Passwd(username, req.OldPasswd, req.NewPasswd); err != nil {
		if err == model.ErrWrongPasswd {
			logger.Err(err).Msgf("[passwd handler] wrong old password")
			c.JSON(http.StatusBadRequest, "Wrong old password")
			return
		}
		logger.Err(err).Msgf("[passwd handler] Unable to update password")
		c.JSON(http.StatusInternalServerError, "Unable to update password")
		return
	}

	// response
	logger.Info().Msgf("[passwd handler] Password updated successfully")
	c.JSON(http.StatusOK, "Password updated successfully")
	return
}

func userUpdateHandler(c *gin.Context) {

	// request

	// users are only allowed to update email for now
	type updateRequest struct {
		NewEmail string `json:"new_email"`
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[update handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}
	req.NewEmail = strings.TrimSpace(req.NewEmail)

	// process
	username := c.GetString("username")
	if err := userLogic.UpdateUserInfo(username, req.NewEmail); err != nil {
		logger.Err(err).Msgf("[update handler] Unable to update user")
		c.JSON(http.StatusInternalServerError, "Unable to update user")
		return
	}

	// response
	logger.Info().Msgf("[update handler] User updated successfully")
	c.JSON(http.StatusOK, "User updated successfully")
	return
}

func getUserHandler(c *gin.Context) {
	// request
	username := c.Param("username")

	// process
	user, err := userLogic.GetUserByName(username)
	if err != nil {
		logger.Err(err).Msgf("[getUserHandler] Unable to retrieve user")
		status := http.StatusInternalServerError
		if err == model.ErrEthAccountNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to retrieve user "+username)
		return
	}

	// response
	user.Password = "" // just to be sure, this field wouldn't be marshalled anyway
	c.JSON(http.StatusOK, user)
	return
}

func getCurrentUserRoles(c *gin.Context) {
	username := c.GetString("username")

	// process
	roles, err := userLogic.GetUserRoles(username)
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrRoleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get user "+username+"'s roles")
		return
	}

	// response

	c.JSON(http.StatusOK, roles)
	return
}
