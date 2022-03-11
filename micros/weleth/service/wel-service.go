package service

import (
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	GotronCommon "github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type WelConsumer struct {
	ContractAddr   string
	WelEthTransDAO dao.IWelEthTransDAO
	abi            abi.ABI
}

func NewWelConsumer(addr string, welEthTransDAO dao.IWelEthTransDAO) *WelConsumer {
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
		ContractAddr:   addr,
		WelEthTransDAO: welEthTransDAO,
		abi:            abi,
	}
}

func (e *WelConsumer) GetConsumer() ([]*welListener.EventConsumer, error) {
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
	}, nil
}

func (e *WelConsumer) DoneReturnParser(t *welListener.Transaction) error {
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Returned",
		t.Log[0].Data,
	)

	rqId := data["requestId"].(*big.Int).String()
	//tokenAddr := data["token"].(common.Address)
	welWalletAddr := GotronCommon.EncodeCheck(data["user"].(common.Address).Bytes())
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	if t.Result == "unconfirmed" {
		tran, _ := e.WelEthTransDAO.SelectTransById(t.Hash)
		if tran == nil {
			return fmt.Errorf("can't find this transaction")
		} else {
			err := e.WelEthTransDAO.UpdateClaimEthWel(rqId, t.Hash, welWalletAddr, amount, fee, model.StatusUnknown)
			if err != nil {
				return err
			}
		}

	} else if t.Result == "confirmed" {
		err := e.WelEthTransDAO.UpdateClaimEthWel(rqId, t.Hash, welWalletAddr, amount, fee, model.StatusSuccess)
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
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Withdraw",
		t.Log[0].Data,
	)
	ethTokenAddr := data["to"].([]byte)
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	var networkID = &big.Int{}
	if t.Result == "unconfirmed" {
		tran, _ := e.WelEthTransDAO.SelectTransByDepositTxHash(t.Hash)

		// somehow, we did not save this deposit to db before
		if tran == nil {
			event := model.WelEthEvent{
				WelTokenAddr: "0x" + GotronCommon.ToHex(t.Log[0].Topics[1]),
				EthTokenAddr: common.BytesToAddress(ethTokenAddr).Hex(),
				NetworkID:    networkID.SetBytes(t.Log[0].Topics[3]).String(),
				DepositAt:    time.Now(),
			}
			m, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("can't gen id")
			}
			event.ID = crypto.Keccak256Hash(m).Big().String()

			event.DepositTxHash = t.Hash
			event.WelWalletAddr = GotronCommon.EncodeCheck(t.Log[0].Topics[2])
			event.DepositAmount = amount
			event.Fee = fee

			_ = e.WelEthTransDAO.CreateWelEthTrans(&event)
		}
	} else if t.Result == "confirmed" {
		err := e.WelEthTransDAO.UpdateDepositWelEthConfirmed(t.Hash, GotronCommon.EncodeCheck(t.Log[0].Topics[2]), amount, fee)
		if err != nil {
			return err
		}
		// emit done deposit event, save to db
	} else {
		return fmt.Errorf("unknown status")
	}

	return nil
}
