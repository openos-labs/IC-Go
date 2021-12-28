package identity

import (
	"encoding/hex"
	"testing"
)

func TestKeys(t *testing.T) {

	pkBytes, _ := hex.DecodeString("833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
	identity := New(false, pkBytes)
	sign,_ := identity.Sign([]byte("vvvvvvvvvvvvvvvvvaaaaaaaaaaaaappppppppppppppyyyyyyyyyyddddddddd"))
	t.Log(hex.EncodeToString(sign))
}
