package idl_test

import "github.com/mix-labs/IC-Go/utils/idl"

func ExampleVariant() {
	result := map[string]idl.Type{
		"ok":  new(idl.Text),
		"err": new(idl.Text),
	}
	test_([]idl.Type{idl.NewVariant(result)}, []interface{}{idl.FieldValue{
		Name:  "ok",
		Value: "good",
	}})
	test_([]idl.Type{idl.NewVariant(result)}, []interface{}{idl.FieldValue{
		Name:  "err",
		Value: "uhoh",
	}})
	// Output:
	// 4449444c016b029cc20171e58eb4027101000004676f6f64
	// 4449444c016b029cc20171e58eb402710100010475686f68
}
