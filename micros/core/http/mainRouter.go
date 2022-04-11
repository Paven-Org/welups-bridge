package router

import (
	"bridge/common"
	"bridge/service-managers/logger"
	"net/http"

	"bridge/micros/core/http/admRouter"
	"bridge/micros/core/http/bridgeRouter"
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

	// because looks like some platform still doesn't support Allow-Headers wildcard
	allowHeaders := []string{
		"Accept",
		"Accept-CH",
		"Accept-Charset",
		"Accept-Datetime",
		"Accept-Encoding",
		"Accept-Ext",
		"Accept-Features",
		"Accept-Language",
		"Accept-Params",
		"Accept-Ranges",
		"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Age",
		"Allow",
		"Alternates",
		"Authentication-Info",
		"Authorization",
		"C-Ext",
		"C-Man",
		"C-Opt",
		"C-PEP",
		"C-PEP-Info",
		"CONNECT",
		"Cache-Control",
		"Compliance",
		"Connection",
		"Content-Base",
		"Content-Disposition",
		"Content-Encoding",
		"Content-ID",
		"Content-Language",
		"Content-Length",
		"Content-Location",
		"Content-MD5",
		"Content-Range",
		"Content-Script-Type",
		"Content-Security-Policy",
		"Content-Style-Type",
		"Content-Transfer-Encoding",
		"Content-Type",
		"Content-Version",
		"Cookie",
		"Cost",
		"DAV",
		"DELETE",
		"DNT",
		"DPR",
		"Date",
		"Default-Style",
		"Delta-Base",
		"Depth",
		"Derived-From",
		"Destination",
		"Differential-ID",
		"Digest",
		"ETag",
		"Expect",
		"Expires",
		"Ext",
		"From",
		"GET",
		"GetProfile",
		"HEAD",
		"HTTP-date",
		"Host",
		"IM",
		"If",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
		"Keep-Alive",
		"Label",
		"Last-Event-ID",
		"Last-Modified",
		"Link",
		"Location",
		"Lock-Token",
		"MIME-Version",
		"Man",
		"Max-Forwards",
		"Media-Range",
		"Message-ID",
		"Meter",
		"Negotiate",
		"Non-Compliance",
		"OPTION",
		"OPTIONS",
		"OWS",
		"Opt",
		"Optional",
		"Ordering-Type",
		"Origin",
		"Overwrite",
		"P3P",
		"PEP",
		"PICS-Label",
		"POST",
		"PUT",
		"Pep-Info",
		"Permanent",
		"Position",
		"Pragma",
		"ProfileObject",
		"Protocol",
		"Protocol-Query",
		"Protocol-Request",
		"Proxy-Authenticate",
		"Proxy-Authentication-Info",
		"Proxy-Authorization",
		"Proxy-Features",
		"Proxy-Instruction",
		"Public",
		"RWS",
		"Range",
		"Referer",
		"Refresh",
		"Resolution-Hint",
		"Resolver-Location",
		"Retry-After",
		"Safe",
		"Sec-Websocket-Extensions",
		"Sec-Websocket-Key",
		"Sec-Websocket-Origin",
		"Sec-Websocket-Protocol",
		"Sec-Websocket-Version",
		"Security-Scheme",
		"Server",
		"Set-Cookie",
		"Set-Cookie2",
		"SetProfile",
		"SoapAction",
		"Status",
		"Status-URI",
		"Strict-Transport-Security",
		"SubOK",
		"Subst",
		"Surrogate-Capability",
		"Surrogate-Control",
		"TCN",
		"TE",
		"TRACE",
		"Timeout",
		"Title",
		"Trailer",
		"Transfer-Encoding",
		"UA-Color",
		"UA-Media",
		"UA-Pixels",
		"UA-Resolution",
		"UA-Windowpixels",
		"URI",
		"Upgrade",
		"User-Agent",
		"Variant-Vary",
		"Vary",
		"Version",
		"Via",
		"Viewport-Width",
		"WWW-Authenticate",
		"Want-Digest",
		"Warning",
		"Width",
		"X-Content-Duration",
		"X-Content-Security-Policy",
		"X-Content-Type-Options",
		"X-CustomHeader",
		"X-DNSPrefetch-Control",
		"X-Forwarded-For",
		"X-Forwarded-Port",
		"X-Forwarded-Proto",
		"X-Frame-Options",
		"X-Modified",
		"X-OTHER",
		"X-PING",
		"X-PINGOTHER",
		"X-Powered-By",
		"X-Requested-With",
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cnf.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     allowHeaders,
		AllowWildcard:    true,
		AllowCredentials: true,
	}))
	router.Use(helmet.NoSniff(),
		helmet.DNSPrefetchControl(),
		helmet.FrameGuard(),
		helmet.XSSFilter(),
		helmet.IENoOpen())

	router.Use(gzip.Gzip(gzip.BestCompression))

	// version routers
	v1 := versionRouters.MkVRouter("v1", router)

	//// add subrouters
	// adm routes
	admRouter.Config(v1, authMW)
	// bridge routes
	bridgeRouter.Config(v1)

	// public user routes
	userRouter.Config(v1, authMW)

	return router
}
