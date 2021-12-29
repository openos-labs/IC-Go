package idl

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Int struct {
	Base uint8
	primType
}

func Int16() *Int {
	return &Int{
		Base: 16,
	}
}

func Int32() *Int {
	return &Int{
		Base: 32,
	}
}

func Int64() *Int {
	return &Int{
		Base: 64,
	}
}

func Int8() *Int {
	return &Int{
		Base: 8,
	}
}

func (n *Int) Decode(r *bytes.Reader) (interface{}, error) {
	if n.Base == 0 {
		return leb128.DecodeSigned(r)
	}
	return readInt(r, int(n.Base/8))
}

func (n Int) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	if n.Base == 0 {
		return leb128.EncodeSigned(big.NewInt(intType))
	}
	intXType := new(big.Int).Set(big.NewInt(intXType))
	intXType = intXType.Add(
		intXType,
		big.NewInt(3-int64(log2(n.Base))),
	)
	return leb128.EncodeSigned(intXType)
}

func (n Int) EncodeValue(v interface{}) ([]byte, error) {
	v_, ok := v.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	if n.Base == 0 {
		return leb128.EncodeSigned(v_)
	}
	{
		exp := big.NewInt(int64(n.Base) - 1)
		lim := big.NewInt(2)
		lim = lim.Exp(lim, exp, nil)
		min := new(big.Int).Set(lim)
		min = min.Mul(min, big.NewInt(-1))
		max := new(big.Int).Set(lim)
		max = max.Add(max, big.NewInt(-1))
		if v_.Cmp(min) < 0 || max.Cmp(v_) < 0 {
			return nil, fmt.Errorf("invalid value: %s", v_)
		}
	}
	return writeInt(v_, int(n.Base/8)), nil
}

func (n Int) String() string {
	if n.Base == 0 {
		return "int"
	}
	return fmt.Sprintf("int%d", n.Base)
}
