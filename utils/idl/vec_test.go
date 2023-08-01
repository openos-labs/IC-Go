package idl_test

import (
	"math/big"

	"github.com/openos-labs/IC-Go/utils/idl"
)

func ExampleVec() {
	test([]idl.Type{idl.NewVec(new(idl.Int))}, []interface{}{
		[]interface{}{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)},
	})
	// Output:
	// 4449444c016d7c01000400010203
}
