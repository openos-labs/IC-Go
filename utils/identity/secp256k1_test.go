package identity

import (
	"bytes"
	"testing"
)

func TestNewSecp256k1Identity(t *testing.T) {
	id, _ := NewRandomSecp256k1Identity()
	data, err := id.ToPEM()
	if err != nil {
		t.Fatal(err)
	}
	id_, err := NewSecp256k1IdentityFromPEM(data)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(id.privateKey.Serialize(), id_.privateKey.Serialize()) {
		t.Error()
	}
	if !bytes.Equal(id.PublicKey(), id_.PublicKey()) {
		t.Error()
	}
}

func TestName(t *testing.T) {
	pri := "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42"
	id, err := NewSecp256k1IdentityFromHex(pri)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(id.Sender().Encode())
}
