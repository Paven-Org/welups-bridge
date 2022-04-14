package service

import (
	"bridge/libs"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"database/sql"
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
	ContractAddr          string
	WelCashinEthTransDAO  dao.IWelCashinEthTransDAO
	EthCashoutWelTransDAO dao.IEthCashoutWelTransDAO
	abi                   abi.ABI
}

func NewWelConsumer(addr string, daos *dao.DAOs) *WelConsumer {
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
		ContractAddr:          addr,
		WelCashinEthTransDAO:  daos.WelCashinEthTransDAO,
		EthCashoutWelTransDAO: daos.EthCashoutWelTransDAO,
		abi:                   abi,
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
	welWalletAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(data["user"].(common.Address).Bytes())) //FIX!
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()

	tran, err := e.EthCashoutWelTransDAO.SelectTransByRqId(rqId)
	if err != nil {
		return err
	}
	if amount != tran.Amount {
		return fmt.Errorf("Claim wrong amount")
	}
	if tran.WelWalletAddr != welWalletAddr {
		return fmt.Errorf("Wrong claim wel wallet address")
	}
	_, err = e.EthCashoutWelTransDAO.GetClaimRequest(rqId)
	if err == sql.ErrNoRows {
		err := e.EthCashoutWelTransDAO.CreateClaimRequest(rqId, tran.ID, model.StatusPending)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if t.Status == "unconfirmed" {
		if tran.ClaimStatus != model.StatusPending { // Invalid state!
			err := e.EthCashoutWelTransDAO.UpdateClaimEthCashoutWel(tran.ID, rqId, model.StatusPending, t.Hash, fee, model.StatusPending)
			if err != nil {
				return err
			}
		}
	} else if t.Status == "confirmed" {
		if tran.ClaimStatus != model.StatusSuccess {
			err := e.EthCashoutWelTransDAO.UpdateClaimEthCashoutWel(tran.ID, rqId, model.RequestSuccess, t.Hash, fee, model.StatusSuccess)
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

	mkEventRecord := func(status string) (*model.WelEthEvent, error) {
		if status != model.StatusSuccess && status != model.StatusUnknown {
			return nil, fmt.Errorf(`Event status is neither "confirmed" nor "unconfirmed"`)
		}

		welTokenAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[1][12:]))
		event := &model.WelEthEvent{
			WelTokenAddr:  welTokenAddr,
			EthWalletAddr: ethWalletAddr.Hex(),
			NetworkID:     networkID.SetBytes(t.Log[logpos].Topics[3]).String(),
			DepositAt:     time.Now(),
		}
		//m, err := json.Marshal(event)
		//if err != nil {
		//	return nil, fmt.Errorf("can't gen id")
		//}
		//event.ID = crypto.Keccak256Hash(m).Big().String()

		event.DepositTxHash = t.Hash
		event.WelWalletAddr, _ = libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
		event.Amount = amount
		event.EthTokenAddr = model.EthTokenFromWel[welTokenAddr]
		event.DepositStatus = status
		event.Fee = fee

		return event, nil
	}

	switch t.Status {
	case "unconfirmed":
		// NOTE: if front end can't get txHash then we will need to fix this
		_, err := e.WelCashinEthTransDAO.SelectTransByDepositTxHash(t.Hash)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			// somehow, we did not save this deposit to db before
			event, err := mkEventRecord(model.StatusUnknown)
			if err != nil {
				logger.Get().Err(err).Msg("[DoneDeposit] can't create new event record")
				return err
			}

			err = e.WelCashinEthTransDAO.CreateWelCashinEthTrans(event)
			if err != nil {
				logger.Get().Err(err).Msg("[DoneDeposit] can't create new transaction")
				return err
			}
		}
		logger.Get().Info().Msg("[DoneDeposit] unconfirmed cashin transaction: " + t.Hash)

	case "confirmed":
		tran, err := e.WelCashinEthTransDAO.SelectTransByDepositTxHash(t.Hash)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			// somehow, we did not save this deposit to db before
			event, err := mkEventRecord(model.StatusSuccess)
			if err != nil {
				logger.Get().Err(err).Msg("[DoneDeposit] can't create new event record")
				return err
			}

			err = e.WelCashinEthTransDAO.CreateWelCashinEthTrans(event)
			if err != nil {
				logger.Get().Err(err).Msg("[DoneDeposit] can't create new transaction")
				return err
			}
		} else {
			if tran.DepositStatus != model.StatusSuccess {
				welWalletAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
				logger.Get().Info().Msg("[DoneDeposit] Update confirmed cashin transaction...")
				err := e.WelCashinEthTransDAO.UpdateDepositWelCashinEthConfirmed(t.Hash, welWalletAddr, amount, fee)
				if err != nil {
					return err
				}
			}
		}

	// emit done deposit event, save to db
	default:
		return fmt.Errorf("unknown status")
	}

	return nil
}
