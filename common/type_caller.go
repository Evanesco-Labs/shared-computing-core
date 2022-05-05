package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type CallerResp struct {
	Status int
	TxHash ethcommon.Hash
	Err    string
}

type CallerReq struct {
	Type  int
	Resp  chan CallerResp
	Param interface{}
}
