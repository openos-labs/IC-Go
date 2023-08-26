package idl_test

import (
	"github.com/openos-labs/IC-Go/utils/idl"
	"math/big"
)

func ExampleRec() {
	test([]idl.Type{idl.NewRec(nil)}, []interface{}{nil})
	test_([]idl.Type{idl.NewRec(map[string]idl.Type{
		"foo": new(idl.Text),
		"bar": new(idl.Int),
	})}, []interface{}{
		map[string]interface{}{
			"foo": "💩",
			"bar": big.NewInt(42),
			"baz": big.NewInt(0),
		},
	})
	// Output:
	// 4449444c016c000100
	// 4449444c016c02d3e3aa027c868eb7027101002a04f09f92a9
}
