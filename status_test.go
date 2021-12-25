package agent_test

import (
	"net/url"
	"testing"

	"github.com/aviate-labs/agent-go"
)

var ic0, _ = url.Parse("https://ic0.app/")

func TestClientStatus(t *testing.T) {
	c := agent.NewClient(agent.ClientConfig{ic0})
	status, _ := c.Status()
	t.Log(status.Version)
	// Output:
	// 0.18.0
}
