package wel

import "github.com/ethereum/go-ethereum/common"

type EventConsumer struct {
	Address    string
	Topic      common.Hash
	ParseEvent EventParser
}

type IEventConsumer interface {
	GetConsumer() ([]*EventConsumer, error)
}

type EventParser func(t *Transaction) error
