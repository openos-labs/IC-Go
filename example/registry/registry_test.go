package registry

import (
	"encoding/binary"
	"fmt"
	agent "github.com/mix-labs/IC-Go"
	"github.com/mix-labs/IC-Go/utils/principal"
	"testing"
)

func TestGetRoutingTable(t *testing.T) {
	a := agent.New(true, "")
	routingTable, err := GetRoutingTable(a)
	if err != nil {
		t.Log("error:", err)
	}
	for i, entry := range routingTable.Entries {
		t.Log("subnet index:", i)
		fmt.Println("subnetID:", principal.New(entry.SubnetId.PrincipalId.Raw).Encode())
		fmt.Println("start canister ID:", principal.New(entry.Range.StartCanisterId.PrincipalId.Raw).Encode())
		fmt.Println("start canister ID raw:", entry.Range.StartCanisterId.PrincipalId.Raw)
		fmt.Println("start canister ID to uint64:", binary.BigEndian.Uint64(entry.Range.StartCanisterId.PrincipalId.Raw[:8]))


		fmt.Println("end canister ID:", principal.New(entry.Range.EndCanisterId.PrincipalId.Raw).Encode())
		fmt.Println("end canister ID raw:", entry.Range.EndCanisterId.PrincipalId.Raw)
		fmt.Println("end canister ID to uint64:", binary.BigEndian.Uint64(entry.Range.EndCanisterId.PrincipalId.Raw[:8]))
	}
}

func TestGetSubnetList(t *testing.T) {
	a := agent.New(true, "")
	subnetList, err := GetSubnetList(a)
	if err != nil {
		t.Log("error:", err)
	}
	for _, entry := range subnetList.Subnets {
		t.Log("subnetID:", principal.New(entry).Encode())
	}
}

func TestGetSubnetRecord(t *testing.T) {
	a := agent.New(true, "")
	subnet, err := GetSubnetRecord(a, "tdb26-jop6k-aogll-7ltgs-eruif-6kk7m-qpktf-gdiqx-mxtrf-vb5e6-eqe")
	if err != nil {
		t.Log("error:", err)
	}

	t.Log("subnet Type:",subnet.SubnetType)
	t.Log("max canister amount",subnet.MaxNumberOfCanisters)
	for _, node := range subnet.Membership {
		t.Log("node:", principal.New(node).Encode())
	}
}

func TestGetNodeInfo(t *testing.T) {
	a := agent.New(true, "")
	node, err := GetNodeInfo(a, "btuxg-lwlbn-43hlo-iag4h-plf64-b3u7d-ugvay-nbvrl-jkhlx-nhvw4-gae")
	if err != nil {
		t.Log("error:", err)
	}
	t.Log("operator:", principal.New(node.NodeOperatorId).Encode())
}

func TestGetOperatorInfo(t *testing.T) {
	a := agent.New(true, "")
	op, err := GetOperatorInfo(a, "redpf-rrb5x-sa2it-zhbh7-q2fsp-bqlwz-4mf4y-tgxmj-g5y7p-ezjtj-5qe")
	if err != nil {
		t.Log("error:", err)
	}
	t.Log("operator:", principal.New(op.NodeOperatorPrincipalId).Encode())
	t.Log("provider:", principal.New(op.NodeProviderPrincipalId).Encode())

	t.Log("Node Allowance:", op.NodeAllowance)
	t.Log("Rewardable Nodes:", op.RewardableNodes)
}
