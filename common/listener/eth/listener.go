package eth

import (
	"context"
	"fmt"
	"math/big"

	"bridge/common/consts"

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
}

func NewEthListener(
	ethInfo consts.IEthInfoRepo,
	ethClient ethclient.Client,
	blockTime uint64,
	logger *zerolog.Logger,
) *EthListener {
	return &EthListener{
		EthInfo:          ethInfo,
		EthClient:        &ethClient,
		EventConsumerMap: make(map[string]*EventConsumer),
		Log:              make(chan types.Log),
		//OpChan:          make(chan SubStatus),
		Logger:    logger,
		blockTime: blockTime,
	}
}

func (s *EthListener) AddFilterQuery(query ethereum.FilterQuery) {
	s.EventFilters = append(s.EventFilters, query)
}

func (s *EthListener) RegisterConsumer(consumer IEventConsumer) error {
	consumerHandler, err := consumer.GetConsumer()
	if err != nil {
		return err
	}
	s.EventConsumerMap[KeyFromBEConsumer(consumerHandler.Address.Hex(), consumerHandler.Topic.Hex())] = consumerHandler

	s.EventFilters = append(s.EventFilters, consumer.GetFilterQuery())
	return nil
}

func (s *EthListener) Start() error {
	if err := s.Scan(5); err != nil {
		s.Logger.Err(err).Msg("[eth_listener] Etherum client scan block err")
		return err
	}

	s.Logger.Info().Msg("[eth_listener] Start listening")
	for {
		select {
		case err := <-s.errC:
			s.Logger.Err(err).Msg("[eth_listener] Etherum client scan block err")

		case vLog := <-s.Log:
			go func(vLog types.Log) {
				s.consumeEvent(vLog)
			}(vLog)

			/*
				case status := <-s.OpChan:
					if status == Stop {
						s.Logger.Info("Blockchain listener stop")
						return nil
					}
				}
			*/
		}
	}
}
func KeyFromBEConsumer(address string, topic string) string {
	return fmt.Sprintf("%s:%s", address, topic)
}

// offset is the number of block to scan before current block to make sure event is confirmed
func (s *EthListener) Scan(offset int64) error {
	sysInfo, err := s.EthInfo.Get()
	if err != nil {
		return err
	}

	header, err := s.EthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	currBlock := header.Number.String()

	scannedBlock := &big.Int{}
	if sysInfo.LastScannedBlock == 0 {
		// set first block
		ok := true
		scannedBlock, ok = scannedBlock.SetString(currBlock, 10)
		if !ok {
			panic("can't parse current block number")
		}
	} else {
		scannedBlock = big.NewInt(sysInfo.LastScannedBlock)
	}
	// scanned a few block before to make sure event confirmed
	scannedBlock.Add(scannedBlock, big.NewInt(1))
	scannedBlock.Sub(scannedBlock, big.NewInt(offset))

	// if last event is more than 10k block away just scan last 10k block
	diff := big.NewInt(0).Sub(header.Number, scannedBlock)
	if diff.Cmp(big.NewInt(10000)) > 0 {
		scannedBlock.Sub(header.Number, big.NewInt(10000))
	}

	//s.Logger.Infof("scan from block %s to %s", scannedBlock.String(), header.Number.String())

	for _, query := range s.EventFilters {
		go func(query ethereum.FilterQuery) {
			query.FromBlock = scannedBlock
			query.ToBlock = header.Number

			//s.Logger.Debugf("subscriber query: %+v", query)

			// get all missed event
			missEvents, err := s.EthClient.FilterLogs(context.Background(), query)
			if err != nil {
				//s.Logger.Warn(err)
				s.errC <- err
			}
			for _, event := range missEvents {
				s.Log <- event
			}
		}(query)
	}

	// update last scan block as current block instead of last event block like before
	sysInfo.LastScannedBlock = header.Number.Int64()
	s.EthInfo.Update(sysInfo)

	// TODO: either push this to delay message queue OR just sleep
	//nextScanTime := time.Now().UTC().Add(time.Second * time.Duration(s.blockTime))

	return nil
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
