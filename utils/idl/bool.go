package idl

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Bool struct {
	primType
}

func (b Bool) Decode(r *bytes.Reader) (interface{}, error) {
	v, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch v {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return nil, fmt.Errorf("invalid bool values: %x", b)
	}
}

func (Bool) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	return leb128.EncodeSigned(big.NewInt(boolType))
}

func (b Bool) EncodeValue(v interface{}) ([]byte, error) {
	v_, ok := v.(bool)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	if v_ {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

func (Bool) String() string {
	return "bool"
}
