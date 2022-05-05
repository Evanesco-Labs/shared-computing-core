package emitter

import (
	"context"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"shared-computing-core/common"
	"shared-computing-core/db"
	"time"
)

type EventTracer struct {
	db        *db.DBService
	dur       time.Duration
	maxHeight int
}

func InitEventTracer(db *db.DBService, dur time.Duration, maxHeight int) (*EventTracer, error) {
	return &EventTracer{
		db:        db,
		dur:       dur,
		maxHeight: maxHeight,
	}, nil
}

func (s *EventTracer) SubscribeChainEvent(chainInfo common.ChainInfo, ch chan common.Event) {
	toHeight := uint64(0)
	duration := s.dur
	timer := time.NewTimer(duration)
	log.Info("subscribe chain events", "chainId", chainInfo.ChainID)
	for {
		select {
		case <-timer.C:
			timer.Reset(duration)
			client, err := ethclient.Dial(chainInfo.RPC)
			if err != nil {
				log.Error("ethclient dial error", "err", err)
				continue
			}
			currentHeight, err := client.BlockNumber(context.Background())
			if err != nil {
				log.Error("client blocknumber error", "err", err)
				continue
			}
			lastHeight, err := s.db.GetEventTraceHeight(chainInfo.ChainID)
			if err != nil {
				log.Error("GetLastEventeHeight error", "err", err)
				continue
			}
			fromHeight := lastHeight + 1
			if currentHeight < uint64(fromHeight) {
				log.Info("check blocks", "current", currentHeight, "from", fromHeight)
				continue
			}
			if (currentHeight - uint64(fromHeight)) > uint64(s.maxHeight) {
				toHeight = uint64(fromHeight) + uint64(s.maxHeight)
			} else {
				toHeight = currentHeight
			}
			log.Info("check blocks", "from", fromHeight, "to", toHeight)
			events, err := getEvents(chainInfo, big.NewInt(fromHeight), big.NewInt(int64(toHeight)))
			if err != nil {
				log.Error("getActivityEvent error", err)
				continue
			}
			for _, event := range events {
				ch <- event
			}
			s.db.PutEventTraceHeight(chainInfo.ChainID, int64(toHeight))
		}
	}
}

func getEvents(chainInfo common.ChainInfo, fromBlock *big.Int, toBlock *big.Int) ([]common.Event, error) {
	contractAddrList := make([]ethcommon.Address, 0)
	for _, contractInfo := range chainInfo.FilterContract {
		contractAddrList = append(contractAddrList, contractInfo.Address)
	}
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: contractAddrList,
	}
	client, err := ethclient.Dial(chainInfo.RPC)
	if err != nil {
		return nil, err
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("FilterLogs error", "err", err)
		return nil, err
	}
	eventList := make([]common.Event, 0)

	for _, vLog := range logs {
		if len(vLog.Topics) == 0 {
			log.Error("the length of Topics is less than 1")
			continue
		}
		found := false
		for _, contractInfo := range chainInfo.FilterContract {
			if found {
				break
			}
			for _, topic := range contractInfo.TopicList {
				if found {
					break
				}
				if vLog.Topics[0] == topic.GetSignature() {
					found = true
					inter, unpackError := topic.Unpack(vLog)
					if unpackError != nil {
						log.Error("UnpackIntoInterface error", unpackError)
						break
					}
					sender := ethcommon.Address{}
					if len(vLog.Topics) > 1 {
						sender = ethcommon.HexToAddress(vLog.Topics[1].Hex())
					}
					eventList = append(eventList, common.Event{
						topic.GetName(),
						vLog.TxHash.String(),
						time.Now().Unix(),
						sender,
						vLog.BlockNumber,
						inter,
					})
				}
			}
		}
	}
	return eventList, nil
}
