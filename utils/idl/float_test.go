package idl_test

import (
	"math/big"

	"github.com/mix-labs/IC-Go/utils/idl"
)

func ExampleFloat32() {
	test([]idl.Type{idl.Float32()}, []interface{}{big.NewFloat(-0.5)})
	test([]idl.Type{idl.Float32()}, []interface{}{big.NewFloat(0)})
	test([]idl.Type{idl.Float32()}, []interface{}{big.NewFloat(0.5)})
	test([]idl.Type{idl.Float32()}, []interface{}{big.NewFloat(3)})
	// Output:
	// 4449444c000173000000bf
	// 4449444c00017300000000
	// 4449444c0001730000003f
	// 4449444c00017300004040
}

func ExampleFloat64() {
	test([]idl.Type{idl.Float64()}, []interface{}{big.NewFloat(-0.5)})
	test([]idl.Type{idl.Float64()}, []interface{}{big.NewFloat(0)})
	test([]idl.Type{idl.Float64()}, []interface{}{big.NewFloat(0.5)})
	test([]idl.Type{idl.Float64()}, []interface{}{big.NewFloat(3)})
	// Output:
	// 4449444c000172000000000000e0bf
	// 4449444c0001720000000000000000
	// 4449444c000172000000000000e03f
	// 4449444c0001720000000000000840
}
