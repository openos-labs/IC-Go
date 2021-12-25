package agent

import (
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/internal/key"
	"github.com/aviate-labs/bip39"
	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/principal-go"
	"github.com/fxamacker/cbor/v2"
	"golang.org/x/crypto/sha3"
	"net/url"
	"time"
)

type Agent struct {
	client        *Client
	key           *key.Pairs
	ingressExpiry time.Duration
	rootKey       []byte		//ICP root key
}

func New() (*Agent,bip39.Mnemonic) {
	var ic0, _ = url.Parse("https://ic0.app/")
	c := NewClient(ClientConfig{ic0})
	//todo:是否需要从ic拉取rootKey信息
	status, _ := c.Status()
	e, _ := bip39.NewEntropy(128)
	m, _ := bip39.English.NewMnemonic(e)
	n, _ := key.New(m, "")
	priv, pubKey, _ := key.Keys(n)
	pair := &key.Pairs{
		PriKey: priv,
		PubKey: pubKey,
	}

	ingressExpiry := time.Second * 10
	return &Agent{
		client:        &c,
		key:           pair,
		ingressExpiry: ingressExpiry,
		rootKey:       status.RootKey,
	},m

}

func (agent *Agent) Sender() *principal.Principal {
	sha3.New224()
	pub := agent.key.PubKey.SerializeUncompressed()
	sender := principal.NewSelfAuthenticating(pub)
	return &sender
}

func (agent *Agent) getExpiryDate() time.Time {
	return time.Now().Add(agent.ingressExpiry)
}

func (agent *Agent) queryEndpoint(canisterID string, data []byte) (*QueryResponse, error) {
	resp, err := agent.client.query(canisterID, data)
	if err != nil {
		return nil, err
	}
	result := new(QueryResponse)
	err = cbor.Unmarshal(resp, result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (agent *Agent) callEndpoint(canisterID string, reqId RequestID, data []byte) (RequestID, error) {
	return agent.client.call(canisterID, reqId, data)
}

func (agent *Agent) readStateEndpoint(canisterID string, data []byte) ([]byte, error) {
	return agent.client.readState(canisterID, data)
}

func (agent *Agent) QueryRaw(cansiterID, methodName string, arg []byte) ([]idl.Type, []interface{}, string, error) {

	canisterID, err := principal.Decode(cansiterID)
	if err != nil {
		return nil, nil, "", err
	}

	//ingressExpiry,err := idl.Encode([]idl.Type{new(idl.Nat)},[]interface{}{big.NewInt(agent.getExpiryDate().UnixNano())})
	req := Request{
		Type:          RequestTypeQuery,
		Sender:        *agent.Sender(),
		CanisterID:    canisterID,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}
	_, data, err := agent.signRequest(&req)
	if err != nil {
		return nil, nil, "", err
	}
	resp, err := agent.queryEndpoint(cansiterID, data)
	if err != nil {
		return nil, nil, "", err
	}
	if resp.Status == "replied" {
		types, values, err := idl.Decode(resp.Reply["arg"])
		if err != nil {
			return nil, nil, "", err
		}
		return types, values, "", nil
	} else if resp.Status == "rejected" {
		return nil, nil, resp.RejectMsg, nil
	}
	return nil, nil, "", nil
}

func (agent *Agent) UpdateRaw(cansiterID, methodName string, arg []byte) ([]idl.Type, []interface{}, error) {
	canisterID, err := principal.Decode(cansiterID)
	if err != nil {
		return nil, nil, err
	}
	req := Request{
		Type:          RequestTypeCall,
		Sender:        *agent.Sender(),
		CanisterID:    canisterID,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().Nanosecond()),
	}
	requestID, data, err := agent.signRequest(&req)
	if err != nil {
		return nil, nil, err
	}
	_, err = agent.callEndpoint(cansiterID, *requestID, data)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("update request id:", hex.EncodeToString(requestID[:]))
	//poll requestID to get result
	//todo:这个时间写成配置之后
	result, err := agent.poll(cansiterID, *requestID, time.Second, time.Second*30)
	if err != nil {
		return nil, nil, err
	}
	types, values, err := idl.Decode(result)
	if err != nil {
		return nil, nil, err
	}
	return types, values, nil

}

func (agent *Agent) poll(cansiterID string, requestID RequestID, delay time.Duration, timeout time.Duration) ([]byte, error) {
	finalStatus := ""
	var finalCert []byte
	timer := time.NewTimer(timeout)
	triker := time.NewTicker(delay)

	for {
		select {
		case <-triker.C:
			status, cert, err := agent.requestStatusRaw(cansiterID, requestID)
			if err != nil {
				fmt.Printf("can not request status raw with error : %v", err)
			}
			finalStatus = string(status)
			finalCert = append(finalCert, cert...)
			if finalStatus == "replied" || finalStatus == "done" || finalStatus == "rejected" {
				break
			}

		case <-timer.C:
			break
		}
	}
	if finalStatus == "replied" {
		paths := [][]byte{[]byte("request_status"), requestID[:], []byte("reply")}
		res, err := LookUp(paths, finalCert)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil,fmt.Errorf("call poll fail with status %v",finalStatus)
}

func (agent *Agent) requestStatusRaw(cansiterID string, requestId RequestID) ([]byte, []byte, error) {
	//todo:回头看看这么编码行不行
	paths := [][]byte{[]byte("request_status")}
	paths = append(paths, requestId[:])
	cert, err := agent.readStateRaw(cansiterID, paths)
	if err != nil {
		return nil, nil, err
	}
	//print(cert)
	paths = append(paths, []byte("status"))
	status, err := LookUp(paths, cert)
	return status, cert, err
}

func (agent *Agent) readStateRaw(cansiterID string, paths [][]byte) ([]byte, error) {
	req := Request{
		Type:          RequestTypeReadState,
		Sender:        *agent.Sender(),
		Paths:         paths,
		IngressExpiry: uint64(agent.getExpiryDate().Nanosecond()),
	}
	_, data, err := agent.signRequest(&req)
	if err != nil {
		return nil, err
	}
	resp, err := agent.readStateEndpoint(cansiterID, data)
	if err != nil {
		return nil, err
	}
	result := struct {
		certificate []byte `cbor:"certificate"`
	}{}
	err = cbor.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result.certificate, nil
}

func (agent *Agent) signRequest(req *Request) (*RequestID, []byte, error) {
	requestID := NewRequestID(*req)
	msg := []byte(IC_REQUEST_DOMAIN_SEPARATOR)
	msg = append(msg, requestID[:]...)
	sig, err := agent.key.Sign(msg)
	if err != nil {
		return nil, nil, err
	}
	envelope := Envelope{
		Content:      req,
		SenderPubkey: agent.key.PubKey.SerializeUncompressed(),
		SenderSig:    sig.Serialize(),
	}
	mashaledEnvelope, err := cbor.Marshal(envelope)
	if err != nil {
		return nil, nil, err
	}
	return &requestID, mashaledEnvelope, nil
}
