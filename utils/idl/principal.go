package idl

import (
	"bytes"
	"fmt"
	"github.com/openos-labs/IC-Go/utils/principal"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Principal struct {
	primType
}

func (t *Principal) Decode(r *bytes.Reader) (interface{}, error) {
	{
		bs := make([]byte, 1)
		n, err := r.Read(bs)
		if err != nil {
			return nil, err
		}
		if n != 1 || bs[0] != 0x01 {
			return nil, fmt.Errorf("invalid func reference: %d", bs)
		}
	}
	l, err := leb128.DecodeUnsigned(r)
	if err != nil {
		return nil, err
	}
	pid := make([]byte, l.Int64())
	n, err := r.Read(pid)
	if err != nil {
		return nil, err
	}
	if n != int(l.Int64()) {
		return nil, fmt.Errorf("invalid principal id: %d", pid)
	}
	return pid, nil
}

func (t Principal) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	return leb128.EncodeSigned(big.NewInt(principalType))
}

func (t Principal) EncodeValue(v interface{}) ([]byte, error) {
	p, ok := v.(principal.Principal)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	l, err := leb128.EncodeUnsigned(big.NewInt(int64(len(p))))
	if err != nil {
		return nil, err
	}
	return concat([]byte{0x01}, l, p), nil
}

func (t Principal) String() string {
	return "principal"
}
