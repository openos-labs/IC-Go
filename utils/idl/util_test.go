package idl_test

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/openos-labs/IC-Go/utils/idl"
)

func TestHash(t *testing.T) {
	if h := idl.Hash("foo"); h.Cmp(big.NewInt(5097222)) != 0 {
		t.Errorf("expected '5097222', got '%s'", h)
	}
	if h := idl.Hash("bar"); h.Cmp(big.NewInt(4895187)) != 0 {
		t.Errorf("expected '4895187', got '%s'", h)
	}
}

func test(types []idl.Type, args []interface{}) {
	e, err := idl.Encode(types, args)
	if err != nil {
		fmt.Println("enc:", err)
		return
	}
	fmt.Printf("%x\n", e)

	ts, vs, err := idl.Decode(e)
	if err != nil {
		fmt.Println("dec:", err)
		return
	}
	if !reflect.DeepEqual(ts, types) {
		fmt.Println("types:", types, ts)
	}
	if !reflect.DeepEqual(vs, args) {
		fmt.Println("args:", args, vs)
	}
}

func test_(types []idl.Type, args []interface{}) {
	e, err := idl.Encode(types, args)
	if err != nil {
		fmt.Println("enc:", err)
		return
	}
	fmt.Printf("%x\n", e)

	if _, _, err := idl.Decode(e); err != nil {
		fmt.Println("dec:", err)
		return
	}
}
