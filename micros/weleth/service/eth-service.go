package service

import (
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	ethListener "bridge/service-managers/listener/eth"
	"database/sql"
	"encoding/json"
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
	ContractAddr   string
	WelEthTransDAO dao.IWelEthTransDAO
	abi            abi.ABI
}

func NewEthConsumer(addr string, welEthTransDAO dao.IWelEthTransDAO) *EthConsumer {
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
		ContractAddr:   addr,
		WelEthTransDAO: welEthTransDAO,
		abi:            abi,
	}
}

func (e *EthConsumer) GetConsumer() ([]*ethListener.EventConsumer, error) {
	return []*ethListener.EventConsumer{
		{
			Address: common.HexToAddress(e.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Imported"].Sig),
			),
			ParseEvent: e.DoneDepositParser,
		},
		{
			Address: common.HexToAddress(e.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Withdraw"].Sig),
			),

			ParseEvent: e.DoneClaimParser,
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

	tran, err := e.WelEthTransDAO.SelectTransById(l.TxHash.Hex())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		event := model.WelEthEvent{
			EthTokenAddr:  common.HexToAddress(l.Topics[1].Hex()).Hex(),
			WelWalletAddr: GotronCommon.EncodeCheck(l.Topics[3].Bytes()),
			NetworkID:     data["networkId"].(*big.Int).String(),
			DepositAt:     time.Now(),
		}
		m, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("can't gen id")
		}
		event.ID = crypto.Keccak256Hash(m).Big().String()

		event.DepositTxHash = txHash
		event.EthWalletAddr = ethWalletAddr
		event.Amount = amount
		event.DepositStatus = model.StatusSuccess

		err = e.WelEthTransDAO.CreateEthWelTrans(&event)
		if err != nil {
			return err
		}
	} else {
		if tran.DepositStatus != model.StatusSuccess {
			err := e.WelEthTransDAO.UpdateDepositEthWelConfirmed(txHash, ethWalletAddr, amount)
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

	tran, err := e.WelEthTransDAO.SelectTransById(reqID.String())
	if err != nil {
		return err
	}
	if amount != tran.Amount {
		return fmt.Errorf("Claim wrong amount")
	}
	if ethWalletAddr != tran.EthWalletAddr {
		return fmt.Errorf("Wrong claim eth wallet address")
	}

	if tran.ClaimStatus != model.StatusSuccess {
		err := e.WelEthTransDAO.UpdateClaimWelEth(reqID.String(), l.TxHash.Hex(), model.StatusSuccess)
		if err != nil {
			return err
		}
	}

	return nil
}
