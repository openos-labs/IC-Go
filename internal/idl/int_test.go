package idl_test

import (
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"math/big"
)

func ExampleInt() {
	test([]idl.Type{new(idl.Int)}, []interface{}{big.NewInt(0)})
	test([]idl.Type{new(idl.Int)}, []interface{}{big.NewInt(42)})
	test([]idl.Type{new(idl.Int)}, []interface{}{big.NewInt(1234567890)})
	test([]idl.Type{new(idl.Int)}, []interface{}{big.NewInt(-1234567890)})
	test([]idl.Type{new(idl.Int)}, []interface{}{func() *big.Int {
		bi, _ := new(big.Int).SetString("60000000000000000", 10)
		return bi
	}()})
	// Output:
	// 4449444c00017c00
	// 4449444c00017c2a
	// 4449444c00017cd285d8cc04
	// 4449444c00017caefaa7b37b
	// 4449444c00017c808098f4e9b5caea00
}

func ExampleInt32() {
	test([]idl.Type{idl.Int32()}, []interface{}{big.NewInt(-1234567890)})
	test([]idl.Type{idl.Int32()}, []interface{}{big.NewInt(-42)})
	test([]idl.Type{idl.Int32()}, []interface{}{big.NewInt(42)})
	test([]idl.Type{idl.Int32()}, []interface{}{big.NewInt(1234567890)})
	// Output:
	// 4449444c0001752efd69b6
	// 4449444c000175d6ffffff
	// 4449444c0001752a000000
	// 4449444c000175d2029649
}

func ExampleInt8() {
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(-129)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(-128)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(-42)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(-1)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(0)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(1)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(42)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(127)})
	test([]idl.Type{idl.Int8()}, []interface{}{big.NewInt(128)})
	// Output:
	// enc: invalid value: -129
	// 4449444c00017780
	// 4449444c000177d6
	// 4449444c000177ff
	// 4449444c00017700
	// 4449444c00017701
	// 4449444c0001772a
	// 4449444c0001777f
	// enc: invalid value: 128
}
