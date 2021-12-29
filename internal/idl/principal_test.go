package idl_test

import (
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"github.com/stopWarByWar/ic-agent/internal/principal"
)

func ExamplePrincipal() {
	p, _ := principal.Decode("gvbup-jyaaa-aaaah-qcdwa-cai")
	test([]idl.Type{new(idl.Principal)}, []interface{}{p})
	// Output:
	// 4449444c000168010a0000000000f010ec0101
}
