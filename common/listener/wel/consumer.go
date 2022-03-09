package wel

import (
	"github.com/ethereum/go-ethereum/common"
)

type EventConsumer struct {
	Address    common.Address
	Topic      common.Hash
	ParseEvent EventParser
}

type IEventConsumer interface {
	GetConsumer() (*EventConsumer, error)
	//GetFilterQuery() ethereum.FilterQuery
}

type EventParser func(t *Transaction) error
