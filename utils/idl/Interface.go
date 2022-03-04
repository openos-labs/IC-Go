package idl

import (
	"bytes"
	"fmt"
)

type Interface struct {
	_Type Type
	primType
}

//func NewInterface(t Type) *Interface {
//	return &Interface{
//		Type: t,
//	}
//}

func (i *Interface) Fill(t Type) {
	i._Type = t
}

func (i *Interface) Decode(r *bytes.Reader) (interface{}, error) {
	return i._Type.Decode(r)
}
func (i Interface) EncodeType(tdt *TypeDefinitionTable) ([]byte, error) {
	return i._Type.EncodeType(tdt)
}
func (i Interface) EncodeValue(value interface{}) ([]byte, error) {
	return i._Type.EncodeValue(value)
}
func (i Interface) String() string {
	return fmt.Sprintf("interface %s", i._Type)
}
