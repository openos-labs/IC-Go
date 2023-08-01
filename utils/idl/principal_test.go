package idl_test

import (
	"github.com/openos-labs/IC-Go/utils/idl"
	"github.com/openos-labs/IC-Go/utils/principal"
)

func ExamplePrincipal() {
	p, _ := principal.Decode("gvbup-jyaaa-aaaah-qcdwa-cai")
	test([]idl.Type{new(idl.Principal)}, []interface{}{p})
	// Output:
	// 4449444c000168010a0000000000f010ec0101
}
