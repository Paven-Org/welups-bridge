package service

import (
	"bridge/libs"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"database/sql"
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

func (e *WelConsumer) DoneReturnParser(t *welListener.Transaction, logpos int) error {
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Returned",
		t.Log[logpos].Data,
	)

	rqId := data["requestId"].(*big.Int).String()
	//tokenAddr := data["token"].(common.Address)
	welWalletAddr := GotronCommon.EncodeCheck(data["user"].(common.Address).Bytes())
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	tran, err := e.WelEthTransDAO.SelectTransById(rqId)
	if err != nil {
		return err
	}
	if amount != tran.Amount {
		return fmt.Errorf("Claim wrong amount")
	}
	if tran.WelWalletAddr != welWalletAddr {
		return fmt.Errorf("Wrong claim wel wallet address")
	}

	if t.Result == "unconfirmed" {
		if tran.ClaimStatus != model.StatusUnknown {
			err := e.WelEthTransDAO.UpdateClaimEthWel(rqId, t.Hash, fee, model.StatusUnknown)
			if err != nil {
				return err
			}
		}
	} else if t.Result == "confirmed" {
		if tran.ClaimStatus != model.StatusSuccess {
			err := e.WelEthTransDAO.UpdateClaimEthWel(rqId, t.Hash, fee, model.StatusSuccess)
			if err != nil {
				return err
			}
		}

		// emit done deposit event, save to db
	} else {
		return fmt.Errorf("unknown status")
	}

	return nil
}

func (e *WelConsumer) DoneDepositParser(t *welListener.Transaction, logpos int) error {
	data := make(map[string]interface{})
	e.abi.UnpackIntoMap(
		data,
		"Withdraw",
		t.Log[logpos].Data,
	)
	ethWalletAddr := data["to"].(common.Address)
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	var networkID = &big.Int{}
	if t.Status == "unconfirmed" {
		// NOTE: if front end can't get txHash then we will need to fix this
		_, err := e.WelEthTransDAO.SelectTransByDepositTxHash(t.Hash)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			// somehow, we did not save this deposit to db before
			event := model.WelEthEvent{
				WelTokenAddr:  GotronCommon.BytesToHexString(t.Log[logpos].Topics[1]),
				EthWalletAddr: ethWalletAddr.Hex(),
				NetworkID:     networkID.SetBytes(t.Log[logpos].Topics[3]).String(),
				DepositAt:     time.Now(),
			}
			m, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("can't gen id")
			}
			event.ID = crypto.Keccak256Hash(m).Big().String()

			event.DepositTxHash = t.Hash
			event.WelWalletAddr, _ = libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
			event.Amount = amount
			event.DepositStatus = model.StatusUnknown
			event.Fee = fee

			_ = e.WelEthTransDAO.CreateWelEthTrans(&event)
		}
	} else if t.Status == "confirmed" {
		tran, err := e.WelEthTransDAO.SelectTransByDepositTxHash(t.Hash)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			// somehow, we did not save this deposit to db before
			event := model.WelEthEvent{
				WelTokenAddr:  "0x" + GotronCommon.ToHex(t.Log[logpos].Topics[1]),
				EthWalletAddr: ethWalletAddr.Hex(),
				NetworkID:     networkID.SetBytes(t.Log[logpos].Topics[3]).String(),
				DepositAt:     time.Now(),
			}
			m, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("can't gen id")
			}
			event.ID = crypto.Keccak256Hash(m).Big().String()

			event.DepositTxHash = t.Hash
			event.WelWalletAddr, _ = libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
			event.Amount = amount
			event.DepositStatus = model.StatusSuccess
			event.Fee = fee

			err = e.WelEthTransDAO.CreateWelEthTrans(&event)
			if err != nil {
				logger.Get().Err(err).Msg("[DoneDeposit] can't create new transaction")
				return err
			}
		} else {
			if tran.DepositStatus != model.StatusSuccess {
				err := e.WelEthTransDAO.UpdateDepositWelEthConfirmed(t.Hash, GotronCommon.EncodeCheck(t.Log[logpos].Topics[2]), amount, fee)
				if err != nil {
					return err
				}
			}
		}
		// emit done deposit event, save to db
	} else {
		return fmt.Errorf("unknown status")
	}

	return nil
}
