package manageUserRouter

import (
	log "bridge/service-managers/logger"
	"net/http"
	"strconv"

	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/m/u", mw... /*,middlewares.Author*/)

	gr.POST("/add", addUser)
	gr.POST("/update/:user", updateUser)
	gr.POST("/remove/:user", removeUser)
	gr.POST("/grant/:role/to/:user", grantRole)
	gr.POST("/revoke/:role/from/:user", revokeRole)
	gr.GET("/haverole/:role/:page", getUsersWithRole)
	gr.GET("/users/:page", getUsers)
	gr.POST("/ban/:user", banUser)
	gr.GET("/roles", getRoles)
	gr.GET("/roles/of/:user", getRolesOfUser)

}

func initialize() {
	logger = log.Get()
	logger.Info().Msg("manage users handlers initialized")
}

func getUsersWithRole(c *gin.Context) {
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
		limit = 15 // default
	} else {
		limit, err = strconv.ParseUint(_limit, 10, 32)
		if err != nil {
			logger.Err(err).Msgf("[get users with role handler] invalid limit")
			limit = 15 // default
		}
	}

	// process
	users, err := userLogic.GetUsersWithRole(role, uint((page-1)*limit), uint(limit))
	if err != nil && err != model.ErrUserNotFound {
		logger.Err(err).Msgf("[get users with role handler] Unable to get users")
		c.JSON(http.StatusInternalServerError, "Unable to get users with role "+role)
		return
	}

	// response

	logger.Info().Msgf("[get users with role handler] Get users with role %s successfully", role)
	c.JSON(http.StatusOK, &users)
	return
}

func getUsers(c *gin.Context) {
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
		limit = 15 // default
	} else {
		limit, err = strconv.ParseUint(_limit, 10, 32)
		if err != nil {
			logger.Err(err).Msgf("[get users handler] invalid limit")
			limit = 15 // default
		}
	}

	// process
	users, err := userLogic.GetUsers(uint((page-1)*limit), uint(limit))
	if err != nil && err != model.ErrUserNotFound {
		logger.Err(err).Msgf("[get users handler] Unable to get users")
		c.JSON(http.StatusInternalServerError, "Unable to get users")
		return
	}

	// response

	logger.Info().Msgf("[get users handler] Get users successfully")
	c.JSON(http.StatusOK, &users)
	return
}

func getRoles(c *gin.Context) {
	// request

	// process
	roles, err := userLogic.GetAllRoles()
	if err != nil && err != model.ErrRoleNotFound {
		logger.Err(err).Msgf("[get roles handler] Unable to get roles")
		c.JSON(http.StatusInternalServerError, "Unable to get roles")
		return
	}

	// response
	logger.Info().Msgf("[get roles handler] get roles successfully")
	c.JSON(http.StatusOK, &roles)
	return

}

func addUser(c *gin.Context) {
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

func updateUser(c *gin.Context) {
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

func removeUser(c *gin.Context) {
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

func banUser(c *gin.Context) {
	// request
	username := c.Param("user")

	// process
	if err := userLogic.AdminUpdateUserInfo(username, "", "", "", model.UserStatusBanned); err != nil {
		logger.Err(err).Msgf("[ban handler] Unable to ban user")
		c.JSON(http.StatusInternalServerError, "Unable to ban user")
		return
	}

	// response
	logger.Info().Msgf("[ban handler] User banned successfully")
	c.JSON(http.StatusOK, "User banned successfully")
	return

}

func getRolesOfUser(c *gin.Context) {
	// request
	username := c.Param("user")

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
