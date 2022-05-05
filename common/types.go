package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractInfo struct {
	Address   ethcommon.Address
	TopicList []Topic
}

type ChainInfo struct {
	RPC            string
	ChainID        int
	FilterContract []ContractInfo
}

type Event struct {
	Type        string
	TxHash      string
	Timestamps  int64
	Sender      ethcommon.Address
	BlockHeight uint64
	Data        interface{}
}

type EventInfo struct {
	Name  string
	Sig   []byte
	Event interface{}
}

type Topic interface {
	GetName() string
	GetSignature() ethcommon.Hash
	Unpack(log types.Log) (interface{}, error)
}