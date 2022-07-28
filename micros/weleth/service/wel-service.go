package service

import (
	"bridge/libs"
	coreEthService "bridge/micros/core/service/eth"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	GotronCommon "github.com/Paven-Org/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
)

type WelConsumer struct {
	ExportContractAddr    string
	ImportContractAddr    string
	WelCashinEthTransDAO  dao.IWelCashinEthTransDAO
	EthCashoutWelTransDAO dao.IEthCashoutWelTransDAO
	exportAbi             abi.ABI

	EthCashinWelTransDAO  dao.IEthCashinWelTransDAO
	WelCashoutEthTransDAO dao.IWelCashoutEthTransDAO
	importAbi             abi.ABI

	tempCli client.Client
}

func NewWelConsumer(iaddr, eaddr string, tempCli client.Client, daos *dao.DAOs) *WelConsumer {
	exportAbiJSON, err := os.Open("abi/wel/Export.json")
	if err != nil {
		panic(err)
	}

	defer exportAbiJSON.Close()

	exportabi, err := abi.JSON(exportAbiJSON)
	if err != nil {
		panic(err)
	}

	importAbiJSON, err := os.Open("abi/wel/Import.json")
	if err != nil {
		panic(err)
	}

	defer importAbiJSON.Close()

	importabi, err := abi.JSON(importAbiJSON)
	if err != nil {
		panic(err)
	}

	return &WelConsumer{
		ExportContractAddr:    eaddr,
		ImportContractAddr:    iaddr,
		WelCashinEthTransDAO:  daos.WelCashinEthTransDAO,
		EthCashoutWelTransDAO: daos.EthCashoutWelTransDAO,
		exportAbi:             exportabi,

		EthCashinWelTransDAO:  daos.EthCashinWelTransDAO,
		WelCashoutEthTransDAO: daos.WelCashoutEthTransDAO,
		importAbi:             importabi,

		tempCli: tempCli,
	}
}

func (e *WelConsumer) GetConsumer() ([]*welListener.EventConsumer, error) {
	return []*welListener.EventConsumer{
		{
			Address: e.ExportContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.exportAbi.Events["Withdraw"].Sig),
			),
			ParseEvent: e.DoneDepositParser,
		},
		{
			Address: e.ExportContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.exportAbi.Events["Returned"].Sig),
			),

			ParseEvent: e.DoneReturnParser,
		},
		{
			Address: e.ImportContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.importAbi.Events["Imported"].Sig),
			),

			ParseEvent: e.DoneImportedParser,
		},
		{
			Address: e.ImportContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(e.importAbi.Events["Withdraw"].Sig),
			),

			ParseEvent: e.DoneIWithdrawParser,
		},
	}, nil
}

func (e *WelConsumer) DoneIWithdrawParser(t *welListener.Transaction, logpos int) error {
	logger.Get().Info().Msgf("[DoneIWithdrawEV] IWithdraw event caught at block %d", t.BlockNumber)
	if t.Status != "confirmed" {
		logger.Get().Info().Msg("[DoneIWithdraw] unconfirmed transaction, skipped")
		return nil
	}
	data := make(map[string]interface{})
	e.importAbi.UnpackIntoMap(
		data,
		"Withdraw",
		t.Log[logpos].Data,
	)

	var tx model.WelCashoutEthTrans

	tx.EthWalletAddr = data["to"].(common.Address).Hex()
	tx.WelWalletAddr, _ = libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[2][12:]))

	tx.Amount = data["amount"].(*big.Int).String()
	tx.CommissionFee = data["fee"].(*big.Int).String()
	_total := big.NewInt(0).Add(data["amount"].(*big.Int), data["fee"].(*big.Int))
	tx.Total = _total.String()

	tx.WelTokenAddr, _ = libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[1][12:]))
	tx.EthTokenAddr = model.EthTokenFromWel[tx.WelTokenAddr]

	tx.NetworkID = (&big.Int{}).SetBytes(t.Log[logpos].Topics[3]).String()

	tx.CashoutStatus = model.WelCashoutEthConfirmed
	tx.DisperseStatus = model.WelCashoutEthUnconfirmed

	tx.WelWithdrawTxHash = t.Hash
	logger.Get().Info().Msgf("IWithdraw transaction to be created: %+v", tx)

	// save tx
	id, err := e.WelCashoutEthTransDAO.CreateWelCashoutEthTrans(&tx)
	if err != nil {
		logger.Get().Err(err).Msgf("[DoneIWithdraw] can't create W2E cashout transaction %s", t.Hash)
		return err
	}
	tx.ID = id

	// send signal to batch tx
	logger.Get().Info().Msgf("IWithdraw transaction to be sent to BatchDisperse: %+v", tx)
	err = e.tempCli.SignalWorkflow(context.Background(), coreEthService.BatchDisperseID, "", coreEthService.BatchDisperseSignal, tx)
	if err != nil {
		logger.Get().Err(err).Msgf("[DoneIWithdraw] Error sending BatchDisperseWF tx %+v", tx)
		return err
	}

	return nil
}

func (e *WelConsumer) DoneImportedParser(t *welListener.Transaction, logpos int) error {
	logger.Get().Info().Msgf("[ImportedEV] Imported event caught at block %d", t.BlockNumber)
	confirmStatus := model.EthCashinWelConfirmed
	if t.Status != model.EthCashinWelConfirmed {
		logger.Get().Info().Msg("[ImportedEV] Imported event unconfirmed")
		//return nil
		confirmStatus = model.EthCashinWelUnconfirmed
	}
	data := make(map[string]interface{})
	e.importAbi.UnpackIntoMap(
		data,
		"Imported",
		t.Log[logpos].Data,
	)

	welTokenAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(t.Log[logpos].Topics[1][12:]))

	_receivers := data["receivers"].([]common.Address)
	receivers := libs.Map(
		func(rawAddr common.Address) string {
			hexAddr := rawAddr.Hex()
			_, hexAddr, _ = strings.Cut(hexAddr, "0x")
			ret, _ := libs.HexToB58("0x41" + hexAddr)
			return ret
		},
		_receivers)
	logger.Get().Info().Msgf("Receivers: %+v", receivers)

	_amounts := data["amounts"].([]*big.Int)
	amounts := libs.Map(
		func(amount *big.Int) string {
			return amount.String()
		},
		_amounts)
	logger.Get().Info().Msgf("Amounts: %+v", amounts)

	fee := data["fee"].(*big.Int).String()
	logger.Get().Info().Msgf("total fee: %s", fee)

	amountOfReceiver := make(map[string]string)
	for i, receiver := range receivers {
		rAmount := amounts[i]
		amountOfReceiver[receiver] = rAmount
	}

	// the rest of this should probably be implemented as a temporal workflow instead, but meh
	trans, err := e.EthCashinWelTransDAO.SelectTransByIssueTxHash(t.Hash)
	if err != nil {
		logger.Get().Err(err).Msgf("[ImportedEV] error while retrieving transactions with issue txhash %s", t.Hash)
		return err
	}
	logger.Get().Info().Msgf("trans: %v", trans)

	for _, tran := range trans {
		if welTokenAddr != tran.WelTokenAddr {
			tran.WelTokenAddr = welTokenAddr
		}
		amount := &big.Int{}
		amount.SetString(amountOfReceiver[tran.WelWalletAddr], 10)
		tran.Amount = amount.String()

		total := big.NewInt(0)
		total.SetString(tran.Total, 10)

		fee := &big.Int{}
		fee = fee.Sub(total, amount)
		tran.CommissionFee = fee.String()

		tran.Status = confirmStatus
		tran.IssuedAt = sql.NullTime{Time: time.Now(), Valid: true}

		logger.Get().Info().Msgf("[ImportedEV] tran to be saved: %+v", tran)
		if err := e.EthCashinWelTransDAO.UpdateEthCashinWelTx(tran); err != nil {
			logger.Get().Err(err).Msgf("[ImportedEV] error while saving transaction %+v", tran)
		}

	}

	return nil
}

func (e *WelConsumer) DoneReturnParser(t *welListener.Transaction, logpos int) error {
	logger.Get().Info().Msgf("[ReturnedEV] Returned event caught at block %d", t.BlockNumber)
	data := make(map[string]interface{})
	e.exportAbi.UnpackIntoMap(
		data,
		"Returned",
		t.Log[logpos].Data,
	)

	rqId := data["requestId"].(*big.Int).String()
	//tokenAddr := data["token"].(common.Address)
	welWalletAddr, _ := libs.HexToB58("0x41" + GotronCommon.Bytes2Hex(data["user"].(common.Address).Bytes())) //FIX!
	amount := data["amount"].(*big.Int).String()
	fee := data["fee"].(*big.Int).String()
	total := big.NewInt(0)
	total = total.Add(data["amount"].(*big.Int), data["fee"].(*big.Int))

	tran, err := e.EthCashoutWelTransDAO.SelectTransByRqId(rqId)
	if err != nil {
		return err
	}
	if total.String() != tran.Amount { // if amount + fee != originally cashed out amount
		return fmt.Errorf("Claim wrong amount")
	}
	if tran.WelWalletAddr != welWalletAddr {
		return fmt.Errorf("Wrong claim wel wallet address")
	}
	_, err = e.EthCashoutWelTransDAO.GetClaimRequest(rqId)
	if err == sql.ErrNoRows {
		err := e.EthCashoutWelTransDAO.CreateClaimRequest(rqId, tran.ID, model.StatusPending, time.Now())
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if t.Status == "unconfirmed" {
		if tran.ClaimStatus != model.StatusPending { // Invalid state!
			err := e.EthCashoutWelTransDAO.UpdateClaimEthCashoutWel(tran.ID, rqId, model.StatusPending, t.Hash, amount, fee, model.StatusPending)
			if err != nil {
				return err
			}
		}
	} else if t.Status == "confirmed" {
		if tran.ClaimStatus != model.StatusSuccess {
			err := e.EthCashoutWelTransDAO.UpdateClaimEthCashoutWel(tran.ID, rqId, model.RequestSuccess, t.Hash, amount, fee, model.StatusSuccess)
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
	logger.Get().Info().Msgf("[EWithdrawEV] EWithdraw event caught at block %d", t.BlockNumber)
	data := make(map[string]interface{})
	e.exportAbi.UnpackIntoMap(
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
