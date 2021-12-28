package agent

import (
	"encoding/hex"
	identity "github.com/aviate-labs/agent-go/internal/key"
	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/principal-go"
	"testing"
)

func TestAgent_QueryRaw(t *testing.T) {
	canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
	agent := New(true, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	methodName := "totalSupply"
	arg, err := idl.Encode([]idl.Type{new(idl.Rec)}, []interface{}{nil})
	if err != nil {
		t.Error(err)
	}
	_, result, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)
	t.Log("errMsg:", errMsg, "err:", err, "result:", result)
	//t.Log("errMsg:", errMsg, "err:", err, "result:", result[0].(big.Int))

}

func TestPrincipal(t *testing.T) {

	pkBytes, _ := hex.DecodeString("833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	identity := identity.New(false, pkBytes)

	p := principal.NewSelfAuthenticating(identity.PubKey.SerializeUncompressed())
	t.Log(p.Encode(), len(identity.PubKey.SerializeUncompressed()))
}
