package welethRouter

import (
	log "bridge/service-managers/logger"
	"fmt"
	"math/big"
	"net/http"

	ethLogic "bridge/micros/core/blogic/eth"
	welLogic "bridge/micros/core/blogic/wel"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Config(router gin.IRouter, mw ...gin.HandlerFunc) {
	initialize()

	gr := router.Group("/weleth", mw... /*,middlewares.Author*/)

	gr.POST("/claim/wel/cashin-to/eth", wel2ethCashin)
	gr.POST("/claim/eth/cashout-to/wel", eth2welCashout)

	gr.POST("/request/eth/cashin-to/wel", eth2welCashin)
	//gr.POST("/claim/wel/cashout-to/eth", wel2ethCashout)

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
		ToAccountAddress string `json:"to_account_address"`
	}
	var req request
	contractVersion := "IMPORTS_ETH_v1"
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Claim W2E cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	tkAddr, amount, reqIDraw, signature, err := ethLogic.ClaimWel2EthCashin(req.TxHash, req.ToAccountAddress, contractVersion)
	if err != nil {
		logger.Err(err).Msgf("[Claim W2E cashin] Unable to generate request ID and signature")
		c.JSON(http.StatusInternalServerError, "Unable to generate request ID and signature")
	}

	reqIDu256 := &big.Int{}
	reqIDu256.SetBytes(reqIDraw)

	// response
	type response struct {
		TokenAddress string `json:"token_address"`
		Amount       string `json:"amount"`
		ReqID        string `json:"request_id"`
		ReqIDHex     string `json:"request_id_hex"`
		ReqIDRaw     []byte `json:"request_id_raw"`
		Signature    []byte `json:"signature"`
		SignatureHex string `json:"signature_hex"`
	}
	resp := response{
		TokenAddress: tkAddr,
		Amount:       amount,
		ReqID:        reqIDu256.String(),
		ReqIDHex:     "0x" + fmt.Sprintf("%x", reqIDraw),
		ReqIDRaw:     reqIDraw,
		Signature:    signature,
		SignatureHex: "0x" + common.Bytes2Hex(signature),
	}

	logger.Info().Msg("[Claim W2E cashin] successfully generated claim request")
	c.JSON(http.StatusOK, resp)
}

func eth2welCashout(c *gin.Context) {
	// request
	type request struct {
		TxHash           string `json:"txhash"`
		ToAccountAddress string `json:"to_account_address"`
	}
	var req request
	contractVersion := "EXPORT_WELUPS_v1"
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Claim E2W cashout] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	tkAddr, amount, reqIDraw, signature, err := welLogic.ClaimEth2WelCashout(req.TxHash, req.ToAccountAddress, contractVersion)
	if err != nil {
		logger.Err(err).Msgf("[Claim E2W cashout] Unable to generate request ID and signature")
		c.JSON(http.StatusInternalServerError, "Unable to generate request ID and signature")
	}

	reqIDu256 := &big.Int{}
	reqIDu256.SetBytes(reqIDraw)

	// response
	type response struct {
		TokenAddress string `json:"token_address"`
		Amount       string `json:"amount"`
		ReqID        string `json:"request_id"`
		ReqIDHex     string `json:"request_id_hex"`
		ReqIDRaw     []byte `json:"request_id_raw"`
		Signature    []byte `json:"signature"`
		SignatureHex string `json:"signature_hex"`
	}
	resp := response{
		TokenAddress: tkAddr,
		Amount:       amount,
		ReqID:        reqIDu256.String(),
		ReqIDHex:     "0x" + fmt.Sprintf("%x", reqIDraw),
		ReqIDRaw:     reqIDraw,
		Signature:    signature,
		SignatureHex: "0x" + common.Bytes2Hex(signature),
	}

	logger.Info().Msg("[Claim E2W cashout] successfully generated claim request")
	c.JSON(http.StatusOK, resp)
}

func eth2welCashin(c *gin.Context) {
	//request
	type request struct {
		From     string `json:"from_eth"`
		To       string `json:"to_wel"`
		Treasury string `json:"eth_treasury"`
		NetId    string `json:"netid"`
		Token    string `json:"token"`
		Amount   string `json:"amount"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[E2W cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	if err := ethLogic.
		WatchTx2TreasuryRequest(
			req.From,
			req.To,
			req.Treasury,
			req.NetId,
			req.Token,
			req.Amount); err != nil {
		logger.Err(err).Msgf("[E2W cashin] Failed to request backend to watch for transaction to treasury")
		c.JSON(http.StatusBadRequest, "Failed to request backend to watch for transaction to treasury")
	}

	// response
	c.JSON(http.StatusOK, fmt.Sprintf("BE is watching for transaction to %s from %s", req.From, req.Treasury))
}

func wel2ethCashout(c *gin.Context) {

}
