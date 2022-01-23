package agent

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/mix-labs/IC-Go/utils/identity"
	"github.com/mix-labs/IC-Go/utils/idl"
	"github.com/mix-labs/IC-Go/utils/principal"
)




type supply struct {
	Ok uint64 `ic:"ok"`
	Err string	`ic:"err"`
}
type Time struct {
	Some big.Int	`ic:"some"`
	None uint8	`ic:"none"`
}
type listing struct {
	Locked Time `ic:locked`
	Price uint64 `ic:"price"`
	Seller principal.Principal `ic:"seller"`
}
type listingTuple struct {
	A uint32 `ic:"0"`
	B listing `ic:"1"`
}
type listings []listingTuple

type TokenIndex uint32
type RegistryTuple struct {
	A TokenIndex `ic:"0"`
	B string `ic:"1"`
}
type Registrys []RegistryTuple
func TestAgent_QueryRaw(t *testing.T) {
	canisterID := "bzsui-sqaaa-aaaah-qce2a-cai"
	
	//agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	agent,err := NewFromPem(false,"./utils/identity/priv.pem")
	if err != nil{
		t.Log(err)
	}
	//methodName := "supply"
	methodName := "listings"
	//methodName := "getRegistry"


	//arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	arg, err := idl.Encode([]idl.Type{new(idl.Text)}, []interface{}{"Motoko"})
	if err != nil {
		t.Error(err)
	}
	Type, result, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)
	//myresult := Registrys{}
	myresult := listings{}
	//myresult := supply{}
	Decode(&myresult, result[0])
	t.Log("errMsg:", errMsg, "err:", err, "result:", myresult, "type:", Type)
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