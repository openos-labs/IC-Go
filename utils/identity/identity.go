package identity

import "github.com/openos-labs/IC-Go/utils/principal"

// Identity is an identity that can sign messages.
type Identity interface {
	// Sender returns the principal of the identity.
	Sender() principal.Principal
	// Sign signs the given message.
	Sign(msg []byte) []byte
	// PublicKey returns the public key of the identity.
	PublicKey() []byte
}
