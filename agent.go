package agent

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/stopWarByWar/ic-agent/internal/identity"
	"github.com/stopWarByWar/ic-agent/internal/idl"
	"github.com/stopWarByWar/ic-agent/internal/principal"
)

type Agent struct {
	client        *Client
	identity      *identity.Identity
	ingressExpiry time.Duration
	rootKey       []byte //ICP root identity
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
	//
	//fmt.Println(hex.EncodeToString(resp))
	//fmt.Println(result)
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
	//fmt.Println("data:", hex.EncodeToString(data))
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
	//fmt.Println("update request id:", hex.EncodeToString(requestID[:]))
	//fmt.Println("data:", hex.EncodeToString(data))
	//data,_ = hex.DecodeString("a367636f6e74656e74a66c726571756573745f747970656463616c6c6673656e646572581d8139de9ec81d50d862a956dd95e8d705e462f9e8df206ff4fe498739026b63616e69737465725f69644a0000000000f010ec01016b6d6574686f645f6e616d65687472616e73666572636172674f4449444c0002687d010080c8afa0256e696e67726573735f6578706972791b16c5481082f22e006d73656e6465725f7075626b6579582c302a300506032b6570032100ec172b93ad5e563bf4932c70e1245034c35467ef2efd4d64ebf819683467e2bf6a73656e6465725f736967584028677b532f1baaa31619923381f3a3d0be33ff4212dfc984dd3649fd0fbec938edcb30f72d59d2f01d1697a6eda3ac1c7de231974c0e3239f37ddbdf5b98f60a")
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
			finalCert = append(finalCert, cert...)
			//fmt.Println(finalCert)
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
	//todo:回头看看这么编码行不行
	paths := [][][]byte{{[]byte("request_status"), requestId[:]}}
	cert, err := agent.readStateRaw(canisterID, paths)
	//fmt.Println("err ", err)
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
	//fmt.Println("Read state Data:", hex.EncodeToString(data))
	//data,_ = hex.DecodeString("a367636f6e74656e74a66c726571756573745f747970656463616c6c6673656e646572581d8139de9ec81d50d862a956dd95e8d705e462f9e8df206ff4fe498739026b63616e69737465725f69644a0000000000f010ec01016b6d6574686f645f6e616d65687472616e73666572636172674f4449444c0002687d010080c8afa0256e696e67726573735f6578706972791b16c52f08c78464006d73656e6465725f7075626b6579582c302a300506032b6570032100ec172b93ad5e563bf4932c70e1245034c35467ef2efd4d64ebf819683467e2bf6a73656e6465725f73696758404093a2371ab4fc3a2e742bed6ed0606f14f35cffd40235e15ae47aa628a0375f252eaee777149d2326ec99900b135a2f2b291df0f3752223e407065947141307")

	if err != nil {
		return nil, err
	}
	resp, err := agent.readStateEndpoint(canisterID, data)
	if err != nil {
		return nil, err
	}
	// result := struct {
	// 	certificate []byte `cbor:"certificate"`
	// }{}
	result := map[string][]byte{}
	//result := []byte{}
	
	err = cbor.Unmarshal(resp, &result)
	if err != nil {
		//return nil, err
		return nil, err
	}
	//fmt.Println(result)
	// result_again := map[string][]byte{}
	// err = cbor.Unmarshal(result["certificate"], &result_again)
	// if err != nil {
	// 	//return nil, err
	// 	return nil, err
	// }

	fmt.Println("result!!!!   ", result["certificate"])
	return result["certificate"], nil
	//return result["certificate"], nil
}

func (agent *Agent) signRequest(req Request) (*RequestID, []byte, error) {
	requestID := NewRequestID(req)
	//fmt.Println(hex.EncodeToString(requestID[:]), "   req_id")
	msg := []byte(IC_REQUEST_DOMAIN_SEPARATOR)
	msg = append(msg, requestID[:]...)
	sig, err := agent.identity.Sign(msg)
	//fmt.Println(hex.EncodeToString(sig))
	if err != nil {
		return nil, nil, err
	}
	envelope := Envelope{
		Content:      req,
		SenderPubkey: agent.identity.PubKeyBytes(),
		SenderSig:    sig,
	}

	//fmt.Println("envelope:",envelope)
	marshaledEnvelope, err := cbor.Marshal(envelope)
	if err != nil {
		return nil, nil, err
	}
	return &requestID, marshaledEnvelope, nil
}
