package agent_test

import (
	"fmt"
	"testing"

	"github.com/aviate-labs/agent-go"
)

func TestNewRequestID(t *testing.T) {
	req := agent.Request{
		Type:       agent.RequestTypeCall,
		CanisterID: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xD2},
		MethodName: "hello",
		Arguments:  []byte("DIDL\x00\xFD*"),
	}
	//t.Log(idl.Decode([]byte("DIDL\x00\xFD*")))
	h := fmt.Sprintf("%x", agent.NewRequestID(req))
	if h != "8781291c347db32a9d8c10eb62b710fce5a93be676474c42babc74c51858f94b" {
		t.Fatal(h)
	}
}

func TestEncodeRequestID(t *testing.T) {
	req := make(map[string]interface{})
	req["request_type"] = agent.RequestTypeCall
	req["canister_id"] = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xD2}
	req["method_name"] = "hello"
	req["arg"] = []byte("DIDL\x00\xFD*")

	id := agent.EncodeRequestID(req)
	h := fmt.Sprintf("%x", id)
	if h != "8781291c347db32a9d8c10eb62b710fce5a93be676474c42babc74c51858f94b" {
		t.Fatal(h)
	}
}
