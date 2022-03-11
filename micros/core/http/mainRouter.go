package router

import (
	"bridge/common"
	"bridge/service-managers/logger"
	"net/http"

	"bridge/micros/core/http/admRouter"
	"bridge/micros/core/http/publicRouter"
	"bridge/micros/core/http/versionRouters"
	"bridge/micros/core/middlewares"

	helmet "github.com/danielkov/gin-helmet"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitMainRouter(cnf common.HttpConf) *gin.Engine {
	router := gin.New()

	// global middlewares...
	router.Use(nice.Recovery(
		func(c *gin.Context, err interface{}) {
			logger.Get().Err(err.(error)).Msg("[Recovery mw] Bad Gateway, error")
			c.JSON(http.StatusBadGateway, err)
		}))

	router.Use(middlewares.LoggerMw())

	router.Use(cors.New(cors.Config{
		AllowOrigins:  cnf.CORSAllowOrigins,
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.NoSniff(),
		helmet.DNSPrefetchControl(),
		helmet.FrameGuard(),
		helmet.XSSFilter(),
		helmet.IENoOpen())

	router.Use(gzip.Gzip(gzip.BestCompression))

	// version routers
	v1 := versionRouters.MkVRouter("v1", router)

	// add subrouters
	admRouter.Config(v1)

	// public routes
	publicRouter.Config(v1)

	return router
}
