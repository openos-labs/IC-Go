package identity

import (
	"encoding/hex"
	"strconv"
	"testing"
)

func TestKeys(t *testing.T) {
	pkBytes, _ := hex.DecodeString("833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	identity := New(false, pkBytes)
	t.Log(hex.EncodeToString([]byte("pyd")))
	sign,_ := identity.Sign([]byte("pyd"))
	t.Log(hex.EncodeToString(sign),len(sign))
}

func TestNew(t *testing.T) {
	var a []byte
	result := strconv.AppendInt(a,-24,8)
	t.Logf("%x",result)
}