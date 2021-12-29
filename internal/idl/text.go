package idl

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"unicode/utf8"

	"github.com/aviate-labs/leb128"
)

type Text struct {
	primType
}

func (t *Text) Decode(r *bytes.Reader) (interface{}, error) {
	n, err := leb128.DecodeUnsigned(r)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, n.Int64())
	i, err := r.Read(bs)
	if err != nil {
		return "", nil
	}
	if i != int(n.Int64()) {
		return nil, io.EOF
	}
	if !utf8.Valid(bs) {
		return nil, fmt.Errorf("invalid utf8 text")
	}

	return string(bs), nil
}

func (t Text) EncodeType(_ *TypeDefinitionTable) ([]byte, error) {
	return leb128.EncodeSigned(big.NewInt(textType))
}

func (t Text) EncodeValue(v interface{}) ([]byte, error) {
	v_, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid argument: %v", v)
	}
	bs, err := leb128.EncodeUnsigned(big.NewInt(int64(len(v_))))
	if err != nil {
		return nil, err
	}
	return append(bs, []byte(v_)...), nil
}

func (t Text) String() string {
	return "text"
}
