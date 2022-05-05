package emitter

import (
	"database/sql"
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"reflect"
	"shared-computing-core/caller"
	"shared-computing-core/common"
	"shared-computing-core/db"
)

var (
	TxExecutedError            = errors.New("transaction event already executed")
	TxEventAlreadyHandledError = errors.New("tx event already handled")
	TxSendFailedError          = errors.New("send transaction failed")
	EventParamsNil             = errors.New("event param nil")
)

type EventHandler struct {
	db          *db.DBService
	ChainList   []common.ChainInfo
	callerCh    chan common.CallerReq
	caller      *caller.Caller
	eventTracer *EventTracer
}

func InitEventHandler(db *db.DBService, chainList []common.ChainInfo, tracer *EventTracer) *EventHandler {
	return &EventHandler{
		db:          db,
		ChainList:   chainList,
		eventTracer: tracer,
	}
}

func (eh *EventHandler) Start() {
	eventCh := make(chan common.Event)
	for _, chainInfo := range eh.ChainList {
		go eh.eventTracer.SubscribeChainEvent(chainInfo, eventCh)
	}
	go func() {
		for {
			event := <-eventCh
			log.Info("try handle event", "name", event.Type)
			InvokeObjectMethod(eh, event.Type, event)
		}
	}()
}

func InvokeObjectMethod(object interface{}, methodname string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range inputs {
		inputs[i] = reflect.ValueOf(args[i])
	}
	method := reflect.ValueOf(object).MethodByName(methodname)
	if method.IsValid() {
		method.Call(inputs)
	} else {
		log.Error("no handle for event", "event", methodname)
	}
}

func (h *EventHandler) IfTxExecuted(txHash string, chainId int) bool {
	sqlStr := ""
	row := h.db.DB.QueryRow(sqlStr, txHash, chainId)
	id := int64(0)
	err := row.Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (h *EventHandler) IfTxEventHandled(txHash string, chainId int, eventType string) bool {
	sqlStr := ""
	row := h.db.DB.QueryRow(sqlStr, txHash, chainId, eventType)
	id := int64(0)
	err := row.Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (h *EventHandler) RecordTxHandle(fromHash string, toHash string, chainId int, status int, errStr string, eventType string) error {
	sqlStr := ""
	_, err := h.db.DB.Exec(sqlStr, fromHash, toHash, chainId, status, errStr, eventType, toHash, status, errStr)
	if err != nil {
		log.Error("RecordTxHandle error", "err", err)
	}
	return nil
}
