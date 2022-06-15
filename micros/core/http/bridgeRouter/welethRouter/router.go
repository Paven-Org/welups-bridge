package welethRouter

import (
	log "bridge/service-managers/logger"
	"fmt"
	"math/big"
	"net/http"

	ethLogic "bridge/micros/core/blogic/eth"
	welLogic "bridge/micros/core/blogic/wel"
	"bridge/micros/weleth/model"

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
	gr.POST("/request/eth/cashin-to/wel/:txid", eth2welCashinByTxId)
	//gr.POST("/claim/wel/cashout-to/eth", wel2ethCashout)

	gr.GET("/transaction/eth/cashin/wel/:eth_txid", getE2WCashinTxByEthTxId)
	gr.GET("/transactions/eth/cashin/wel", getE2WCashinTx)
	gr.GET("/transactions/wel/cashout/eth", getW2ECashoutTx)
	gr.GET("/transactions/wel/cashin/eth", getW2ECashinTx)
	gr.GET("/transactions/eth/cashout/wel", getE2WCashoutTx)
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
	c.JSON(http.StatusOK, fmt.Sprintf("BE is watching for transaction to %s from %s", req.Treasury, req.From))
}

func eth2welCashinByTxId(c *gin.Context) {
	//request
	type request struct {
		To    string `json:"to_wel"`
		NetId string `json:"netid"`
		Token string `json:"token"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[E2W cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	txhash := c.Param("txid")
	if len(txhash) <= 0 {
		err := fmt.Errorf("Invalid request payload")
		logger.Err(err).Msgf("[E2W cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	if err := ethLogic.
		WatchTx2TreasuryRequestByTxhash(
			txhash,
			req.To,
			req.NetId,
			req.Token); err != nil {
		logger.Err(err).Msgf("[E2W cashin] Failed to request backend to watch for transaction to treasury")
		c.JSON(http.StatusBadRequest, "Failed to request backend to watch for transaction to treasury")
	}

	// response
	c.JSON(http.StatusOK, fmt.Sprintf("BE is confirming transaction to treasury with transaction id %s", txhash))
}

func getE2WCashinTxByEthTxId(c *gin.Context) {
	//request
	txhash := c.Param("eth_txid")
	if len(txhash) <= 0 {
		err := fmt.Errorf("Invalid request payload")
		logger.Err(err).Msgf("[E2W cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	tx, err := ethLogic.GetE2WCashinTransByEthTxHash(txhash)
	if err != nil {
		logger.Err(err).Msgf("[E2W cashin] Failed to get E2W cashin transaction with eth side transaction id " + txhash)
		c.JSON(http.StatusBadRequest, "Failed to get E2W cashin transaction with eth side transaction id "+txhash)
	}

	// response
	c.JSON(http.StatusOK, tx)
}

func getW2ECashinTx(c *gin.Context) {
	// request
	type request struct {
		Sender   string `json:"from_wel"`
		Receiver string `json:"to_eth"`
		Status   string `json:"withdraw_status"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Get W2E cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	txs, err := welLogic.GetW2ECashinTrans(req.Sender, req.Receiver, req.Status)
	if err != nil {
		logger.Err(err).Msgf("[Get W2E cashin] Unable to get W2E cashin transactions")
		c.JSON(http.StatusInternalServerError, "Unable to get W2E cashin transactions")
	}

	// response

	logger.Info().Msg("[Get W2E cashin] successfully get W2E cashin transactions")
	c.JSON(http.StatusOK, txs)
}

func getE2WCashoutTx(c *gin.Context) {
	// request
	type request struct {
		Sender   string `json:"from_eth"`
		Receiver string `json:"to_wel"`
		Status   string `json:"withdraw_status"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Get E2W cashout] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	txs, err := ethLogic.GetE2WCashoutTrans(req.Sender, req.Receiver, req.Status)
	if err != nil {
		logger.Err(err).Msgf("[Get E2W cashout] Unable to get E2W cashout transactions")
		c.JSON(http.StatusInternalServerError, "Unable to get E2W cashout transactions")
	}

	// response

	logger.Info().Msg("[Get E2W cashout] successfully get E2W cashout transactions")
	c.JSON(http.StatusOK, txs)
}

func getE2WCashinTx(c *gin.Context) {
	// request
	type request struct {
		Sender   string `json:"from_eth"`
		Receiver string `json:"to_wel"`
		Status   string `json:"cashin_tx_status"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Get E2W cashin] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	txs, tx2tr, err := ethLogic.GetE2WCashinTrans(req.Sender, req.Receiver, req.Status)
	if err != nil {
		logger.Err(err).Msgf("[Get E2W cashin] Unable to get E2W cashin transactions")
		c.JSON(http.StatusInternalServerError, "Unable to get E2W cashin transactions")
	}

	// response
	type response struct {
		CashinTx     []model.EthCashinWelTrans `json:"cashin_tx"`
		TxToTreasury []model.TxToTreasury      `json:"to_treasury_tx"`
	}
	resp := response{
		CashinTx:     txs,
		TxToTreasury: tx2tr,
	}

	logger.Info().Msg("[Get E2W cashin] successfully get E2W cashin transactions")
	c.JSON(http.StatusOK, resp)
}

func getW2ECashoutTx(c *gin.Context) {
	// request
	type request struct {
		Sender   string `json:"from_wel"`
		Receiver string `json:"to_eth"`
		Status   string `json:"status"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Err(err).Msgf("[Get W2E cashout] Invalid request payload")
		c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// process
	txs, err := welLogic.GetW2ECashoutTrans(req.Sender, req.Receiver, req.Status)
	if err != nil {
		logger.Err(err).Msgf("[Get W2E cashout] Unable to get W2E cashout transactions")
		c.JSON(http.StatusInternalServerError, "Unable to get W2E cashout transactions")
	}

	// response

	logger.Info().Msg("[Get W2E cashout] successfully get W2E cashout transactions")
	c.JSON(http.StatusOK, txs)
}

func wel2ethCashout(c *gin.Context) {

}
