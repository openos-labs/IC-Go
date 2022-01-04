package agent

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stopWarByWar/ic-agent/internal/identity"
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"github.com/stopWarByWar/ic-agent/internal/principal"
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
	t.Log("errMsg:", errMsg, "err:", err, "result:", result[0].(*big.Int))
}

func TestAgent_UpdateRaw(t *testing.T) {
	canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	
	//fmt.Println(hex.EncodeToString(agent.identity.PubKeyBytes()))
	//t.Log("time:", uint64(agent.getExpiryDate().UnixNano()))
	//envelope := new(Envelope)
	//data, _ := hex.DecodeString("a167636f6e74656e74a66c726571756573745f747970656463616c6c6673656e64657241046e696e67726573735f6578706972791a2f449dd06b63616e69737465725f69644a0000000000f010ec01016b6d6574686f645f6e616d65687472616e73666572636172674f4449444c0002687d010080c8afa025")
	//cbor.Unmarshal(data, envelope)
	//t.Log("envelope ", envelope.Content)
	//
	//t.Log("sender", envelope.Content.Sender.Encode())
	//t.Log("type", envelope.Content.Type)
	//t.Log("ingress expiry", envelope.Content.IngressExpiry)
	//t.Log("method", envelope.Content.MethodName)
	//t.Log("arg", envelope.Content.Arguments)
	//t.Log("canister", envelope.Content.CanisterID.Encode())

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

func TestCbor(t *testing.T) {
	canisterID, _ := principal.Decode("gvbup-jyaaa-aaaah-qcdwa-cai")
	agent := New(true, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")

	//t.Log("time:", uint64(agent.getExpiryDate().UnixNano()))
	req := Request{
		Type:          "call",
		Sender:        agent.Sender(),
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
		CanisterID:    canisterID,
		MethodName:    "transfer",
		Arguments:     []byte("i love vivian"),
	}

	envelope := Envelope{
		req,
		[]byte{},
		[]byte{},
	}

	data, _ := cbor.Marshal(envelope)
	resp := new(Envelope)
	cbor.Unmarshal(data, resp)
	t.Log("sender", resp.Content.Sender.Encode())
	t.Log("type", resp.Content.Type)
	t.Log("ingress expiry", resp.Content.IngressExpiry)
	t.Log("method", resp.Content.MethodName)
	t.Log("arg", resp.Content.Arguments)
	t.Log("canister", resp.Content.CanisterID.Encode())
}