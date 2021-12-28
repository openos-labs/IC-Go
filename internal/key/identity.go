package identity

import (
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcec"
)

func New(anonymous bool, pkBytes []byte) *Identity {
	if anonymous == true {
		return &Identity{
			Anonymous: anonymous,
		}
	}
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return &Identity{
		anonymous,
		privKey,
		pubkey,
	}
}

type Identity struct {
	Anonymous bool
	PriKey    *btcec.PrivateKey
	PubKey    *btcec.PublicKey
}

func (identity *Identity) Sign(m []byte) ([]byte, error) {
	if identity.Anonymous == true {
		return []byte{}, nil
	}
	hashByte := sha256.Sum256(m)
	sign, err := identity.PriKey.Sign(hashByte[:])
	if err != nil {
		return nil, err
	}
	return sign.Serialize(), nil
}

func (identity *Identity) PubKeyBytes() []byte {
	var senderPubKey []byte
	if identity.Anonymous == false {
		return identity.PubKey.SerializeUncompressed()
	}
	return senderPubKey

}
