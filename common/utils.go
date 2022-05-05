package common

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func EventSignatureHash(funSignature string) ethcommon.Hash {
	return crypto.Keccak256Hash([]byte(funSignature))
}

func StringToInt(amountStr string) (*big.Int, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if amount == nil || !ok {
		return nil, errors.New("string to int error" + amountStr)
	}
	return amount, nil
}

func BigToRound16(num *big.Int) *big.Int {
	n := new(big.Int).Div(num, big.NewInt(1e+16))
	nPlus := new(big.Int).Add(n, big.NewInt(int64(1)))
	return new(big.Int).Mul(big.NewInt(1e+16), nPlus)
}

func BigToFloatStr(num *big.Int) string {
	n := new(big.Int).Div(num, big.NewInt(1e+18))
	mod := new(big.Int).Mod(num, big.NewInt(1e+18))
	modInt := new(big.Int).Div(mod, big.NewInt(1e+16))
	nStr := n.String() + "." + modInt.String()
	return nStr
}

func VerifySig(msg string, sigHex string, addr ethcommon.Address) bool {
	data := []byte(msg)
	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		return false
	}
	if len(sig) != crypto.SignatureLength {
		return false
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return false
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return false
	}
	signer := crypto.PubkeyToAddress(*rpk)
	if signer != addr {
		return false
	}
	return true
}

func Sig2String(sig []byte) string {
	return "0x" + ethcommon.Bytes2Hex(sig)
}
