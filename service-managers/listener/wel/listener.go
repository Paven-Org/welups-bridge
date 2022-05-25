package wel

import (
	"context"
	"fmt"
	"time"

	"bridge/common/consts"
	"bridge/service-managers/daemon"

	GotronCommon "github.com/Paven-Org/gotron-sdk/pkg/common"
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
		s.Logger.Err(err).Msg("[wel listener] Unable to get consumer")
		return err
	}

	for i := 0; i < len(consumerHandler); i++ {
		// remove 0x from the topic
		s.EventConsumerMap[KeyFromBEConsumer(consumerHandler[i].Address, consumerHandler[i].Topic.Hex()[2:])] = consumerHandler[i]
	}
	return nil
}

func KeyFromBEConsumer(address string, topic string) string {
	return fmt.Sprintf("%s:%s", address, topic)
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
				s.Logger.Info().Msg("[wel_listener] Blockchain listener stop")
				return
			}
		}

	}
	return fn, nil
}

// offset is the number of block to scan before current block to make sure event is confirmed
func (s *WelListener) Scan(parentContext context.Context) (fn consts.Daemon, err error) {
	fn = func() {
		s.Logger.Info().Msgf("[wel listener] Begin scan...")
		for {
			select {
			case <-parentContext.Done():
				s.Logger.Info().Msg("[wel listener] Stop listening")
				return
			default:
				sysInfo, err := s.WelInfo.Get()
				if err != nil {
					s.Logger.Err(err).Msg("[wel_listener] can't get system info")
					continue
				}
				lastScanned := sysInfo.LastScannedBlock
				//s.Logger.Info().Msgf("[wel listener] sysinfo: %d", lastScanned)

				header, err := s.TransHandler.Client.GetNowBlock()
				if err != nil {
					s.Logger.Err(err).Msg("[wel_listener] can't get head by number, possibly due to rpc node failure")
					continue
				}
				headNum := header.BlockHeader.RawData.Number
				//s.Logger.Info().Msgf("[wel listener] scan from lastScanned %d to headNum %d", lastScanned, headNum)

				brange := headNum - lastScanned + 1
				if brange > 600000 {
					brange = 600000
					lastScanned = headNum - brange + 1
				}
				//s.Logger.Info().Msgf("[wel listener] block range: %d", brange)
				if brange > 100 {
					// partition the range to 100-long chunks
					for begin := lastScanned; headNum-begin > 0; begin += 100 {
						// probably should fire 1 goroutine for each partition, but i'm not sure gotron
						// sdk's client is threadsafe or not
						var limit int64
						if headNum-begin < 99 {
							limit = headNum - begin
						} else {
							limit = 99
						}
						//fmt.Println("from: ", begin, " to: ", begin+limit)
						s.TransHandler.GetInfoListTransactionRange(begin+limit, limit+1, "", s.Trans, s.errC)
						sysInfo.LastScannedBlock = begin + 99
						s.WelInfo.Update(sysInfo)
					}
				} else {
					//fmt.Println("from: ", headNum-brange, " to: ", headNum)
					s.TransHandler.GetInfoListTransactionRange(headNum, brange, "", s.Trans, s.errC)
					sysInfo.LastScannedBlock = headNum
					s.WelInfo.Update(sysInfo)
				}

				// TODO: either push this to delay message queue to run OR just sleep
				consts.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
			}
		}
	}
	return fn, nil
}

func (s *WelListener) matchEvent(tran *Transaction) (consumer *EventConsumer, position int) {
	//TODO: add topic
	position = -1
	ctrAddress := tran.Contract.Parameter.Raw["ContractAddress"]
	if ctrAddress == nil {
		return nil, -1
	}
	for i, log := range tran.Log {
		//fmt.Println("[matchEvent] at log ", i)
		key := KeyFromBEConsumer(ctrAddress.(string), GotronCommon.Bytes2Hex(log.Topics[0]))
		//fmt.Println("[matchEvent] key: ", key)
		consumer, isExisted := s.EventConsumerMap[key]
		if isExisted {
			return consumer, i
		}
	}

	return nil, -1
}

func (s *WelListener) consumeEvent(t *Transaction) {
	consumer, position := s.matchEvent(t)
	if position >= 0 {
		err := consumer.ParseEvent(t, position)
		if err != nil {
			s.Logger.Err(err).Msgf("[wel_listener] Consume event error, tx with event: %v", t)
		}
	}
}
