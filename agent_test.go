package agent

import (
	"bytes"
	"encoding/hex"
	"github.com/aviate-labs/agent-go/internal/key"
	"github.com/aviate-labs/bip39"
	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/principal-go"
	"github.com/fxamacker/cbor/v2"
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

	_, _, errMsg, err := agent.QueryRaw(canisterID, methodName, arg)
	t.Log("errMsg:", errMsg, "err:", err)

}

type Envelope1 struct {
	Content      *Request `cbor:"content,omitempty"`
	SenderPubkey []byte   `cbor:"sender_pubkey,omitempty"`
	SenderSig    []byte   `cbor:"sender_sig,omitempty"`
}

type Envelope2 struct {
	Content      map[string]interface{} `cbor:"content,omitempty"`
	SenderPubkey []byte                 `cbor:"sender_pubkey,omitempty"`
	SenderSig    []byte                 `cbor:"sender_sig,omitempty"`
}

func TestNew(t *testing.T) {
	a := make(map[string]interface{})
	b := make(map[string]string)
	a["a"] = "a"
	a["b"] = "b"
	b["a"] = "a"
	b["b"] = "b"

	a1, _ := cbor.Marshal(a)
	b1, _ := cbor.Marshal(b)
	t.Log("a:", a1, "\nb:", b1, "\n", bytes.Equal(a1, b1))
}

func TestNew2(t *testing.T) {
	canisterID, _ := principal.Decode("gvbup-jyaaa-aaaah-qcdwa-cai")
	agent, _ := New()
	methodName := "totalSupply"
	arg, err := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
	if err != nil {
		t.Error(err)
	}

	req := Request{
		Type:          RequestTypeQuery,
		Sender:        *agent.Sender(),
		CanisterID:    canisterID,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}

	request := make(map[string]interface{})
	request["type"] = req.Type
	request["sender"] = req.Sender
	request["canister_id"] = req.CanisterID
	request["method_name"] = req.MethodName
	request["ingress_expiry"] = req.IngressExpiry
	request["arg"] = req.Arguments

	requestID := NewRequestID(req)
	msg := []byte(IC_REQUEST_DOMAIN_SEPARATOR)
	msg = append(msg, requestID[:]...)
	//sig, _ := agent.key.Sign(msg)

	//envelope1 := Envelope1{
	//	Content:      &req,
	//	SenderPubkey: agent.key.PubKey.SerializeUncompressed(),
	//	SenderSig:    sig.Serialize(),
	//}
	//envelope2 := Envelope2{
	//	Content:      request,
	//	SenderPubkey: agent.key.PubKey.SerializeUncompressed(),
	//	SenderSig:    sig.Serialize(),
	//}
	mashaledEnvelope1, _ := cbor.Marshal(request)
	mashaledEnvelope2, _ := cbor.Marshal(req)
	t.Log("a: ", hex.EncodeToString(mashaledEnvelope1), "\nb:", hex.EncodeToString(mashaledEnvelope2), "\n", bytes.Equal(mashaledEnvelope2, mashaledEnvelope1))
}

func TestPrincipal(t *testing.T) {
	e, _ := bip39.NewEntropy(128)
	m, _ := bip39.English.NewMnemonic(e)
	n, _ := key.New(m, "")
	_, pub, _ := key.Keys(n)
	//p := principal.NewSelfAuthenticating(pub.SerializeUncompressed())
	//pubByte,_ := hex.DecodeString("3056301006072a8648ce3d020106052b8104000a0342000480ef3bac9d68cf374cbc9c9943e180043a94c462ef8270274e57089d5dcdba1b8fbf83b7546ccc1b3781377776e9c4710fdf533ed9bcba8d8ebb32a14a2aba11")
	pubByte := pub.SerializeCompressed()
	pubByte1 := pub.SerializeUncompressed()
	p := principal.NewSelfAuthenticating(pubByte)
	p1 := principal.NewSelfAuthenticating(pubByte1)
	t.Log(p.Encode(),p1.Encode(),len(pubByte))
}
