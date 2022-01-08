package idl_test

import (
	"github.com/mix-labs/IC-Go/utils/idl"
)

func ExampleBool() {
	test([]idl.Type{new(idl.Bool)}, []interface{}{true})
	test([]idl.Type{new(idl.Bool)}, []interface{}{false})
	test([]idl.Type{new(idl.Bool)}, []interface{}{0})
	test([]idl.Type{new(idl.Bool)}, []interface{}{"false"})
	// Output:
	// 4449444c00017e01
	// 4449444c00017e00
	// enc: invalid argument: 0
	// enc: invalid argument: false
}
