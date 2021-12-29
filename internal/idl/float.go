package idl

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/aviate-labs/leb128"
)

type Float struct {
	Base uint8
	primType
}

func Float32() *Float {
	return &Float{
		Base: 32,
	}
}

func Float64() *Float {
	return &Float{
		Base: 64,
	}
}

func (f *Float) Decode(r *bytes.Reader) (interface{}, error) {
	return readFloat(r, int(f.Base/8))
}

func (f Float) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	floatXType := new(big.Int).Set(big.NewInt(floatXType))
	if f.Base == 64 {
		floatXType.Add(floatXType, big.NewInt(-1))
	}
	return leb128.EncodeSigned(floatXType)
}

func (f Float) EncodeValue(v interface{}) ([]byte, error) {
	v_, ok := v.(*big.Float)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	return writeFloat(v_, int(f.Base/8)), nil
}

func (f Float) String() string {
	return fmt.Sprintf("float%d", f.Base)
}
