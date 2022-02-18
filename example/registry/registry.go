package registry

import (
	"errors"
	"github.com/golang/protobuf/proto"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/example/registry/proto/pb"
)

const CanisterId = "rwlgt-iiaaa-aaaaa-aaaaa-cai"

func getValue(agent *agent.Agent, key string) (*pb.RegistryGetValueResponse, error) {
	getValueResponse := new(pb.RegistryGetValueResponse)
	requestKey := []byte(key)
	request := &pb.RegistryGetValueRequest{
		Version: nil,
		Key:     requestKey,
	}
	requestBuf, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, ErrMsg, err := agent.QueryRaw(CanisterId, "get_value", requestBuf)
	if ErrMsg != "" {
		return nil, errors.New(ErrMsg)
	} else if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp, getValueResponse)
	if err != nil {
		return nil, err
	}
	return getValueResponse, nil
}

func GetRoutingTable(agent *agent.Agent) (*pb.RoutingTable, error) {
	routingTable := new(pb.RoutingTable)
	resp, err := getValue(agent, "routing_table")
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp.Value, routingTable)
	if err != nil {
		return nil, err
	}
	return routingTable, nil
}

func GetSubnetList(agent *agent.Agent) (*pb.SubnetListRecord, error) {
	subnetList := new(pb.SubnetListRecord)
	resp, err := getValue(agent, "subnet_list")
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp.Value, subnetList)
	if err != nil {
		return nil, err
	}
	return subnetList, nil
}

func GetSubnetRecord(agent *agent.Agent, subnetID string) (*pb.SubnetRecord, error) {
	subnetRecord := new(pb.SubnetRecord)
	key := "subnet_record_" + subnetID
	resp, err := getValue(agent, key)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp.Value, subnetRecord)
	if err != nil {
		return nil, err
	}
	return subnetRecord, nil
}

func GetNodeInfo(agent *agent.Agent, nodeID string) (*pb.NodeRecord, error) {
	nodeRecord := new(pb.NodeRecord)
	key := "node_record_" + nodeID
	resp, err := getValue(agent, key)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp.Value, nodeRecord)
	if err != nil {
		return nil, err
	}
	return nodeRecord, nil
}

func GetOperatorInfo(agent *agent.Agent, operatorID string) (*pb.NodeOperatorRecord, error) {
	operatorRecord := new(pb.NodeOperatorRecord)
	key := "node_operator_record_" + operatorID
	resp, err := getValue(agent, key)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(resp.Value, operatorRecord)
	if err != nil {
		return nil, err
	}
	return operatorRecord, nil
}
