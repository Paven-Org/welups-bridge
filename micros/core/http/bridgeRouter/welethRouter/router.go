package welethRouter

import (
	log "bridge/service-managers/logger"
	"net/http"

	ethLogic "bridge/micros/core/blogic/eth"
	welLogic "bridge/micros/core/blogic/wel"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/weleth", mw... /*,middlewares.Author*/)

	gr.POST("/claim/wel/cashin-to/eth", wel2ethCashin)
	gr.POST("/claim/eth/cashout-to/wel", eth2welCashout)

	gr.POST("/claim/eth/cashin-to/wel", eth2welCashin)
	gr.POST("/claim/wel/cashout-to/eth", wel2ethCashout)

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
	// request
	type request struct {
		TxHash           string `json:"txhash"`
		ToTokenAddr      string `json:"to_token_address"`
		ToAccountAddress string `json:"to_account_address"`
		Amount           string `json:"amount"`
	}
	var req request
	contractVersion := "IMPORTS_ETH_v1"
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Claim W2E cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	reqID, signature, err := ethLogic.ClaimWel2EthCashin(req.TxHash, req.ToTokenAddr, req.ToAccountAddress, req.Amount, contractVersion)
	if err != nil {
		logger.Err(err).Msgf("[Claim W2E cashin] Unable to generate request ID and signature")
		c.JSON(http.StatusInternalServerError, "Unable to generate request ID and signature")
	}

	// response
	type response struct {
		ReqID     []byte `json:"request_id"`
		Signature []byte `json:"signature"`
	}
	resp := response{
		ReqID:     reqID,
		Signature: signature,
	}

	logger.Info().Msg("[Claim W2E cashin] successfully generated claim request")
	c.JSON(http.StatusOK, resp)
}

func eth2welCashout(c *gin.Context) {
	// request
	type request struct {
		TxHash           string `json:"txhash"`
		ToTokenAddr      string `json:"to_token_address"`
		ToAccountAddress string `json:"to_account_address"`
		Amount           string `json:"amount"`
	}
	var req request
	contractVersion := "EXPORT_WELUPS_v1"
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Claim E2W cashout] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	reqID, signature, err := welLogic.ClaimEth2WelCashout(req.TxHash, req.ToTokenAddr, req.ToAccountAddress, req.Amount, contractVersion)
	if err != nil {
		logger.Err(err).Msgf("[Claim E2W cashout] Unable to generate request ID and signature")
		c.JSON(http.StatusInternalServerError, "Unable to generate request ID and signature")
	}

	// response
	type response struct {
		ReqID     []byte `json:"request_id"`
		Signature []byte `json:"signature"`
	}
	resp := response{
		ReqID:     reqID,
		Signature: signature,
	}

	logger.Info().Msg("[Claim W2E cashin] successfully generated claim request")
	c.JSON(http.StatusOK, resp)
}

func eth2welCashin(c *gin.Context) {

}

func wel2ethCashout(c *gin.Context) {

}
