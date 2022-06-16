package eth

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"bridge/common/consts"
	"bridge/micros/weleth/model"
	"bridge/service-managers/daemon"
	"bridge/service-managers/logger"

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
	TxMonitors       map[common.Address]ITxMonitor
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
		TxMonitors:       make(map[common.Address]ITxMonitor),
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
		s.Logger.Err(err).Msg("[eth listener] Unable to get consumer")
		return err
	}
	for i := 0; i < len(consumerHandler); i++ {
		logger.Get().Debug().Msgf("Key for comsumer: %+v, +%v", KeyFromBEConsumer(consumerHandler[i].Address.Hex(), consumerHandler[i].Topic.Hex()))
		s.EventConsumerMap[KeyFromBEConsumer(consumerHandler[i].Address.Hex(), consumerHandler[i].Topic.Hex())] = consumerHandler[i]
	}

	s.EventFilters = append(s.EventFilters, consumer.GetFilterQuery()...)
	return nil
}

func (s *EthListener) RegisterTxMonitor(monitor ITxMonitor) error {
	if monitor == nil {
		err := fmt.Errorf("Nil monitor")
		s.Logger.Err(err).Msg("[eth listener] Register nil monitor")
		return err
	}
	address := monitor.MonitoredAddress()
	if _, ok := s.TxMonitors[address]; !ok {
		s.TxMonitors[address] = monitor
		s.Logger.Info().Msg("Monitor for " + fmt.Sprintf("0x%x", address) + " Added")
		return nil
	}
	return fmt.Errorf("Monitor for " + fmt.Sprintf("0x%x", address) + " already existed")
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
				logger.Get().Debug().Msgf("[eth_listner] handling event log: %+v", vLog)
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
func (s *EthListener) Scan(parentContext context.Context) (daemon consts.Daemon, err error) {
	// scanners

	blockTransferScanner := func(from *big.Int, to *big.Int) {
		for i := from; i.Cmp(to) < 1; i = i.Add(i, big.NewInt(1)) {
			logger.Get().Info().Msgf("[eth_listener] start scanning transfer in block %s", i.String())
			currBlock, err := s.EthClient.BlockByNumber(context.Background(), i)
			if err != nil {
				s.Logger.Err(err).Msg("[eth_listener]] Ethereum tx scan err")
				s.errC <- err
				continue
			}
			trans := currBlock.Transactions()
			for _, t := range trans {
				if t.To() == nil { // contract deployment, ignore
					continue
				}
				var isContractCall = false
				// temporary cutoff the ERC20 checking path because go-eth's performance is horrible
				if len(t.Data()) > 0 {
					//isContractCall = true
					continue
				}
				//receipt, err := s.EthClient.TransactionReceipt(context.Background(), t.Hash())
				//if err != nil {
				//	s.Logger.Err(err).Msg("[eth_listener]] Ethereum tx receipt retrieval err")
				//	continue
				//}
				//if receipt.Status != 1 {
				//	//s.Logger.Debug().Msgf("[eth_listener]] Skipping failed Ethereum tx %s...", t.Hash().Hex())
				//	continue
				//}
				// check if this is an ERC20 transfer
				//if len(t.Data()) > 0 {
				//	isContractCall = true
				//	abiJson := strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`)
				//	abi, _ := abi.JSON(abiJson)
				//	topic := crypto.Keccak256Hash([]byte(abi.Events["Transfer"].Sig))
				//	for _, log := range receipt.Logs {
				//		if len(log.Topics) < 3 {
				//			//s.Logger.Debug().Msgf("[eth listener] topics length too short in log: %+v\n", log)
				//			continue
				//		}
				//		if log.Topics[0] == topic {
				//			//s.Logger.Info().Msgf("[eth listener] ERC20 transfer tx detected: %+v\n", t)
				//			_to := common.HexToAddress(log.Topics[2].Hex())
				//			if monitor, ok := s.TxMonitors[_to]; ok {
				//				data := make(map[string]interface{})
				//				abi.UnpackIntoMap(data, "Transfer", log.Data)
				//				from := common.HexToAddress(log.Topics[1].Hex()).Hex()
				//				to := _to.Hex()
				//				contract := t.To().Hex()
				//				amount := data["amount"].(*big.Int).String()
				//				monitor.TxParse(t, from, to, contract, amount)
				//			}
				//		}
				//	}
				//}

				if monitor, ok := s.TxMonitors[*t.To()]; ok && !isContractCall {
					logger.Get().Debug().Msgf("tran to: %+v", *t.To())
					msg, err := t.AsMessage(types.LatestSignerForChainID(t.ChainId()), nil)
					if err != nil {
						s.Logger.Err(err).Msg("Unable to convert transaction to message")
						continue
					}
					from := msg.From().Hex()
					to := msg.To().Hex()
					amount := t.Value().String()
					monitor.TxParse(t, from, to, model.EthereumTk, amount)
				}
			}
		}
	}

	eventScanner := func(query ethereum.FilterQuery, from, to *big.Int) {
		query.FromBlock = from //scannedBlock
		query.ToBlock = to     //currBlock

		s.Logger.Debug().Msgf("[eth_listener] query %+v", query)

		// get all event
		events, err := s.EthClient.FilterLogs(context.Background(), query)
		if err != nil {
			s.Logger.Err(err).Msg("[eth_listener]] Ethereum filter query err")
			s.errC <- err
		}
		logger.Get().Debug().Msgf("[eth_listener] events: %+v", events)
		for _, event := range events {
			logger.Get().Debug().Msgf("[eth_listener] queuing event: %+v", event)
			s.Log <- event
		}
	}

	// main scanning daemon
	daemon = func() {
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
					continue
				}

				header, err := s.EthClient.HeaderByNumber(parentContext, nil)
				if err != nil {
					s.Logger.Err(err).Msg("[eth_listener] can't get head by number, possibly due to rpc node failure")
					continue
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

						// tx scan
						wg := sync.WaitGroup{}
						if len(s.TxMonitors) > 0 {
							wg.Add(1)
							go func(from *big.Int, to *big.Int) {
								blockTransferScanner(from, to)
								wg.Done()
							}(begin, until)
						}

						// events scan
						for _, query := range s.EventFilters {
							wg.Add(1)
							go func(query ethereum.FilterQuery) {
								eventScanner(query, begin, until)
								wg.Done()
							}(query)
						}
						wg.Wait()
						// update last scan block
						sysInfo.LastScannedBlock = until.Int64()
						s.EthInfo.Update(sysInfo)
					}
				} else {
					//s.Logger.Info().Msg(fmt.Sprintf("[eth_listener] scan from block %s to %s", scannedBlock.String(), currBlock.String()))
					// tx scan
					wg := sync.WaitGroup{}
					if len(s.TxMonitors) > 0 {
						wg.Add(1)
						go func(from *big.Int, to *big.Int) {
							blockTransferScanner(from, to)
							wg.Done()
						}(scannedBlock, currBlock)
					}

					// events scan
					for _, query := range s.EventFilters {
						wg.Add(1)
						go func(query ethereum.FilterQuery) {
							eventScanner(query, scannedBlock, currBlock)
							wg.Done()
						}(query)
					}
					wg.Wait()
					// update last scan block
					sysInfo.LastScannedBlock = currBlock.Int64()
					s.EthInfo.Update(sysInfo)
				}

				// TODO: either push this to delay message queue to run OR just sleep
				//consts.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
			}
		}
	}
	return daemon, nil
}

func (s *EthListener) matchEvent(vLog types.Log) (*EventConsumer, bool) {
	key := KeyFromBEConsumer(vLog.Address.Hex(), vLog.Topics[0].Hex())
	logger.Get().Debug().Msgf("key for event: %+v", key)
	consumer, isExisted := s.EventConsumerMap[key]

	if !isExisted {
		key = KeyFromBEConsumer(common.Address{}.Hex(), vLog.Topics[0].Hex())
		consumer, isExisted = s.EventConsumerMap[key]
	}
	logger.Get().Debug().Msgf("consumer if existed: %+v", consumer)

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
