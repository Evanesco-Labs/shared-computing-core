package emitter

import (
	"errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"shared-computing-core/common"
	"shared-computing-core/db"
)

var (
	NotOwnerError   = errors.New("not owner")
	StillInUseError = errors.New("nft still in use")
)

type RpcHandler struct {
	db       *db.DBService
	callerCh chan common.CallerReq
}

func InitRpcHandler(db *db.DBService, callerCh chan common.CallerReq) *RpcHandler {
	return &RpcHandler{
		db:       db,
		callerCh: callerCh,
	}
}

func (rh *RpcHandler) RecordTx(hash string, chainId int, addr ethcommon.Address, desc string) {
	sqlStr := ""
	_, err := rh.db.DB.Exec(sqlStr, hash, addr.String(), chainId, desc)
	if err != nil {
		log.Error("record rpc call tx error", "err", err)
	}
	return
}
