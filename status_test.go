package agent_test

import (
	agent "github.com/openos-labs/IC-Go"
	"testing"
)

func TestClientStatus(t *testing.T) {
	c := agent.NewClient("https://ic0.app")
	status, _ := c.Status()
	t.Log(status.Version)
	// Output:
	// 0.18.0
}
