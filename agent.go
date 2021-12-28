package agent

import (
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/internal/key"
	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/principal-go"
	"github.com/fxamacker/cbor/v2"
	"time"
)

type Agent struct {
	client        *Client
	identity      *identity.Identity
	ingressExpiry time.Duration
	rootKey       []byte //ICP root key
}

func New(anonymous bool, privKey string) *Agent {
	c := NewClient("https://ic0.app")
	//todo:是否需要从ic拉取rootKey信息
	status, _ := c.Status()
	pbBytes, _ := hex.DecodeString(privKey)
	id := identity.New(anonymous, pbBytes)

	ingressExpiry := time.Second * 10
	return &Agent{
		client:        &c,
		identity:      id,
		ingressExpiry: ingressExpiry,
		rootKey:       status.RootKey,
	}
}

func (agent *Agent) Sender() *principal.Principal {
	if agent.identity.Anonymous == true {
		return &principal.AnonymousID
	}
	sender := principal.NewSelfAuthenticating(agent.identity.PubKeyBytes())
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

	fmt.Println(hex.EncodeToString(resp))
	fmt.Println(result)
	return result, nil
}

func (agent *Agent) callEndpoint(canisterID string, reqId RequestID, data []byte) (RequestID, error) {
	return agent.client.call(canisterID, reqId, data)
}

func (agent *Agent) readStateEndpoint(canisterID string, data []byte) ([]byte, error) {
	return agent.client.readState(canisterID, data)
}

func (agent *Agent) QueryRaw(canisterID, methodName string, arg []byte) ([]idl.Type, []interface{}, string, error) {

	canisterIDPrincipal, err := principal.Decode(canisterID)
	if err != nil {
		return nil, nil, "", err
	}
	req := Request{
		Type:          RequestTypeQuery,
		Sender:        *agent.Sender(),
		CanisterID:    canisterIDPrincipal,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}
	_, data, err := agent.signRequest(req)
	if err != nil {
		return nil, nil, "", err
	}
	//data,_ = hex.DecodeString("a367636f6e74656e74a66c726571756573745f747970656571756572796673656e646572581d8139de9ec81d50d862a956dd95e8d705e462f9e8df206ff4fe498739026b63616e69737465725f69644a0000000000f010ec01016b6d6574686f645f6e616d65646e616d6563617267464449444c00006e696e67726573735f6578706972791b16c4df04ce1440006d73656e6465725f7075626b6579582c302a300506032b6570032100ec172b93ad5e563bf4932c70e1245034c35467ef2efd4d64ebf819683467e2bf6a73656e6465725f73696758408cd1def386a59e59d29adf00cbc735dc8a3dba48c715e6f69a6f9e0ede937c404e6fd7c1bd2583f49a63b597fa1d90fd3feefc88d4eb731ff2c96af51ea3ac02")
	resp, err := agent.queryEndpoint(canisterID, data)
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

func (agent *Agent) UpdateRaw(canisterID, methodName string, arg []byte) ([]idl.Type, []interface{}, error) {
	canisterIDPrincipal, err := principal.Decode(canisterID)
	if err != nil {
		return nil, nil, err
	}
	req := Request{
		Type:          RequestTypeCall,
		Sender:        *agent.Sender(),
		CanisterID:    canisterIDPrincipal,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().Nanosecond()),
	}

	requestID, data, err := agent.signRequest(req)
	if err != nil {
		return nil, nil, err
	}
	_, err = agent.callEndpoint(canisterID, *requestID, data)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("update request id:", hex.EncodeToString(requestID[:]))
	//poll requestID to get result
	//todo:这个时间写成配置之后
	result, err := agent.poll(canisterID, *requestID, time.Second, time.Second*30)
	if err != nil {
		return nil, nil, err
	}
	types, values, err := idl.Decode(result)
	if err != nil {
		return nil, nil, err
	}
	return types, values, nil

}

func (agent *Agent) poll(canisterID string, requestID RequestID, delay time.Duration, timeout time.Duration) ([]byte, error) {
	finalStatus := ""
	var finalCert []byte
	timer := time.NewTimer(timeout)
	ticker := time.NewTicker(delay)

	for {
		select {
		case <-ticker.C:
			status, cert, err := agent.requestStatusRaw(canisterID, requestID)
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
	return nil, fmt.Errorf("call poll fail with status %v", finalStatus)
}

func (agent *Agent) requestStatusRaw(canisterID string, requestId RequestID) ([]byte, []byte, error) {
	//todo:回头看看这么编码行不行
	paths := [][]byte{[]byte("request_status")}
	paths = append(paths, requestId[:])
	cert, err := agent.readStateRaw(canisterID, paths)
	if err != nil {
		return nil, nil, err
	}
	//print(cert)
	paths = append(paths, []byte("status"))
	status, err := LookUp(paths, cert)
	return status, cert, err
}

func (agent *Agent) readStateRaw(canisterID string, paths [][]byte) ([]byte, error) {
	req := Request{
		Type:          RequestTypeReadState,
		Sender:        *agent.Sender(),
		Paths:         paths,
		IngressExpiry: uint64(agent.getExpiryDate().Nanosecond()),
	}

	_, data, err := agent.signRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := agent.readStateEndpoint(canisterID, data)
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

func (agent *Agent) signRequest(req Request) (*RequestID, []byte, error) {
	requestID := NewRequestID(req)
	msg := []byte(IC_REQUEST_DOMAIN_SEPARATOR)
	msg = append(msg, requestID[:]...)
	sig, err := agent.identity.Sign(msg)
	if err != nil {
		return nil, nil, err
	}
	envelope := Envelope{
		Content:      req,
		SenderPubkey: agent.identity.PubKeyBytes(),
		SenderSig:    sig,
	}
	marshaledEnvelope, err := cbor.Marshal(envelope)
	if err != nil {
		return nil, nil, err
	}
	return &requestID, marshaledEnvelope, nil
}
