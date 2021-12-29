package idl_test

import (
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"math/big"
)

func ExampleVec() {
	test([]idl.Type{idl.NewVec(new(idl.Int))}, []interface{}{
		[]interface{}{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)},
	})
	// Output:
	// 4449444c016d7c01000400010203
}
