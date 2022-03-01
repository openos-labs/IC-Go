package agent_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	agent "github.com/mix-labs/IC-Go"
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

func TestEncodeList(t *testing.T) {
	a := [][]byte{[]byte("i"),[]byte("love"),[]byte("you")}
	res := encodeList(a)
	t.Log(hex.EncodeToString([]byte("i")))
	t.Log(hex.EncodeToString([]byte("love")))
	t.Log(hex.EncodeToString([]byte("you")))
	t.Log(hex.EncodeToString(res[:]))
}

func encodeList(paths [][]byte) [32]byte {
	var pathsBytes []byte
	for _, path := range paths {
		pathBytes := sha256.Sum256(path)
		pathsBytes = append(pathsBytes, pathBytes[:]...)
	}
	return sha256.Sum256(pathsBytes)
}

