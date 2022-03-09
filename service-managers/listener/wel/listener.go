package wel

import (
	"context"
	"fmt"
	"time"

	"bridge/common/consts"
	"bridge/service-managers/daemon"

	GotronCommon "github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/rs/zerolog"
)

type WelListener struct {
	TransHandler     *TransHandler
	WelInfo          consts.IWelInfoRepo
	EventConsumerMap map[string]*EventConsumer
	Logger           *zerolog.Logger
	Trans            chan *Transaction
	errC             chan error
	blockTime        uint64
	blockOffset      int64
}

func NewWelListener(
	welInfo consts.IWelInfoRepo,
	transHandler *TransHandler,
	blockTime uint64,
	blockOffset int64,
	logger *zerolog.Logger,
) *WelListener {
	return &WelListener{
		WelInfo:          welInfo,
		TransHandler:     transHandler,
		EventConsumerMap: make(map[string]*EventConsumer),
		Trans:            make(chan *Transaction),
		errC:             make(chan error),
		Logger:           logger,
		blockTime:        blockTime,
		blockOffset:      blockOffset,
	}
}

func (s *WelListener) RegisterConsumer(consumer IEventConsumer) error {
	consumerHandler, err := consumer.GetConsumer()
	if err != nil {
		return err
	}

	s.EventConsumerMap[KeyFromBEConsumer(GotronCommon.EncodeCheck(consumerHandler.Address))] = consumerHandler
	return nil
}

func KeyFromBEConsumer(address string) string {
	return fmt.Sprintf("%s", address)
}

func (s *WelListener) Start(ctx context.Context) {
	daemon.BootstrapDaemons(ctx, s.Handling, s.Scan)
}

func (s *WelListener) Handling(parentContext context.Context) (fn consts.Daemon, err error) {
	fn = func() {
		s.Logger.Info().Msg("[wel_listener] Start handling")
		for {
			select {
			case err := <-s.errC:
				s.Logger.Err(err).Msg("[wel_listener] client scan block err")

			case vLog := <-s.Trans:
				go func(t *Transaction) {
					s.consumeEvent(vLog)
				}(vLog)

			case <-parentContext.Done():
				s.Logger.Info().Msg("[eth_listener] Blockchain listener stop")
			}
		}

	}
	return fn, nil
}

// offset is the number of block to scan before current block to make sure event is confirmed
func (s *WelListener) Scan(parentContext context.Context) (fn consts.Daemon, err error) {
	fn = func() {
		sysInfo, err := s.WelInfo.Get()
		if err != nil {
			s.Logger.Err(err).Msg("[wel_listener] can't get system info")
		}

		header, err := s.TransHandler.Client.GetNowBlock()
		if err != nil {
			s.Logger.Err(err).Msg("[wel_listener] can't get head by number")
		}
		headNum := header.BlockHeader.RawData.Number

		s.Logger.Info().Msg(fmt.Sprintf("[wel_listener] scan from block %v to %v", sysInfo.LastScannedBlock, headNum))
		if headNum-sysInfo.LastScannedBlock > 10000 {
			s.TransHandler.GetInfoListTransactionRange(headNum, 10000, "", s.Trans, s.errC)
		} else {
			s.TransHandler.GetInfoListTransactionRange(headNum, s.blockOffset, "", s.Trans, s.errC)
		}

		// update last scan block
		sysInfo.LastScannedBlock = headNum
		s.WelInfo.Update(sysInfo)

		// TODO: either push this to delay message queue to run OR just sleep
		consts.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
	}
	return fn, nil
}

func (s *WelListener) matchEvent(tran *Transaction) (*EventConsumer, bool) {
	key := KeyFromBEConsumer(tran.ContractAddress)
	consumer, isExisted := s.EventConsumerMap[key]
	if isExisted {
		return consumer, isExisted
	}

	return nil, false
}

func (s *WelListener) consumeEvent(t *Transaction) {
	consumer, isExisted := s.matchEvent(t)
	if isExisted {
		err := consumer.ParseEvent(t)
		if err != nil {
			s.Logger.Err(err).Msg("[wel_listener] Consume event error")
		}
	}
}
