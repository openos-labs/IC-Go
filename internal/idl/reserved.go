package idl

import (
	"bytes"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Reserved struct {
	primType
}

func (Reserved) Decode(*bytes.Reader) (interface{}, error) {
	return nil, nil
}

func (Reserved) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	return leb128.EncodeSigned(big.NewInt(reservedType))
}

func (Reserved) EncodeValue(_ interface{}) ([]byte, error) {
	return []byte{}, nil
}

func (Reserved) String() string {
	return "reserved"
}
