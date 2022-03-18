package middlewares

import (
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	log "bridge/service-managers/logger"
	"context"
	"fmt"
	"net/http"
	"strings"

	casbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func MkAuthMW(enforcer *casbin.Enforcer, rm *manager.RedisManager) gin.HandlerFunc {
	logger := log.Get()
	redis, err := rm.GetRedisClient(manager.StdAuthDBName)
	if err != nil {
		logger.Err(err).Msgf("[AuthMW] Failed to get redis connection")
		return func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, "Unable to authorize user due to internal service error")
			return
		}
	}
	logger.Info().Msg("[AuthMW] connected to redis server")

	return func(c *gin.Context) {

		// get and process credentials
		tokenS := c.GetHeader("Authorization")
		tokenS = strings.TrimPrefix(tokenS, "Bearer ")
		logger.Debug().Msg("[AuthMW] Token string: " + tokenS)

		claims, err := userLogic.ParseTokenToClaims(tokenS)
		if err != nil {
			logger.Err(err).Msgf("[AuthMW] Unable to parse JWT")
			c.JSON(http.StatusBadRequest, "Unable to parse JWT")
			return
		}
		username := claims.Username
		sessionID := claims.Session

		cookie, err := c.Cookie(sessionID)
		if err != nil {
			logger.Err(err).Msgf("[AuthMW] No cookie in request")
			c.JSON(http.StatusBadRequest, "No cookie in request")
			return
		}

		ctx := context.Background()
		sessionSecret, err := redis.
			Get(ctx,
				fmt.Sprintf("session:user_%s:%s", username, sessionID)).
			Result()

		if err != nil {
			logger.Err(err).Msgf("[AuthMW] Error while authorizing user %s", username)
			c.JSON(http.StatusServiceUnavailable, "Unable to authorize user due to internal service error")
			return
		}

		if cookie != sessionSecret {
			err := model.ErrInconsistentCredentials
			logger.Err(err).Msgf("[AuthMW] Error while authorizing user %s", username)
			c.JSON(http.StatusUnauthorized, "Unable to authorize user")
			return
		}
		// save claims into context
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)
		roles, err := userLogic.GetUserRoles(claims.Username)
		if err != nil {
			logger.Err(err).Msgf("[AuthMW] Error while authorizing user %s", username)
			c.JSON(http.StatusServiceUnavailable, "Unable to authorize user due to internal service error")
			return
		}
		c.Set("roles", roles)

		// enforcing rbac policies
		action := c.Request.Method
		obj := c.FullPath()

		var authorized bool
		for _, role := range roles {
			authorized = false
			var err error
			logger.Debug().Msg("Role: " + role)
			authorized, err = enforcer.Enforce(role, obj, action)
			if err != nil {
				logger.Debug().Msgf("[AuthMW] Failed to authorize role %s, error: %s", role, err.Error())
				continue
			}
			if authorized {
				logger.Info().Msgf("[AuthMW] role %s authorized", role)
				break
			}
		}

		if !authorized {
			logger.Debug().Msgf("[AuthMW] Failed to authorize user %s", claims.Username)
			c.JSON(http.StatusUnauthorized, "Unable to authorize user")
			return
		}
		// next
		c.Next()
	}
}
