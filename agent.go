package agent

import (
	"encoding/hex"
	"fmt"
	"github.com/openos-labs/IC-Go/utils/identity"
	"github.com/openos-labs/IC-Go/utils/idl"
	"github.com/openos-labs/IC-Go/utils/principal"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type Agent struct {
	client        *Client
	identity      *identity.Identity
	ingressExpiry time.Duration
	rootKey       []byte //ICP root identity
}

func NewFromPem(anonymous bool, pemPath string) (*Agent, error) {
	var id *identity.Identity
	c := NewClient("https://ic0.app")
	//todo:是否需要从ic拉取rootKey信息
	status, _ := c.Status()

	if anonymous == true {
		id = &identity.Identity{
			Anonymous: anonymous,
		}
	}
	privKey, err := identity.FromPem(pemPath)
	if err != nil {
		return nil, err
	}
	pubKey := privKey.Public()

	id = &identity.Identity{
		anonymous,
		privKey,
		pubKey,
	}
	ingressExpiry := time.Second * 10
	return &Agent{
		client:        &c,
		identity:      id,
		ingressExpiry: ingressExpiry,
		rootKey:       status.RootKey,
	}, nil
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

func (agent *Agent) Identity() *identity.Identity {
	return agent.identity
}

func (agent *Agent) Sender() principal.Principal {
	if agent.identity.Anonymous == true {
		return principal.AnonymousID
	}
	sender := principal.NewSelfAuthenticating(agent.identity.PubKeyBytes())
	return sender
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

func (agent *Agent) QueryRaw(canisterID, methodName string, arg []byte) ([]idl.Type, []interface{}, string, error) {
	canisterIDPrincipal, err := principal.Decode(canisterID)
	if err != nil {
		return nil, nil, "", err
	}
	req := Request{
		Type:          RequestTypeQuery,
		Sender:        agent.Sender(),
		CanisterID:    canisterIDPrincipal,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}
	_, data, err := agent.signRequest(req)
	if err != nil {
		return nil, nil, "", err
	}

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
		Sender:        agent.Sender(),
		CanisterID:    canisterIDPrincipal,
		MethodName:    methodName,
		Arguments:     arg,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}

	requestID, data, err := agent.signRequest(req)
	if err != nil {
		return nil, nil, err
	}

	_, err = agent.callEndpoint(canisterID, *requestID, data)
	if err != nil {
		return nil, nil, err
	}
	//poll requestID to get result
	//todo:这个时间写成配置之后
	result, err := agent.poll(canisterID, *requestID, time.Second, time.Second*10)
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
	stopped := true
	for stopped {
		select {
		case <-ticker.C:
			status, cert, err := agent.requestStatusRaw(canisterID, requestID)
			if err != nil {
				fmt.Printf("can not request status raw with error : %v\n", err)
			}
			finalStatus = string(status)
			finalCert = cert
			if finalStatus == "replied" || finalStatus == "done" || finalStatus == "rejected" {
				stopped = false
			}
		case <-timer.C:
			stopped = false
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
	defer timer.Stop()
	defer ticker.Stop()
	return nil, fmt.Errorf("call poll fail with status %v", finalStatus)
}

func (agent *Agent) requestStatusRaw(canisterID string, requestId RequestID) ([]byte, []byte, error) {

	paths := [][][]byte{{[]byte("request_status"), requestId[:]}}
	cert, err := agent.readStateRaw(canisterID, paths)

	if err != nil {
		return nil, nil, err
	}
	path := [][]byte{[]byte("request_status"), requestId[:], []byte("status")}
	status, err := LookUp(path, cert)
	return status, cert, err
}

func (agent *Agent) readStateRaw(canisterID string, paths [][][]byte) ([]byte, error) {
	req := Request{
		Type:          RequestTypeReadState,
		Sender:        agent.Sender(),
		Paths:         paths,
		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
	}

	_, data, err := agent.signRequest(req)

	if err != nil {
		return nil, err
	}
	resp, err := agent.readStateEndpoint(canisterID, data)
	if err != nil {
		return nil, err
	}

	result := map[string][]byte{}

	err = cbor.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result["certificate"], nil
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

func (agent *Agent) GetCanisterControllers(canisterID string) ([]principal.Principal, error) {
	info, err := agent.getCanisterInfo(canisterID, "controllers")
	if err != nil {
		return nil, err
	}
	var mResult [][]byte
	var result []principal.Principal
	err = cbor.Unmarshal(info, &mResult)
	if err != nil {
		return nil, err
	}
	for _, p := range mResult {
		result = append(result, principal.New(p))
	}
	return result, nil
}

func (agent *Agent) GetCanisterModule(canisterID string) ([]byte, error) {
	return agent.getCanisterInfo(canisterID, "module_hash")
}

func (agent Agent) getCanisterInfo(canisterID, subPath string) ([]byte, error) {
	canisterBytes, err := principal.Decode(canisterID)
	if err != nil {
		return nil, err
	}
	paths := [][][]byte{{[]byte("canister"), canisterBytes, []byte(subPath)}}
	cert, err := agent.readStateRaw(canisterID, paths)
	if err != nil {
		return nil, err
	}
	path := [][]byte{[]byte("canister"), canisterBytes, []byte(subPath)}
	return LookUp(path, cert)
}
