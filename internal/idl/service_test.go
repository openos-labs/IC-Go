package idl_test

import (
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"github.com/stopWarByWar/ic-agent/internal/principal"
)

func ExampleService() {
	p, _ := principal.Decode("w7x7r-cok77-xa")
	test(
		[]idl.Type{idl.NewService(
			map[string]*idl.Func{
				"foo": idl.NewFunc(
					[]idl.Type{new(idl.Text)},
					[]idl.Type{new(idl.Nat)},
					nil,
				),
			},
		)},
		[]interface{}{
			p,
		},
	)
	// Output:
	// 4449444c026a0171017d00690103666f6f0001010103caffee
}
