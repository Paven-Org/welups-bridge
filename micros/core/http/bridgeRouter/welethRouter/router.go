package welethRouter

import (
	log "bridge/service-managers/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/weleth", mw... /*,middlewares.Author*/)

	gr.POST("/claim/wel/cashin-to/eth/:wel-txid", wel2ethCashin)
	gr.POST("/claim/eth/cashout-to/wel/:eth-txid", eth2welCashout)

	gr.POST("/claim/eth/cashin-to/wel/:eth-txid", eth2welCashin)
	gr.POST("/claim/wel/cashout-to/eth/:wel-txid", wel2ethCashout)

	//	gr.GET("/transaction/cashin/from/eth/:txid")
	//	gr.GET("/transaction/cashout/to/eth/:txid")
	//	gr.GET("/transaction/cashin/from/wel/:txid")
	//	gr.GET("/transaction/cashout/to/wel/:txid")
}

func initialize() {
	logger = log.Get()
	logger.Info().Msg("weleth bridge handlers initialized")
}

func wel2ethCashin(c *gin.Context) {

}

func eth2welCashout(c *gin.Context) {

}

func eth2welCashin(c *gin.Context) {

}

func wel2ethCashout(c *gin.Context) {

}
