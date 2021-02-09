package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	client http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
	}
}

func (c *HttpClient) Do(
	ctx context.Context,
	shop, path, method string,
	headers map[string]string,
	in, out interface{},
) error {
	req, err := c.buildRequest(ctx, shop, path, method, in)
	if err != nil {
		return err
	}
	for k := range headers {
		req.Header.Add(k, headers[k])
	}
	return c.executeRequest(req, out)
}

func (c *HttpClient) executeRequest(req *http.Request, out interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%d: %s", resp.StatusCode, string(respBody))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *HttpClient) buildRequest(ctx context.Context, shop, path, method string, in interface{}) (*http.Request, error) {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("https://%s%s", shop, path),
		bytes.NewReader(reqBody),
	)
}
