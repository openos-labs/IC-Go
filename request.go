package agent

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"sort"

	"github.com/aviate-labs/leb128"
	"github.com/aviate-labs/principal-go"
)

var (
	typeKey          = sha256.Sum256([]byte("request_type"))
	canisterIDKey    = sha256.Sum256([]byte("canister_id"))
	nonceKey         = sha256.Sum256([]byte("nonce"))
	methodNameKey    = sha256.Sum256([]byte("method_name"))
	argumentsKey     = sha256.Sum256([]byte("arg"))
	ingressExpiryKey = sha256.Sum256([]byte("ingress_expiry"))
	senderKey        = sha256.Sum256([]byte("sender"))
)

type RequestType = string

const (
	RequestTypeCall      RequestType = "call"
	RequestTypeQuery     RequestType = "query"
	RequestTypeReadState RequestType = "read_state"
)

// DOCS: https://smartcontracts.org/docs/interface-spec/index.html#http-call
type Request struct {
	Type RequestType `cbor:"request_type,omitempty"`
	// The user who issued the request.
	Sender principal.Principal `cbor:"sender,omitempty"`
	// Arbitrary user-provided data, typically randomly generated. This can be
	// used to create distinct requests with otherwise identical fields.
	Nonce []byte `cbor:"nonce,omitempty"`
	// An upper limit on the validity of the request, expressed in nanoseconds
	// since 1970-01-01 (like ic0.time()).
	IngressExpiry uint64 `cbor:"ingress_expiry,omitempty"`
	// The principal of the canister to call.
	CanisterID principal.Principal `cbor:"canister_id"`
	// Name of the canister method to call.
	MethodName string `cbor:"method_name,omitempty"`
	// Argument to pass to the canister method.
	Arguments []byte `cbor:"arg,omitempty"`
	// Paths (sequence of paths): A list of paths, where a path is itself a sequence of blobs.
	Paths [][]byte `cbor:"paths,omitempty"`
}

type RequestID [32]byte

func EncodeRequestID(request map[string]interface{}) RequestID {
	var hashes [][]byte
	for k, v := range request {
		switch k {
		case "request_type":
			typeKey := sha256.Sum256([]byte("request_type"))
			typeValue := sha256.Sum256([]byte(v.(string)))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "sender":
			typeKey := sha256.Sum256([]byte("sender"))
			typeValue := sha256.Sum256(v.(principal.Principal))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "nonce":
			typeKey := sha256.Sum256([]byte("nonce"))
			typeValue := sha256.Sum256(v.([]byte))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "ingress_expiry":
			typeKey := sha256.Sum256([]byte("ingress_expiry"))
			typeValue := sha256.Sum256(encodeLEB128(v.(uint64)))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "canister_id":
			typeKey := sha256.Sum256([]byte("canister_id"))
			typeValue := sha256.Sum256(v.(principal.Principal))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "method_name":
			typeKey := sha256.Sum256([]byte("method_name"))
			typeValue := sha256.Sum256([]byte(v.(string)))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		case "arg":
			typeKey := sha256.Sum256([]byte("arg"))
			typeValue := sha256.Sum256(v.([]byte))
			result := append(typeKey[:],typeValue[:]...)
			hashes = append(hashes,result)
		default:
		}
	}
	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i], hashes[j]) == -1
	})
	return sha256.Sum256(bytes.Join(hashes, nil))
}

// DOCS: https://smartcontracts.org/docs/interface-spec/index.html#request-id
func NewRequestID(req Request) RequestID {
	var (
		typeHash       = sha256.Sum256([]byte(req.Type))
		canisterIDHash = sha256.Sum256(req.CanisterID)
		methodNameHash = sha256.Sum256([]byte(req.MethodName))
		argumentsHash  = sha256.Sum256(req.Arguments)
	)
	hashes := [][]byte{
		append(typeKey[:], typeHash[:]...),
		append(canisterIDKey[:], canisterIDHash[:]...),
		append(methodNameKey[:], methodNameHash[:]...),
		append(argumentsKey[:], argumentsHash[:]...),
	}
	if len(req.Sender) != 0 {
		senderHash := sha256.Sum256(req.Sender)
		hashes = append(hashes, append(senderKey[:], senderHash[:]...))
	}
	if req.IngressExpiry != 0 {
		ingressExpiryHash := sha256.Sum256(encodeLEB128(req.IngressExpiry))
		hashes = append(hashes, append(ingressExpiryKey[:], ingressExpiryHash[:]...))
	}
	if len(req.Nonce) != 0 {
		nonceHash := sha256.Sum256(req.Nonce)
		hashes = append(hashes, append(nonceKey[:], nonceHash[:]...))
	}
	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i], hashes[j]) == -1
	})
	return sha256.Sum256(bytes.Join(hashes, nil))
}

func encodeLEB128(i uint64) []byte {
	bi := big.NewInt(int64(i))
	e, _ := leb128.EncodeUnsigned(bi)
	return e
}

type Envelope struct {
	Content      *Request `cbor:"content,omitempty"`
	SenderPubkey []byte   `cbor:"sender_pubkey,omitempty"`
	SenderSig    []byte   `cbor:"sender_sig,omitempty"`
}
