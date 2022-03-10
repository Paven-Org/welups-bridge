package wel

type EventConsumer struct {
	Address    string
	ParseEvent EventParser
}

type IEventConsumer interface {
	GetConsumer() (*EventConsumer, error)
}

type EventParser func(t *Transaction) error
