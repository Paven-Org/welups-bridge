package manageUserRouter

import (
	log "bridge/service-managers/logger"
	"net/http"
	"strconv"

	userLogic "bridge/micros/core/blogic/user"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/m/eth", mw... /*,middlewares.Author*/)
	//AddEthAccount(address string, status string)
	//GetAllEthAccounts(offset uint, size uint)
	//GetAllRoles()
	//GetEthAccount(address string)
	//GetEthAccountRoles(address string)
	//GetEthAccountsWithRole(role string, offset uint, size uint)
	//GetEthPrikeyIfExists(address string)
	//GrantRole(address string, role string)
	//RemoveEthAccount(address string)
	//RevokeRole(address string, role string)
	//SetEthAccountStatus(address string, status string)
	//SetPriKey(address string, key string)
	//UnsetPrikey(address string)
	gr.POST("/add", addEthAccount)
	gr.POST("/set-status/:acc/:status", setStatus)
	gr.POST("/set-prikey/:acc", setKey)
	gr.POST("/unset-prikey/:acc", unsetKey)
	gr.POST("/remove/:acc", removeEth)
	gr.POST("/grant/:role/to/:acc", grantRole)
	gr.POST("/revoke/:role/from/:acc", revokeRole)
	gr.POST("/getroles/:acc", getUserRoles)
	gr.GET("/haverole/:role/:page", getAccsWithRole)
	gr.GET("/accounts/:page", getAccs)
	gr.GET("/info/:account", getAcc)
	gr.GET("/roles", getRoles)

}

func initialize() {
	logger = log.Get()
	logger.Info().Msg("manage users handlers initialized")
}

func getAccsWithRole(c *gin.Context) {
	// request
	var page, limit uint64

	role := c.Param("role")

	_page := c.Param("page")
	page, err := strconv.ParseUint(_page, 10, 32)
	if err != nil {
		logger.Err(err).Msgf("[get users with role handler] invalid page")
		page = 1 // default
		return
	}

	_limit := c.Query("limit")
	if _limit == "" {
		limit = 10 // default
	} else {
		limit, err = strconv.ParseUint(_limit, 10, 32)
		if err != nil {
			logger.Err(err).Msgf("[get users with role handler] invalid limit")
			limit = 10 // default
			return
		}

	}

	// process
	users, err := userLogic.GetUsersWithRole(role, uint((page-1)*limit), uint(limit))
	if err != nil {
		logger.Err(err).Msgf("[get users with role handler] Unable to get users")
		c.JSON(http.StatusInternalServerError, "Unable to get users with role "+role)
		return
	}

	// response

	logger.Info().Msgf("[get users with role handler] Get users with role %s successfully", role)
	c.JSON(http.StatusOK, &users)
	return
}

func getAccs(c *gin.Context) {
	// request
	var page, limit uint64

	_page := c.Param("page")
	page, err := strconv.ParseUint(_page, 10, 32)
	if err != nil {
		logger.Err(err).Msgf("[get users handler] invalid page")
		page = 1 // default
		return
	}

	_limit := c.Query("limit")
	if _limit == "" {
		limit = 10 // default
	} else {
		limit, err = strconv.ParseUint(_limit, 10, 32)
		if err != nil {
			logger.Err(err).Msgf("[get users handler] invalid limit")
			limit = 10 // default
			return
		}

	}

	// process
	users, err := userLogic.GetUsers(uint((page-1)*limit), uint(limit))
	if err != nil {
		logger.Err(err).Msgf("[get users handler] Unable to get users")
		c.JSON(http.StatusInternalServerError, "Unable to get users")
		return
	}

	// response

	logger.Info().Msgf("[get users handler] Get users %s successfully")
	c.JSON(http.StatusOK, &users)
	return
}

func getRoles(c *gin.Context) {
	// request

	// process
	roles, err := userLogic.GetAllRoles()
	if err != nil {
		logger.Err(err).Msgf("[get roles handler] Unable to get roles")
		c.JSON(http.StatusInternalServerError, "Unable to get roles")
		return
	}

	// response
	logger.Info().Msgf("[get roles handler] get roles successfully")
	c.JSON(http.StatusOK, &roles)
	return

}

func addEthAccount(c *gin.Context) {
	type addUserReq struct {
		Username string
		Email    string
		Password string
	}

	var req addUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[add user handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := userLogic.AddUser(req.Username, req.Email, req.Password); err != nil {
		logger.Err(err).Msgf("[add user handler] Unable to add user")
		c.JSON(http.StatusInternalServerError, "Unable to add user")
		return
	}

	logger.Info().Msgf("[add user handler] User successfully added")
	c.JSON(http.StatusOK, "User successfully added")
	return
}

func setStatus(c *gin.Context) {
	// request
	type updateRequest struct {
		Username    string `json:"username"`
		NewUsername string `json:"new_username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Status      string `json:"status"`
	}
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[update handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}
	// process
	if err := userLogic.AdminUpdateUserInfo(req.Username, req.NewUsername, req.Email, req.Password, req.Status); err != nil {
		logger.Err(err).Msgf("[update handler] Unable to update user")
		c.JSON(http.StatusInternalServerError, "Unable to update user")
		return
	}

	// response
	logger.Info().Msgf("[update handler] User updated successfully")
	c.JSON(http.StatusOK, "User updated successfully")
	return

}

func setKey(c *gin.Context) {
	// request
	type updateRequest struct {
		Username    string `json:"username"`
		NewUsername string `json:"new_username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Status      string `json:"status"`
	}
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[update handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}
	// process
	if err := userLogic.AdminUpdateUserInfo(req.Username, req.NewUsername, req.Email, req.Password, req.Status); err != nil {
		logger.Err(err).Msgf("[update handler] Unable to update user")
		c.JSON(http.StatusInternalServerError, "Unable to update user")
		return
	}

	// response
	logger.Info().Msgf("[update handler] User updated successfully")
	c.JSON(http.StatusOK, "User updated successfully")
	return

}

func unsetKey(c *gin.Context) {
	// request
	type updateRequest struct {
		Username    string `json:"username"`
		NewUsername string `json:"new_username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Status      string `json:"status"`
	}
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[update handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}
	// process
	if err := userLogic.AdminUpdateUserInfo(req.Username, req.NewUsername, req.Email, req.Password, req.Status); err != nil {
		logger.Err(err).Msgf("[update handler] Unable to update user")
		c.JSON(http.StatusInternalServerError, "Unable to update user")
		return
	}

	// response
	logger.Info().Msgf("[update handler] User updated successfully")
	c.JSON(http.StatusOK, "User updated successfully")
	return

}

func removeEth(c *gin.Context) {
	// request
	username := c.Param("user")

	// process
	if err := userLogic.RemoveUser(username); err != nil {
		logger.Err(err).Msgf("[getUserHandler] Unable to remove user")
		c.JSON(http.StatusInternalServerError, "Unable to remove user "+username)
		return
	}
	// response
	logger.Info().Msgf("[getUserHandler] removed user")
	c.JSON(http.StatusOK, "removed user "+username)
	return
}

func grantRole(c *gin.Context) {
	// request
	username := c.Param("user")
	role := c.Param("role")

	// process

	if err := userLogic.GrantRole(username, role); err != nil {
		logger.Err(err).Msgf("[getUserHandler] Unable to grant role %s to user %s", role, username)
		c.JSON(http.StatusInternalServerError, "Unable to grant role")
		return
	}

	// response
	logger.Info().Msgf("[getUserHandler] Granted role %s to user %s", role, username)
	c.JSON(http.StatusOK, "Granted role")
	return
}

func revokeRole(c *gin.Context) {
	// request
	username := c.Param("user")
	role := c.Param("role")

	// process

	if err := userLogic.RevokeRole(username, role); err != nil {
		logger.Err(err).Msgf("[getUserHandler] Unable to revoke role %s to user %s", role, username)
		c.JSON(http.StatusInternalServerError, "Unable to revoke role")
		return
	}

	// response
	logger.Info().Msgf("[getUserHandler] Revoked role %s to user %s", role, username)
	c.JSON(http.StatusOK, "Revoked role")
	return
}

func getUserRoles(c *gin.Context) {

}

func getAcc(c *gin.Context) {
	// request
	username := c.Param("username")

	// process
	user, err := userLogic.GetUserByName(username)
	if err != nil {
		logger.Err(err).Msgf("[getUserHandler] Unable to retrieve user")
		c.JSON(http.StatusNotFound, "Unable to retrieve user "+username)
		return
	}

	// response
	user.Password = "" // just to be sure, this field wouldn't be marshalled anyway
	c.JSON(http.StatusOK, user)
	return
}
