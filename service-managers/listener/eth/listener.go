package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"bridge/common/consts"
	"bridge/service-managers/daemon"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type EthListener struct {
	EthInfo          consts.IEthInfoRepo
	Log              chan types.Log
	EthClient        *ethclient.Client
	EventFilters     []ethereum.FilterQuery
	EventConsumerMap map[string]*EventConsumer
	Logger           *zerolog.Logger
	errC             chan error
	blockTime        uint64
	blockOffset      int64
}

func NewEthListener(
	ethInfo consts.IEthInfoRepo,
	ethClient *ethclient.Client,
	blockTime uint64,
	blockOffset int64,
	logger *zerolog.Logger,
) *EthListener {
	return &EthListener{
		EthInfo:          ethInfo,
		EthClient:        ethClient,
		EventConsumerMap: make(map[string]*EventConsumer),
		Log:              make(chan types.Log),
		errC:             make(chan error),
		Logger:           logger,
		blockTime:        blockTime,
		blockOffset:      blockOffset,
	}
}

func KeyFromBEConsumer(address string, topic string) string {
	return fmt.Sprintf("%s:%s", address, topic)
}

func (s *EthListener) AddFilterQuery(query ethereum.FilterQuery) {
	s.EventFilters = append(s.EventFilters, query)
}

func (s *EthListener) RegisterConsumer(consumer IEventConsumer) error {
	consumerHandler, err := consumer.GetConsumer()
	if err != nil {
		return err
	}
	for i := 0; i < len(consumerHandler); i++ {
		s.EventConsumerMap[KeyFromBEConsumer(consumerHandler[i].Address.Hex(), consumerHandler[i].Topic.Hex())] = consumerHandler[i]
	}

	s.EventFilters = append(s.EventFilters, consumer.GetFilterQuery())
	return nil
}

func (s *EthListener) Start(ctx context.Context) {
	daemon.BootstrapDaemons(ctx, s.Handling, s.Scan)
}

func (s *EthListener) Handling(parentContext context.Context) (fn consts.Daemon, err error) {
	fn = func() {
		s.Logger.Info().Msg("[eth_listener] Start handling")
		for {
			select {
			case err := <-s.errC:
				s.Logger.Err(err).Msg("[eth_listener] Ethereum client scan block err")

			case vLog := <-s.Log:
				go func(vLog types.Log) {
					s.consumeEvent(vLog)
				}(vLog)

			case <-parentContext.Done():
				s.Logger.Info().Msg("[eth_listener] Blockchain listener stop")
				return
			}
		}

	}
	return fn, nil
}

// offset is the number of block to scan before current block to make sure event is confirmed
func (s *EthListener) Scan(parentContext context.Context) (fn consts.Daemon, err error) {
	fn = func() {
		s.Logger.Info().Msgf("[eth listener] Begin scan...")
		for {
			select {
			case <-parentContext.Done():
				s.Logger.Info().Msg("[eth listener] Stop listening")
				return
			default:
				sysInfo, err := s.EthInfo.Get()
				if err != nil {
					s.Logger.Err(err).Msg("[eth_listener] can't get system info")
				}

				header, err := s.EthClient.HeaderByNumber(parentContext, nil)
				if err != nil {
					s.Logger.Err(err).Msg("[eth_listener] can't get head by number")
				}
				currBlock := header.Number

				var scannedBlock *big.Int
				if sysInfo.LastScannedBlock <= 0 {
					// set first block
					scannedBlock = big.NewInt(s.blockOffset)
				} else {
					scannedBlock = big.NewInt(sysInfo.LastScannedBlock)
				}

				// scanned a offset - 1 block before to sure event confirmed
				scannedBlock = scannedBlock.Sub(scannedBlock, big.NewInt(s.blockOffset-1))

				// if last scanned block is more than $BIGNUM blocks away just scan last $BIGNUM blocks
				diff := big.NewInt(0).Sub(currBlock, scannedBlock)
				if diff.Cmp(big.NewInt(600000)) > 0 {
					diff = big.NewInt(600000)
					scannedBlock = scannedBlock.Sub(currBlock, diff)
				}

				if diff.Cmp(big.NewInt(100)) > 0 {
					// scan in 100-chunk
					for begin := scannedBlock; currBlock.Cmp(begin) > 0; begin = begin.Add(begin, big.NewInt(100)) {
						limit := big.NewInt(99)
						curDiff := big.NewInt(0)
						if curDiff.Sub(currBlock, begin).Cmp(limit) < 0 {
							limit = currBlock.Sub(currBlock, begin)
						}
						until := big.NewInt(0)
						until = until.Add(begin, limit)
						//s.Logger.Info().Msg(fmt.Sprintf("[eth_listener] scan from block %s to %s", begin.String(), until.String()))

						for _, query := range s.EventFilters {
							go func(query ethereum.FilterQuery) {
								query.FromBlock = begin
								query.ToBlock = until

								//s.Logger.Debug().Msg(fmt.Sprintf("[eth_listener] query %v", query))

								// get all event
								events, err := s.EthClient.FilterLogs(context.Background(), query)
								if err != nil {
									s.Logger.Err(err).Msg("[eth_listener]] Ethereum filter query err")
									s.errC <- err
								}
								for _, event := range events {
									s.Log <- event
								}
							}(query)
						}
						// update last scan block
						sysInfo.LastScannedBlock = until.Int64()
						s.EthInfo.Update(sysInfo)
					}
				} else {
					//s.Logger.Info().Msg(fmt.Sprintf("[eth_listener] scan from block %s to %s", scannedBlock.String(), currBlock.String()))
					for _, query := range s.EventFilters {
						go func(query ethereum.FilterQuery) {
							query.FromBlock = scannedBlock
							query.ToBlock = currBlock

							//s.Logger.Debug().Msg(fmt.Sprintf("[eth_listener] query %v", query))

							// get all event
							events, err := s.EthClient.FilterLogs(context.Background(), query)
							if err != nil {
								s.Logger.Err(err).Msg("[eth_listener]] Ethereum filter query err")
								s.errC <- err
							}
							for _, event := range events {
								s.Log <- event
							}
						}(query)
					}
					// update last scan block
					sysInfo.LastScannedBlock = currBlock.Int64()
					s.EthInfo.Update(sysInfo)
				}

				// TODO: either push this to delay message queue to run OR just sleep
				consts.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
			}
		}
	}
	return fn, nil
}

func (s *EthListener) matchEvent(vLog types.Log) (*EventConsumer, bool) {
	key := KeyFromBEConsumer(vLog.Address.Hex(), vLog.Topics[0].Hex())
	consumer, isExisted := s.EventConsumerMap[key]

	if !isExisted {
		key = KeyFromBEConsumer(common.Address{}.Hex(), vLog.Topics[0].Hex())
		consumer, isExisted = s.EventConsumerMap[key]
	}

	return consumer, isExisted
}

func (s *EthListener) consumeEvent(vLog types.Log) {
	consumer, isExisted := s.matchEvent(vLog)
	if isExisted {
		err := consumer.ParseEvent(vLog)
		if err != nil {
			s.Logger.Err(err).Msg("[eth_client] Consume event error")
		}
	}
}
