package service

import (
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	welListener "bridge/service-managers/listener/wel"
	"fmt"
	"strconv"
)

type DoneDepositConsumer struct {
	ContractAddr  string
	VaultAddr     string
	WelDepositDAO dao.IWelTransDAO
}

func NewDoneDepositConsumer(addr, vaultAddr string, welDepositDAO dao.IWelTransDAO) *DoneDepositConsumer {
	return &DoneDepositConsumer{
		ContractAddr:  addr,
		VaultAddr:     vaultAddr,
		WelDepositDAO: welDepositDAO,
	}
}

func (e *DoneDepositConsumer) GetConsumer() *welListener.EventConsumer {
	return &welListener.EventConsumer{
		Address:    e.ContractAddr,
		ParseEvent: e.Parser,
	}
}

func (e *DoneDepositConsumer) Parser(t *welListener.Transaction) error {
	// filter if the to wallet is our vault wallet
	if t.Contract.Parameter.Raw["to_address"].(string) == e.VaultAddr {
		if t.Result == "unconfirmed" {
			tran, _ := e.WelDepositDAO.SelectTransByTxHash(t.Hash)
			// somehow, we did not save this deposit to db before
			if tran == nil {
				_ = e.WelDepositDAO.CreateTrans(&model.DoneDepositEvent{
					ID:        "", // TODO: use keccak256 later
					DepositID: "",
					TxHash:    t.Hash,
					FromAddr:  t.Contract.Parameter.Raw["owner_address"].(string),
					Amount:    strconv.FormatUint(t.Contract.Parameter.Raw["amount"].(uint64), 10),
					Decimal:   t.Contract.Parameter.Raw["decimal"].(uint),
				})
			}
		} else if t.Result == "confirmed" {
			err := e.WelDepositDAO.UpdateVerified(t.Hash)
			if err != nil {
				return err
			}
			// emit done deposit event, save to db
		} else {
			return fmt.Errorf("unknown status")
		}
	}

	return nil
}
