package idl

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Empty struct {
	primType
}

func (Empty) Decode(*bytes.Reader) (interface{}, error) {
	return nil, fmt.Errorf("cannot decode empty type")
}

func (Empty) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	return leb128.EncodeSigned(big.NewInt(emptyType))
}

func (Empty) EncodeValue(_ interface{}) ([]byte, error) {
	return []byte{}, nil
}

func (Empty) String() string {
	return "empty"
}
