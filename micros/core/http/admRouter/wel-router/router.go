package welRouter

import (
	log "bridge/service-managers/logger"
	"fmt"
	"net/http"
	"strconv"

	welLogic "bridge/micros/core/blogic/wel"
	"bridge/micros/core/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/m/wel", mw... /*,middlewares.Author*/)
	gr.POST("/add", addWelAccount)
	gr.POST("/set-status/:acc/:status", setStatus)
	gr.POST("/set/authenticator-prikey", setKey)
	gr.POST("/unset/authenticator-prikey", unsetKey)
	gr.POST("/remove/:acc", removeWelAccount)
	gr.POST("/grant/:role/to/:acc", grantRole)
	gr.POST("/revoke/:role/from/:acc", revokeRole)
	gr.GET("/roles/of/:acc", getAccRoles)
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
	accs, err := welLogic.GetWelAccountsWithRole(role, uint((page-1)*limit), uint(limit))
	if err != nil {
		logger.Err(err).Msgf("[get accs with role handler] Unable to get welups accounts")
		status := http.StatusInternalServerError
		if err == model.ErrWelAccountNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get welups accounts with role "+role)
		return
	}

	// response

	logger.Info().Msgf("[get accs with role handler] Get welups accs with role %s successfully", role)
	c.JSON(http.StatusOK, &accs)
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
	accs, err := welLogic.GetAllWelAccounts(uint((page-1)*limit), uint(limit))
	if err != nil {
		logger.Err(err).Msgf("[get accs handler] Unable to get accs")
		status := http.StatusInternalServerError
		if err == model.ErrWelAccountNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get accs")
		return
	}

	// response

	logger.Info().Msgf("[get accs handler] Get accs %s successfully")
	c.JSON(http.StatusOK, &accs)
	return
}

func getRoles(c *gin.Context) {
	// request

	// process
	roles, err := welLogic.GetAllRoles()
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrWelRoleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get roles")
		return
	}

	// response
	logger.Info().Msgf("[get roles handler] get roles successfully")
	c.JSON(http.StatusOK, &roles)
	return

}

func addWelAccount(c *gin.Context) {
	type addWelAccReq struct {
		Address string
		Status  string
	}

	var req addWelAccReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[add wel account handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := welLogic.AddWelAccount(req.Address, req.Status); err != nil {
		logger.Err(err).Msgf("[add wel account handler] Unable to add wel account")
		c.JSON(http.StatusInternalServerError, "Unable to add wel account")
		return
	}

	logger.Info().Msgf("[add wel account handler] Welups account successfully added")
	c.JSON(http.StatusOK, "Welups account successfully added")
	return
}

func setStatus(c *gin.Context) {
	// request
	acc := c.Param("acc")
	status := c.Param("status")
	if acc == "" || status == "" {
		logger.Err(fmt.Errorf("URI parameters unavailable")).Msgf("[update handler] Invalid request parameters")
		c.JSON(http.StatusBadRequest, "Invalid request parameters")
		return
	}
	// process
	if err := welLogic.SetWelAccountStatus(acc, status); err != nil {
		logger.Err(err).Msgf("[set status handler] Unable to update account %s's status to %s", acc, status)
		c.JSON(http.StatusInternalServerError, "Unable to set status")
		return
	}

	// response
	logger.Info().Msgf("[set status handler] Status updated successfully")
	c.JSON(http.StatusOK, "Status updated successfully")
	return

}

func setKey(c *gin.Context) {
	// request
	type setkeyRequest struct {
		AuthenticatorKey string `json:"authenticator_key"`
	}
	var req setkeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[setkey handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}
	// process
	if err := welLogic.SetCurrentAuthenticator(req.AuthenticatorKey); err != nil {
		logger.Err(err).Msgf("[setkey handler] Unable to set authenticator key")
		c.JSON(http.StatusInternalServerError, "Unable to set authenticator key")
		return
	}

	// response
	logger.Info().Msgf("[setkey handler] Authenticator key set successfully")
	c.JSON(http.StatusOK, "Authenticator key set successfully")
	return

}

func unsetKey(c *gin.Context) {
	// request
	// process
	if err := welLogic.UnsetCurrentAuthenticator(); err != nil {
		logger.Err(err).Msgf("[unsetkey handler] Unable to unset authenticator key")
		c.JSON(http.StatusInternalServerError, "Unable to unset authenticator key")
		return
	}

	// response
	logger.Info().Msgf("[unsetkey handler] Authenticator key unset successfully")
	c.JSON(http.StatusOK, "Authenticator key unset successfully")
	return

}

func removeWelAccount(c *gin.Context) {
	// request
	acc := c.Param("acc")

	// process
	if err := welLogic.RemoveWelAccount(acc); err != nil {
		logger.Err(err).Msgf("[remove wel account handler] Unable to remove welups account")
		c.JSON(http.StatusInternalServerError, "Unable to remove welups account "+acc)
		return
	}
	// response
	logger.Info().Msgf("[remove wel account handler] removed welups account")
	c.JSON(http.StatusOK, "removed welups account "+acc)
	return
}

func grantRole(c *gin.Context) {
	// request
	acc := c.Param("acc")
	role := c.Param("role")

	type adminKey struct {
		AdminKey string `json:"admin_key"`
	}

	var req adminKey
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[grantRole handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// process

	if txid, err := welLogic.GrantRole(acc, role, req.AdminKey); err != nil {
		logger.Err(err).Msgf("[grantRole handler] Unable to grant role %s to account %s with txid %s", role, acc, txid)
		c.JSON(http.StatusInternalServerError, "Unable to grant role")
		return
	}

	// response
	logger.Info().Msgf("[grantRole handler] Granted role %s to account %s", role, acc)
	c.JSON(http.StatusOK, "Granted role")
	return
}

func revokeRole(c *gin.Context) {
	// request
	acc := c.Param("acc")
	role := c.Param("role")

	type adminKey struct {
		AdminKey string `json:"admin_key"`
	}

	var req adminKey
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[revokeRole handler] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
		return
	}

	// process

	if txid, err := welLogic.RevokeRole(acc, role, req.AdminKey); err != nil {
		logger.Err(err).Msgf("[revokeRole handler] Unable to revoke role %s to account %s with txid %s", role, acc, txid)
		c.JSON(http.StatusInternalServerError, "Unable to revoke role")
		return
	}

	// response
	logger.Info().Msgf("[revokeRole handler] Revokeed role %s to account %s", role, acc)
	c.JSON(http.StatusOK, "Revoked role")
	return
}

func getAccRoles(c *gin.Context) {
	// request
	acc := c.Param("acc")
	// process
	roles, err := welLogic.GetWelAccountRoles(acc)
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrWelRoleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get welups account "+acc+"'s roles")
		return
	}

	// response

	c.JSON(http.StatusOK, roles)
	return
}

func getAcc(c *gin.Context) {
	// request
	acc := c.Param("acc")

	// process
	account, err := welLogic.GetWelAccount(acc)
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrWelAccountNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, "Unable to get welups account "+acc)
		return
	}

	// response

	c.JSON(http.StatusOK, account)
	return
}
