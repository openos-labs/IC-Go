package agent

import (
	"github.com/aviate-labs/candid-go/idl"
	"math/big"
	"testing"
)

func TestAgent(t *testing.T) {
	canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
	agent, _ := New()
	methodName := "totalSupply"
	arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	if err != nil {
		t.Error(err)
	}
	_, values, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)
	if err != nil {
		t.Errorf("can not get rusult with errMsg: %v,err: %v", errMsg, err)
	} else {
		t.Errorf("the result is %v", values[0].(big.Int))
	}
}
