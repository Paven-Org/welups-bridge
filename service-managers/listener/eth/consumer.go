package eth

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventConsumer struct {
	Address    common.Address
	Topic      common.Hash
	ParseEvent EventParser
}

type IEventConsumer interface {
	GetConsumer() ([]*EventConsumer, error)
	GetFilterQuery() ethereum.FilterQuery
}

type EventParser func(types.Log) error

//-----------------------------------------------------------//
type ITxMonitor interface {
	MonitoredAddress() common.Address
	TxParse(t *types.Transaction, from, to, tokenAddr, amount string) error
}
