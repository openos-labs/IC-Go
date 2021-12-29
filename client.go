package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client http.Client
	host   string
}

func NewClient(host string) Client {
	return Client{
		client: http.Client{},
		host:   host,
	}
}

func (c *Client) Status() (Status, error) {
	raw, err := c.get("/api/v2/status")
	if err != nil {
		return Status{}, err
	}
	status := new(Status)
	err = status.UnmarshalCBOR(raw)
	return *status, err
}

func (c *Client) query(canisterId string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	endpoint := c.host + "/api/v2/canister/" + canisterId + "/query"
	fmt.Println("post url:", endpoint)
	resp, err := c.client.Post(endpoint, "application/cbor", buffer)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		fmt.Println(
			"status:", resp.Status, "\n",
			"StatusCode:", resp.StatusCode, "\n",
			"Proto:", resp.Proto, "\n",
			"ProtoMajor:", resp.ProtoMajor, "\n",
			"ProtoMinor:", resp.ProtoMinor, "\n",
			"Header:", resp.Header, "\n",
			"Body:", resp.Body, "\n",
			"ContentLength:", resp.ContentLength, "\n",
			"TransferEncoding:", resp.TransferEncoding, "\n",
			"Request:", resp.Request,
		)
		return nil, fmt.Errorf("fail to post ic with status: %v", resp.Status)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) call(canisterId string, reqId RequestID, data []byte) (RequestID, error) {
	buffer := bytes.NewBuffer(data)
	endpoint := c.host + "/api/v2/canister/" + canisterId + "/call"
	fmt.Println("endpoint:", endpoint)
	contentType := "application/cbor"
	resp, err := c.client.Post(endpoint, contentType, buffer)
	if err != nil {
		return reqId, err
	}
	if resp.StatusCode != 200 {
		return reqId, fmt.Errorf("fail to call ic with status: %v", resp.Status)
	}
	return reqId, nil
}

func (c *Client) readState(canisterId string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	endpoint := c.host + "/api/v2/canister/" + canisterId + "/read_state"
	fmt.Println("endpoint:", endpoint)
	contentType := "application/cbor"
	resp, err := c.client.Post(endpoint, contentType, buffer)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("fail to read state with status: %v", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (c Client) get(path string) ([]byte, error) {
	a := c.host + path
	resp, err := c.client.Get(a)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}