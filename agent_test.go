package agent

import (
	"encoding/hex"
	"fmt"
	"github.com/mix-labs/IC-Go/utils/identity"
	"math/big"
	"reflect"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/mix-labs/IC-Go/utils/idl"
	"github.com/mix-labs/IC-Go/utils/principal"
)

func TestAgent_QueryRaw(t *testing.T) {
	canisterID := "bzsui-sqaaa-aaaah-qce2a-cai"
	//agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	agent,err := NewFromPem(false,"/Users/panyeda/go/src/github.com/mix-labs/IC-Go/utils/identity/priv.pem")
	if err != nil{
		t.Log(err)
	}
	methodName := "supply"
	//methodName := "listings"
	//arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	arg, err := idl.Encode([]idl.Type{new(idl.Text)}, []interface{}{"Motoko"})
	if err != nil {
		t.Error(err)
	}
	_, result, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)

	myresult := map[string]interface{}{}
	//myresult["err"] = map[string]interface{}{"InvalidToken":"", "Other":""}
	myresult["ok"] = 0

	//fmt.Println(reflect.ValueOf(result[0]))

	_result, ok := result[0].(*idl.FieldValue)
	if !ok {
		return
	}


	for key, _ := range(myresult) {
		if idl.Hash(key).String() == _result.Name {
			fmt.Println(reflect.TypeOf(_result.Value))
			final_value, ok := _result.Value.(*big.Int)
			if !ok {
				fmt.Println("error type")
				return
			}
			myresult[key] = final_value.Int64()
		}

	}
	t.Log("errMsg:", errMsg, "err:", err, "result:", myresult)
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
	t.Log("errMsg:", err, "result:", result[0])
}

func TestAgent_GetCanisterController(t *testing.T) {
	canisterID := "tviwd-niaaa-aaaaf-qaaaa-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	controllersBytes, err := agent.GetCanisterControllers(canisterID)
	if err != nil{
		t.Log(err)
	}
	for _, controllerBytes := range controllersBytes{
		controller := principal.New(controllerBytes).Encode()
		t.Log("controller:",controller)
	}
}

func TestAgent_GetCanisterInfo_ModuleHash(t *testing.T) {
	canisterID := "krnha-raaaa-aaaaf-qac5q-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	moduleHash, err := agent.GetCanisterModuleHash(canisterID)
	if err != nil{
		t.Log(err)
	}
	t.Log("module_hash:",hex.EncodeToString(moduleHash))
}

func TestAgent_GetCanisterSubnet(t *testing.T) {
	canisterID := "krnha-raaaa-aaaaf-qac5q-cai"
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	subnet, err := agent.GetCanisterSubnet(canisterID)
	if err != nil{
		t.Log(err)
	}
	t.Log("subnet:",hex.EncodeToString(subnet))
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
