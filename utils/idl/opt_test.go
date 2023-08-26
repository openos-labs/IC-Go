package idl_test

import (
	"github.com/openos-labs/IC-Go/utils/idl"
	"math/big"
)

func ExampleOpt() {
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{nil})
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{big.NewInt(1)})
	// Output:
	// 4449444c016e7d010000
	// 4449444c016e7d01000101
}
