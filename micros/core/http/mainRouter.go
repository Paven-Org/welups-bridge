package router

import (
	"bridge/common"
	"bridge/service-managers/logger"
	"net/http"

	"bridge/micros/core/http/admRouter"
	userRouter "bridge/micros/core/http/userRouter"
	"bridge/micros/core/http/versionRouters"
	"bridge/micros/core/middlewares"

	helmet "github.com/danielkov/gin-helmet"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Init main router
// authMW is special and should be constructed separately since it uses certain
// infrastructures the http server doesn't need to be aware of
func InitMainRouter(cnf common.HttpConf, authMW gin.HandlerFunc) *gin.Engine {
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
	admRouter.Config(v1, authMW)

	// public routes
	userRouter.Config(v1, authMW)

	return router
}
