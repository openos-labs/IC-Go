package agent_test

import (
	"testing"

	"github.com/aviate-labs/agent-go"
)

func TestClientStatus(t *testing.T) {
	c := agent.NewClient("https://ic0.app")
	status, _ := c.Status()
	t.Log(status.Version)
	// Output:
	// 0.18.0
}
