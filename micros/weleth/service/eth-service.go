package service

import (
	"bridge/libs"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	ethListener "bridge/service-managers/listener/eth"
	"bridge/service-managers/logger"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"time"

	GotronCommon "github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthConsumer struct {
	ContractAddr          string
	WelCashinEthTransDAO  dao.IWelCashinEthTransDAO
	EthCashoutWelTransDAO dao.IEthCashoutWelTransDAO
	abi                   abi.ABI
}

func NewEthConsumer(addr string, daos *dao.DAOs) *EthConsumer {
	importAbiJSON, err := os.Open("abi/eth/Import.json")
	if err != nil {
		panic(err)
	}

	defer importAbiJSON.Close()

	abi, err := abi.JSON(importAbiJSON)
	if err != nil {
		panic(err)
	}

	return &EthConsumer{
		ContractAddr:          addr,
		WelCashinEthTransDAO:  daos.WelCashinEthTransDAO,
		EthCashoutWelTransDAO: daos.EthCashoutWelTransDAO,
		abi:                   abi,
	}
}

func (e *EthConsumer) GetConsumer() ([]*ethListener.EventConsumer, error) {
	return []*ethListener.EventConsumer{
		{
			Address: common.HexToAddress(e.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Imported"].Sig),
			),
			ParseEvent: e.DoneClaimParser,
		},
		{
			Address: common.HexToAddress(e.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Withdraw"].Sig),
			),

			ParseEvent: e.DoneDepositParser,
		},
	}, nil
}

func (e *EthConsumer) GetFilterQuery() ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(e.ContractAddr)},
		Topics: [][]common.Hash{{
			crypto.Keccak256Hash(
				[]byte(e.abi.Events["Withdraw"].Sig),
			),
			crypto.Keccak256Hash(
				[]byte(e.abi.Events["Imported"].Sig),
			),
		}},
	}

}

func (e *EthConsumer) DoneDepositParser(l types.Log) error {
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
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
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
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

func MkTreasuryMonitor(address string) ethListener.ITxMonitor {
	return &TreasuryMonitor{
		treasury_address: address,
	}
}

func (tm *TreasuryMonitor) MonitoredAddress() common.Address {
	return common.HexToAddress(tm.treasury_address)
}

func (tm *TreasuryMonitor) TxParse(t *types.Transaction, from, to, tokenAddr, amount string) error {
	fmt.Printf("transaction: %x\n", t.Hash())
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
	fmt.Printf("record tx to treasury: %+v\n", tx2treasury)
	if err := tm.EthCashinWelDAO.CreateTx2Treasury(tx2treasury); err != nil {
		logger.Get().Err(err).Msg("Unable to record transaction to treasury")
		return err
	}
	logger.Get().Info().Msg("Recorded transaction to treasury")
	return nil
}
