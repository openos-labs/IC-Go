package idl

import (
	"bytes"
	"fmt"

	"github.com/aviate-labs/leb128"
	"github.com/openos-labs/IC-Go/utils/principal"

	"math/big"
	"sort"
	"strings"
)

type Method struct {
	Name string
	Func *Func
}

type Service struct {
	methods []Method
}

func NewService(methods map[string]*Func) *Service {
	var service Service
	for k, v := range methods {
		service.methods = append(service.methods, Method{
			Name: k,
			Func: v,
		})
	}
	sort.Slice(service.methods, func(i, j int) bool {
		return Hash(service.methods[i].Name).Cmp(Hash(service.methods[j].Name)) < 0
	})
	return &service
}

func (s Service) AddTypeDefinition(tdt *TypeDefinitionTable) error {
	for _, f := range s.methods {
		if err := f.Func.AddTypeDefinition(tdt); err != nil {
			return err
		}
	}

	id, err := leb128.EncodeSigned(big.NewInt(serviceType))
	if err != nil {
		return err
	}
	l, err := leb128.EncodeUnsigned(big.NewInt(int64(len(s.methods))))
	if err != nil {
		return err
	}
	var vs []byte
	for _, f := range s.methods {
		id := []byte(f.Name)
		l, err := leb128.EncodeUnsigned(big.NewInt(int64(len((id)))))
		if err != nil {
			return nil
		}
		t, err := f.Func.EncodeType(tdt)
		if err != nil {
			return nil
		}
		vs = append(vs, concat(l, id, t)...)
	}

	tdt.Add(s, concat(id, l, vs))
	return nil
}

func (s Service) Decode(r *bytes.Reader) (interface{}, error) {
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
	pid := make(principal.Principal, l.Int64())
	n, err := r.Read(pid)
	if err != nil {
		return nil, err
	}
	if n != int(l.Int64()) {
		return nil, fmt.Errorf("invalid principal id: %d", pid)
	}
	return pid, nil
}

func (s Service) EncodeType(tdt *TypeDefinitionTable) ([]byte, error) {
	idx, ok := tdt.Indexes[s.String()]
	if !ok {
		return nil, fmt.Errorf("missing type index for: %s", s)
	}
	return leb128.EncodeSigned(big.NewInt(int64(idx)))
}

func (s Service) EncodeValue(v interface{}) ([]byte, error) {
	p, ok := v.(principal.Principal)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	l, err := leb128.EncodeUnsigned(big.NewInt(int64(len(p))))
	if err != nil {
		return nil, err
	}
	return concat([]byte{0x01}, l, []byte(p)), nil
}

func (Service) Fill(Type) {

}
func (s Service) String() string {
	var methods []string
	for _, m := range s.methods {
		methods = append(methods, fmt.Sprintf("%s:%s", m.Name, m.Func.String()))
	}
	return fmt.Sprintf("service {%s}", strings.Join(methods, "; "))
}
