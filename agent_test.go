package agent

import (
	"encoding/hex"
	"github.com/stopWarByWar/ic-agent/internal/identity"
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"github.com/stopWarByWar/ic-agent/internal/principal"
	"math/big"
	"testing"
)

func TestAgent_QueryRaw(t *testing.T) {
	canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	methodName := "totalSupply"
	arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	if err != nil {
		t.Error(err)
	}
	_, result, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)
	t.Log("errMsg:", errMsg, "err:", err, "result:", result)
	//t.Log("errMsg:", errMsg, "err:", err, "result:", result[0].(big.Int))
}

func TestAgent_UpdateRaw(t *testing.T) {
	canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	methodName := "transfer"
	var argType []idl.Type
	var argValue []interface{}

	p, _ := principal.Decode("aaaaa-aa")
	argType = append(argType, new(idl.Principal))
	argType = append(argType, new(idl.Nat))
	argValue = append(argValue, p)
	argValue = append(argValue, big.NewInt(10000000000))

	arg, _ := idl.Encode(argType, argValue)
	_, result, err := agent.UpdateRaw(canisterID, methodName, arg)
	t.Log("errMsg:", err, "result:", result)
}

func TestPrincipal(t *testing.T) {
	pkBytes, _ := hex.DecodeString("833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	identity := identity.New(false, pkBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log(p.Encode(), len(identity.PubKeyBytes()))
}
