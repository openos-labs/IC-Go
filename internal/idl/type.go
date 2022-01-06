package idl

import (
	"bytes"
	"fmt"
)

var (
	nullType      int64 = -1  // 0x7f
	boolType      int64 = -2  // 0x7e
	natType       int64 = -3  // 0x7d
	intType       int64 = -4  // 0x7c
	natXType      int64 = -5  // 0x7b-0x78
	intXType      int64 = -9  // 0x77-0x73
	floatXType    int64 = -13 // 0x72
	textType      int64 = -15 // 0x71
	reservedType  int64 = -16 // 0x70
	emptyType     int64 = -17 // 0x6f
	optType       int64 = -18 // 0x6e
	vecType       int64 = -19 // 0x6d
	recType       int64 = -20 // 0x6c
	varType       int64 = -21 // 0x6b
	funcType      int64 = -22 // 0x6a
	serviceType   int64 = -23 // 0x69
	principalType int64 = -24 // 0x68
)

type PrimType interface {
	prim()
}

type Type interface {
	// AddTypeDefinition adds itself to the definition table if it is not a primitive type.
	AddTypeDefinition(*TypeDefinitionTable) error

	// Decodes the value from the reader.
	Decode(*bytes.Reader) (interface{}, error)

	// Encodes the type.
	EncodeType(*TypeDefinitionTable) ([]byte, error)

	// Encodes the value.
	EncodeValue(v interface{}) ([]byte, error)

	fmt.Stringer
}

func getType(t int64) (Type, error) {
	// if t >= 0 {
	// 	if int(t) >= len(tds) {
	// 		return nil, fmt.Errorf("type index out of range: %d", t)
	// 	}
	// 	return tds[t], nil
	// }

	switch t {
	case nullType:
		return new(Null), nil
	case boolType:
		return new(Bool), nil
	case natType:
		return new(Nat), nil
	case intType:
		return new(Int), nil
	case natXType:
		return Nat8(), nil
	case natXType - 1:
		return Nat16(), nil
	case natXType - 2:
		return Nat32(), nil
	case natXType - 3:
		return Nat64(), nil
	case intXType:
		return Int8(), nil
	case intXType - 1:
		return Int16(), nil
	case intXType - 2:
		return Int32(), nil
	case intXType - 3:
		return Int64(), nil
	case floatXType:
		return Float32(), nil
	case floatXType - 1:
		return Float64(), nil
	case textType:
		return new(Text), nil
	case reservedType:
		return new(Reserved), nil
	case emptyType:
		return new(Empty), nil
	case principalType:
		return new(Principal), nil
	default:
		if t < -24 {
			return nil, &FormatError{
				Description: "type: out of range",
			}
		}
		return nil, &FormatError{
			Description: "type: not primitive",
		}
	}
}

type primType struct{}

func (primType) AddTypeDefinition(_ *TypeDefinitionTable) error {
	return nil // No need to add primitive types to the type definition table.
}

func (primType) prim() {}
