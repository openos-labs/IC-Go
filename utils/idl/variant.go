package idl

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/aviate-labs/leb128"
)

type FieldValue struct {
	Name  string
	Value interface{}
}

type Variant struct {
	Fields []Field
}

func NewVariant(fields map[string]Type) *Variant {
	var variant Variant
	for k, v := range fields {
		variant.Fields = append(variant.Fields, Field{
			Name: k,
			Type: v,
		})
	}
	sort.Slice(variant.Fields, func(i, j int) bool {
		return Hash(variant.Fields[i].Name).Cmp(Hash(variant.Fields[j].Name)) < 0
	})
	return &variant
}

func (v Variant) AddTypeDefinition(tdt *TypeDefinitionTable) error {
	for _, f := range v.Fields {
		if err := f.Type.AddTypeDefinition(tdt); err != nil {
			return err
		}
	}

	id, err := leb128.EncodeSigned(big.NewInt(varType))
	if err != nil {
		return err
	}
	l, err := leb128.EncodeUnsigned(big.NewInt(int64(len(v.Fields))))
	if err != nil {
		return err
	}
	var vs []byte
	for _, f := range v.Fields {
		id, err := leb128.EncodeUnsigned(Hash(f.Name))
		if err != nil {
			return nil
		}
		t, err := f.Type.EncodeType(tdt)
		if err != nil {
			return nil
		}
		vs = append(vs, concat(id, t)...)
	}

	tdt.Add(v, concat(id, l, vs))
	return nil
}

func (v Variant) Decode(r *bytes.Reader) (interface{}, error) {
	id, err := leb128.DecodeUnsigned(r)
	if err != nil {
		return nil, err
	}
	if id.Cmp(big.NewInt(int64(len(v.Fields)))) >= 0 {
		return nil, fmt.Errorf("invalid variant index: %v", id)
	}
	v_, err := v.Fields[int(id.Int64())].Type.Decode(r)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{v.Fields[int(id.Int64())].Name:v_}, nil
}

// func (v Variant) Decode(r *bytes.Reader) (interface{}, error) {
// 	id, err := leb128.DecodeUnsigned(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if id.Cmp(big.NewInt(int64(len(v.Fields)))) >= 0 {
// 		return nil, fmt.Errorf("invalid variant index: %v", id)
// 	}
// 	v_, err := v.Fields[int(id.Int64())].Type.Decode(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &FieldValue{
// 		Name:  v.Fields[int(id.Int64())].Name,
// 		Value: v_,
// 	}, nil
// }

func (v Variant) EncodeType(tdt *TypeDefinitionTable) ([]byte, error) {
	idx, ok := tdt.Indexes[v.String()]
	if !ok {
		return nil, fmt.Errorf("missing type index for: %s", v)
	}
	return leb128.EncodeSigned(big.NewInt(int64(idx)))
}

func (v Variant) EncodeValue(value interface{}) ([]byte, error) {
	fs, ok := value.(FieldValue)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	for i, f := range v.Fields {
		if f.Name == fs.Name {
			id, err := leb128.EncodeUnsigned(big.NewInt(int64(i)))
			if err != nil {
				return nil, err
			}
			v_, err := f.Type.EncodeValue(fs.Value)
			if err != nil {
				return nil, err
			}
			return concat(id, v_), nil
		}
	}
	return nil, fmt.Errorf("unknown variant: %v", value)
}

func (v Variant) String() string {
	var s []string
	for _, f := range v.Fields {
		s = append(s, fmt.Sprintf("%s:%s", f.Name, f.Type.String()))
	}
	return fmt.Sprintf("variant {%s}", strings.Join(s, "; "))
}
