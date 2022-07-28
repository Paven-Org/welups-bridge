package service

import (
	"bridge/libs"
	coreEthService "bridge/micros/core/service/eth"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	ethListener "bridge/service-managers/listener/eth"
	"bridge/service-managers/logger"
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"time"

	GotronCommon "github.com/Paven-Org/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
)

type EthConsumer struct {
	ImportContractAddr  string
	MulsendContractAddr string

	WelCashinEthTransDAO  dao.IWelCashinEthTransDAO
	EthCashoutWelTransDAO dao.IEthCashoutWelTransDAO
	WelCashoutEthTransDAO dao.IWelCashoutEthTransDAO

	importAbi  abi.ABI
	mulsendAbi abi.ABI

	tempCli client.Client
}

func NewEthConsumer(iaddr, msaddr string, tempCli client.Client, daos *dao.DAOs) *EthConsumer {
	importAbiJSON, err := os.Open("abi/eth/Import.json")
	if err != nil {
		panic(err)
	}

	defer importAbiJSON.Close()

	importAbi, err := abi.JSON(importAbiJSON)
	if err != nil {
		panic(err)
	}

	mulsendAbiJSON, err := os.Open("abi/eth/MultiSender.json")
	if err != nil {
		panic(err)
	}

	defer mulsendAbiJSON.Close()

	mulsendAbi, err := abi.JSON(mulsendAbiJSON)
	if err != nil {
		panic(err)
	}
	return &EthConsumer{
		ImportContractAddr:  iaddr,
		MulsendContractAddr: msaddr,

		WelCashinEthTransDAO:  daos.WelCashinEthTransDAO,
		EthCashoutWelTransDAO: daos.EthCashoutWelTransDAO,
		WelCashoutEthTransDAO: daos.WelCashoutEthTransDAO,

		importAbi:  importAbi,
		mulsendAbi: mulsendAbi,

		tempCli: tempCli,
	}
}

func (e *EthConsumer) GetConsumer() ([]*ethListener.EventConsumer, error) {
	logger.Get().Info().Msgf("Mulsend's disperse address %s", e.MulsendContractAddr)
	logger.Get().Info().Msgf("Mulsend's disperse signature %s", crypto.Keccak256Hash([]byte(e.mulsendAbi.Events["Disperse"].Sig)))
	return []*ethListener.EventConsumer{
		{
			Address: common.HexToAddress(e.ImportContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.importAbi.Events["Imported"].Sig),
			),
			ParseEvent: e.DoneClaimParser,
		},
		{
			Address: common.HexToAddress(e.ImportContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.importAbi.Events["Withdraw"].Sig),
			),
			ParseEvent: e.DoneDepositParser,
		},
		{
			Address: common.HexToAddress(e.MulsendContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.mulsendAbi.Events["Decline"].Sig),
			),
			ParseEvent: e.DeclineParser,
		},
		{
			Address: common.HexToAddress(e.MulsendContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.mulsendAbi.Events["Disperse"].Sig),
			),
			ParseEvent: e.DisperseParser,
		},
	}, nil
}

func (e *EthConsumer) GetFilterQuery() []ethereum.FilterQuery {
	return []ethereum.FilterQuery{
		ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(e.ImportContractAddr)},
			Topics: [][]common.Hash{{
				crypto.Keccak256Hash(
					[]byte(e.importAbi.Events["Withdraw"].Sig),
				),
				crypto.Keccak256Hash(
					[]byte(e.importAbi.Events["Imported"].Sig),
				),
			}},
		},
		ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(e.MulsendContractAddr)},
			Topics: [][]common.Hash{{
				crypto.Keccak256Hash(
					[]byte(e.mulsendAbi.Events["Decline"].Sig),
				),
				crypto.Keccak256Hash(
					[]byte(e.mulsendAbi.Events["Disperse"].Sig),
				),
			}},
		},
	}

}

func (e *EthConsumer) DisperseParser(l types.Log) error {
	logger.Get().Info().Msgf("[DisperseEV] Disperse event caught at block %d", l.BlockNumber)
	data := make(map[string]interface{})
	e.mulsendAbi.UnpackIntoMap(
		data,
		"Disperse",
		l.Data,
	)
	logger.Get().Info().Msgf("data: %+v", data)
	ethTx := l.TxHash.Hex()
	logger.Get().Info().Msgf("eth tx: %s", ethTx)

	ethToken := common.HexToAddress(l.Topics[1].Hex()).Hex()
	logger.Get().Info().Msgf("eth token: %s", ethToken)

	_receivers := data["receivers"].([]common.Address)
	receivers := libs.Map(
		func(rawAddr common.Address) string {
			hexAddr := rawAddr.Hex()
			return hexAddr
		},
		_receivers)
	logger.Get().Info().Msgf("Receivers: %+v", receivers)

	_remains := data["remains"].([]*big.Int)
	remains := libs.Map(
		func(remain *big.Int) string {
			return remain.String()
		},
		_remains)
	logger.Get().Info().Msgf("Remains: %+v", remains)

	remainOfReceiver := make(map[string]string)
	for i, receiver := range receivers {
		rRemain := remains[i]
		remainOfReceiver[receiver] = rRemain
	}

	// inspect and update corresponding records in DB
	trans, err := e.WelCashoutEthTransDAO.SelectTransByDisperseTxHash(ethTx)
	if err != nil {
		logger.Get().Err(err).Msgf("[DisperseEV] error while retrieving transactions with disperse txhash %s", ethTx)
		return err
	}
	logger.Get().Info().Msgf("trans: %+v", trans)
	for _, tran := range trans {
		if ethToken != tran.EthTokenAddr {
			tran.EthTokenAddr = ethToken
		}

		remain := &big.Int{}
		remain.SetString(remainOfReceiver[tran.EthWalletAddr], 10)
		if remain.Cmp(big.NewInt(0)) != 0 {
			logger.Get().Err(err).Msgf("[DisperseEV] transaction possibly declined: %+v", tran)
			continue
		}
		//tran.Remain = remain.String()

		//total := big.NewInt(0)
		//total.SetString(tran.Total, 10)

		//fee := &big.Int{}
		//fee = fee.Sub(total, remain)
		//tran.CommissionFee = fee.String()

		tran.DisperseStatus = model.WelCashoutEthConfirmed
		tran.DispersedAt = sql.NullTime{Time: time.Now(), Valid: true}

		logger.Get().Info().Msgf("[DisperseEV] tran to be updated: %+v", tran)
		if err := e.WelCashoutEthTransDAO.UpdateWelCashoutEthTx(tran); err != nil {
			logger.Get().Err(err).Msgf("[DisperseEV] error while updating transaction %+v", tran)
		}

	}
	return nil
}

func (e *EthConsumer) DeclineParser(l types.Log) error {
	logger.Get().Info().Msgf("[DeclineEV] Decline event caught at block %d", l.BlockNumber)
	data := make(map[string]interface{})
	e.mulsendAbi.UnpackIntoMap(
		data,
		"Decline",
		l.Data,
	)
	ethTx := l.TxHash.Hex()
	logger.Get().Info().Msgf("eth tx: %s", ethTx)

	ethWalletAddr := common.HexToAddress(l.Topics[1].Hex()).Hex()
	logger.Get().Info().Msgf("ethWalletAddr: %s", ethWalletAddr)

	amount := data["amount"].(*big.Int).String()
	logger.Get().Info().Msgf("amount: %s", amount)

	// set corresponding record in DB: disperse_status = retry
	trans, err := e.WelCashoutEthTransDAO.SelectTransByDisperseTxHashEthAddrAmount(ethTx, ethWalletAddr, amount)
	if err != nil {
		logger.Get().Err(err).Msgf("[DeclineEV] error while retrieving transactions with disperse txhash %s", ethTx)
		return err
	}
	logger.Get().Info().Msgf("trans: %+v", trans)
	for _, tran := range trans {
		if tran.DisperseStatus == model.WelCashoutEthConfirmed {
			continue
		}
		tran.DisperseStatus = model.WelCashoutEthRetry

		logger.Get().Info().Msgf("[DeclineEV] tran to be updated: %+v", tran)
		if err := e.WelCashoutEthTransDAO.UpdateWelCashoutEthTx(tran); err != nil {
			logger.Get().Err(err).Msgf("[DeclineEV] error while updating transaction %+v", tran)
			return err
		}
		// signal WF to re-batch tx
		err = e.tempCli.SignalWorkflow(context.Background(), coreEthService.BatchDisperseID, "", coreEthService.BatchDisperseSignal, tran)
		if err != nil {
			logger.Get().Err(err).Msgf("[DeclineEV] Error sending BatchDisperseWF tx %+v", tran)
			return err
		}

	}
	return nil
}

func (e *EthConsumer) DoneDepositParser(l types.Log) error {
	logger.Get().Info().Msgf("[DepositEV] Deposit event caught at block %d", l.BlockNumber)
	data := make(map[string]interface{})
	e.importAbi.UnpackIntoMap(
		data,
		"Withdraw",
		l.Data,
	)
	amount := data["amount"].(*big.Int).String()
	txHash := l.TxHash.Hex()
	ethWalletAddr := common.HexToAddress(l.Topics[2].Hex()).Hex()

	tran, err := e.EthCashoutWelTransDAO.SelectTransByDepositTxHash(l.TxHash.Hex())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		welWalletAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(l.Topics[3].Bytes()[12:]))
		event := model.EthWelEvent{
			EthTokenAddr: common.HexToAddress(l.Topics[1].Hex()).Hex(),
			//WelWalletAddr: GotronCommon.EncodeCheck(l.Topics[3].Bytes()),
			WelWalletAddr: welWalletAddr,
			NetworkID:     data["networkId"].(*big.Int).String(),
			DepositAt:     time.Now(),
		}
		//m, err := json.Marshal(event)
		//if err != nil {
		//	return fmt.Errorf("can't gen id")
		//}
		//event.ID = crypto.Keccak256Hash(m).Big().String()

		event.DepositTxHash = txHash
		event.EthWalletAddr = ethWalletAddr
		event.WelTokenAddr = model.WelTokenFromEth[event.EthTokenAddr]
		event.Amount = amount
		event.DepositStatus = model.StatusSuccess

		err = e.EthCashoutWelTransDAO.CreateEthCashoutWelTrans(&event)
		if err != nil {
			return err
		}
	} else {
		if tran.DepositStatus != model.StatusSuccess {
			err := e.EthCashoutWelTransDAO.UpdateDepositEthCashoutWelConfirmed(txHash, ethWalletAddr, amount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *EthConsumer) DoneClaimParser(l types.Log) error {
	logger.Get().Info().Msgf("[ClaimEV] Claim event caught at block %d", l.BlockNumber)
	data := make(map[string]interface{})
	e.importAbi.UnpackIntoMap(
		data,
		"Imported",
		l.Data,
	)

	amount := data["amount"].(*big.Int).String()
	ethWalletAddr := common.HexToAddress(l.Topics[3].Hex()).Hex()
	var reqID = &big.Int{}
	reqID.SetBytes(l.Topics[1].Bytes())
	rqId := reqID.String()

	tran, err := e.WelCashinEthTransDAO.SelectTransByRqId(rqId)
	if err != nil {
		return err
	}
	if amount != tran.Amount {
		return fmt.Errorf("Claim wrong amount")
	}
	if ethWalletAddr != tran.EthWalletAddr {
		return fmt.Errorf("Wrong claim eth wallet address")
	}
	_, err = e.WelCashinEthTransDAO.GetClaimRequest(rqId)
	if err == sql.ErrNoRows {
		err := e.WelCashinEthTransDAO.CreateClaimRequest(rqId, tran.ID, model.StatusPending)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if tran.ClaimStatus != model.StatusSuccess {
		err := e.WelCashinEthTransDAO.UpdateClaimWelCashinEth(tran.ID, reqID.String(), model.RequestSuccess, l.TxHash.Hex(), model.StatusSuccess)
		if err != nil {
			return err
		}
	}

	return nil
}

//---------------------------------------------------------------//
type TreasuryMonitor struct {
	treasury_address string
	EthCashinWelDAO  dao.IEthCashinWelTransDAO
}

func MkTreasuryMonitor(address string, daos *dao.DAOs) ethListener.ITxMonitor {
	return &TreasuryMonitor{
		treasury_address: address,
		EthCashinWelDAO:  daos.EthCashinWelTransDAO,
	}
}

func (tm *TreasuryMonitor) MonitoredAddress() common.Address {
	return common.HexToAddress(tm.treasury_address)
}

func (tm *TreasuryMonitor) TxParse(t *types.Transaction, from, to, tokenAddr, amount string) error {
	logger.Get().Info().Msgf("transaction to treasury: %x", t.Hash())
	tx2treasury := &model.TxToTreasury{}

	tx2treasury.TxID = t.Hash().Hex()
	tx2treasury.FromAddress = from
	tx2treasury.TreasuryAddr = to
	tx2treasury.TokenAddr = tokenAddr
	tx2treasury.Amount = amount

	tx_fee := t.Cost()
	tx_fee = tx_fee.Sub(tx_fee, t.Value())
	tx2treasury.TxFee = tx_fee.String()
	tx2treasury.Status = "unconfirmed"
	tx2treasury.CreatedAt = time.Now()
	logger.Get().Info().Msgf("record tx to treasury: %+v\n", tx2treasury)
	if err := tm.EthCashinWelDAO.CreateTx2Treasury(tx2treasury); err != nil {
		logger.Get().Err(err).Msg("Unable to record transaction to treasury")
		return err
	}
	logger.Get().Info().Msg("Recorded transaction to treasury")
	return nil
}
