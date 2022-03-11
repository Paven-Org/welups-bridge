package service

import (
	"bridge/common/utils"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	GotronCommon "github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type WelConsumer struct {
	ContractAddr  string
	WelDepositDAO dao.IWelTransDAO
	abi           abi.ABI
}

func NewDoneDepositConsumer(addr, vaultAddr string, welDepositDAO dao.IWelTransDAO) *WelConsumer {
	exportAbiJSON, err := os.Open("abi/wel/Export.json")
	if err != nil {
		panic(err)
	}

	defer exportAbiJSON.Close()

	abi, err := abi.JSON(exportAbiJSON)
	if err != nil {
		panic(err)
	}

	return &WelConsumer{
		ContractAddr:  addr,
		WelDepositDAO: welDepositDAO,
		abi:           abi,
	}
}

func (e *WelConsumer) GetConsumer() []*welListener.EventConsumer {
	return []*welListener.EventConsumer{
		{
			Address: e.ContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Withdraw"].Sig),
			),
			ParseEvent: e.DoneDepositParser,
		},
		{
			Address: e.ContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.abi.Events["Returned"].Sig),
			),

			ParseEvent: e.DoneReturnParser,
		},
	}
}

func (e *WelConsumer) DoneReturnParser(t *welListener.Transaction) error {
	// filter if the to wallet is our vault wallet
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Withdraw",
		t.Log[0].Data,
	)
	ethTokenAddr := data["to"].([]byte)
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	if t.Result == "unconfirmed" {
		tran, _ := e.WelDepositDAO.SelectTransByTxHash(t.Hash)

		// somehow, we did not save this deposit to db before
		if tran == nil {
			event := model.DoneDepositEvent{
				TxHash:       t.Hash,
				WelTokenAddr: GotronCommon.ToHex(t.Log[0].Topics[1]),
				FromAddr:     GotronCommon.EncodeCheck(t.Log[0].Topics[2]),
				NetworkID:    binary.BigEndian.Uint64(t.Log[0].Topics[3]),
				EthTokenAddr: common.BytesToAddress(ethTokenAddr).Hex(),
			}
			m, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("can't gen id")
			}
			event.ID = string(utils.HashKeccak256(m))

			_ = e.WelDepositDAO.CreateTrans(&event)
		}
	} else if t.Result == "confirmed" {
		err := e.WelDepositDAO.UpdateVerified(t.Hash, amount, fee)
		if err != nil {
			return err
		}
		// emit done deposit event, save to db
	} else {
		return fmt.Errorf("unknown status")
	}

	return nil
}

func (e *WelConsumer) DoneDepositParser(t *welListener.Transaction) error {
	// filter if the to wallet is our vault wallet
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Withdraw",
		t.Log[0].Data,
	)
	ethTokenAddr := data["to"].([]byte)
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	if t.Result == "unconfirmed" {
		tran, _ := e.WelDepositDAO.SelectTransByTxHash(t.Hash)

		// somehow, we did not save this deposit to db before
		if tran == nil {
			event := model.DoneDepositEvent{
				TxHash:       t.Hash,
				WelTokenAddr: GotronCommon.ToHex(t.Log[0].Topics[1]),
				FromAddr:     GotronCommon.EncodeCheck(t.Log[0].Topics[2]),
				NetworkID:    binary.BigEndian.Uint64(t.Log[0].Topics[3]),
				EthTokenAddr: common.BytesToAddress(ethTokenAddr).Hex(),
			}
			m, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("can't gen id")
			}
			event.ID = string(utils.HashKeccak256(m))

			_ = e.WelDepositDAO.CreateTrans(&event)
		}
	} else if t.Result == "confirmed" {
		err := e.WelDepositDAO.UpdateVerified(t.Hash, amount, fee)
		if err != nil {
			return err
		}
		// emit done deposit event, save to db
	} else {
		return fmt.Errorf("unknown status")
	}

	return nil
}
