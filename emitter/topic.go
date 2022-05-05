package emitter

import (
	"errors"
)

var (
	TopicNotEnoughError = errors.New("topic not enough")
)

type AbiUnPacker interface {
	UnpackIntoInterface(v interface{}, name string, data []byte) error
}
