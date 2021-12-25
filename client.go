package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	client http.Client
	config ClientConfig
}

func NewClient(cfg ClientConfig) Client {
	return Client{
		client: http.Client{},
		config: cfg,
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
	endpoint := c.url("/api/v2/canister/" + canisterId + "/query")
	contentType := "application/cbor"
	resp, err := c.client.Post(endpoint, contentType, buffer)
	if err != nil {
		return nil, err
	}

	fmt.Println(
		"status:",resp.Status,"\n",
			"StatusCode:",resp.StatusCode,"\n",
			"Proto:",resp.Proto,"\n",
			"ProtoMajor:",resp.ProtoMajor,"\n",
			"ProtoMinor:",resp.ProtoMinor,"\n",
			"Header:",resp.Header,"\n",
			"Body:",resp.Body,"\n",
			"ContentLength:",resp.ContentLength,"\n",
			"TransferEncoding:",resp.TransferEncoding,"\n",
			"Request:",resp.Request,
		)
	return io.ReadAll(resp.Body)
}

func (c *Client) call(canisterId string, reqId RequestID, data []byte) (RequestID, error) {
	buffer := bytes.NewBuffer(data)
	endpoint := c.url("/api/v2/canister/" + canisterId + "/call")
	contentType := "application/cbor"
	_, err := c.client.Post(endpoint, contentType, buffer)
	if err != nil {
		return reqId, err
	}
	return reqId, nil
}

func (c *Client) readState(canisterId string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	endpoint := c.url("/api/v2/canister/" + canisterId + "/read_state")
	contentType := "application/cbor"
	resp, err := c.client.Post(endpoint, contentType, buffer)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (c Client) get(path string) ([]byte, error) {
	resp, err := c.client.Get(c.url(path))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (c Client) url(p string) string {
	url := c.config.Host
	url.Path = path.Join(url.Path, p)
	return url.String()
}

type ClientConfig struct {
	Host *url.URL
}
