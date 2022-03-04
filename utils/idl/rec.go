package idl

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/aviate-labs/leb128"
)

type Rec struct {
	Fields []Field
}

func NewRec(fields map[string]Type) *Rec {
	var rec Rec
	for k, v := range fields {
		rec.Fields = append(rec.Fields, Field{
			Name: k,
			Type: v,
		})
	}
	sort.Slice(rec.Fields, func(i, j int) bool {
		return Hash(rec.Fields[i].Name).Cmp(Hash(rec.Fields[j].Name)) < 0
	})
	return &rec
}

func (r Rec) AddTypeDefinition(tdt *TypeDefinitionTable) error {
	for _, f := range r.Fields {
		if err := f.Type.AddTypeDefinition(tdt); err != nil {
			return err
		}
	}

	id, err := leb128.EncodeSigned(big.NewInt(recType))
	if err != nil {
		return err
	}
	l, err := leb128.EncodeUnsigned(big.NewInt(int64(len(r.Fields))))
	if err != nil {
		return err
	}
	var vs []byte
	for _, f := range r.Fields {
		l, err := leb128.EncodeUnsigned(Hash(f.Name))
		if err != nil {
			return nil
		}
		t, err := f.Type.EncodeType(tdt)
		if err != nil {
			return nil
		}
		vs = append(vs, concat(l, t)...)
	}

	tdt.Add(r, concat(id, l, vs))
	return nil
}

func (r Rec) Decode(r_ *bytes.Reader) (interface{}, error) {
	rec := make(map[string]interface{})
	for _, f := range r.Fields {
		v, err := f.Type.Decode(r_)
		if err != nil {
			return nil, err
		}
		rec[f.Name] = v
	}
	if len(rec) == 0 {
		return nil, nil
	}
	return rec, nil
}

func (r Rec) EncodeType(tdt *TypeDefinitionTable) ([]byte, error) {
	idx, ok := tdt.Indexes[r.String()]
	if !ok {
		return nil, fmt.Errorf("missing type index for: %s", r)
	}
	return leb128.EncodeSigned(big.NewInt(int64(idx)))
}

func (r Rec) EncodeValue(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	fs, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	var vs_ []interface{}
	for _, f := range r.Fields {
		vs_ = append(vs_, fs[f.Name])
	}
	var vs []byte
	for i, f := range r.Fields {
		v_, err := f.Type.EncodeValue(vs_[i])
		if err != nil {
			return nil, err
		}
		vs = append(vs, v_...)
	}
	return vs, nil
}

func (Rec) Fill(Type) {

}

func (r Rec) String() string {
	var s []string
	for _, f := range r.Fields {
		s = append(s, fmt.Sprintf("%s:%s", f.Name, f.Type.String()))
	}
	return fmt.Sprintf("record {%s}", strings.Join(s, "; "))
}
