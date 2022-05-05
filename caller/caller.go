package caller

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"shared-computing-core/common"
	"time"
)

var TxTimeOutError = errors.New("transaction pending timeout, pls retry")

type Caller struct {
	reqCh chan common.CallerReq
	priv  []*ecdsa.PrivateKey
}

func InitCaller(privStrList []string) (*Caller, error) {
	caller := Caller{
		reqCh: make(chan common.CallerReq),
		priv:  make([]*ecdsa.PrivateKey, 0),
	}
	for _, privStr := range privStrList {
		priv, err := crypto.HexToECDSA(privStr)
		if err != nil {
			return nil, err
		}
		caller.priv = append(caller.priv, priv)
	}
	return &caller, nil
}

func (c *Caller) GetRequestCh() chan common.CallerReq {
	return c.reqCh
}

func (c *Caller) Start() {
	for _, priv := range c.priv {
		c.StartPrivWorker(priv)
	}
}

func (c *Caller) StartPrivWorker(priv *ecdsa.PrivateKey) {
	go func() {
		for {
			req := <-c.reqCh
			switch req.Type {
			}
			time.Sleep(time.Second)
		}
	}()
}
