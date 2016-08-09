package http

import (
	"bytes"
	"encoding/json"
	"io"
	basehttp "net/http"
	"time"

	"github.com/mantzas/adaptlog"
)

// RestClient interface
type RestClient interface {
	Get(url string) (*basehttp.Response, error)
	Post(url string, payload interface{}) (*basehttp.Response, error)
	Put(url string, payload interface{}) (*basehttp.Response, error)
	Delete(url string) (*basehttp.Response, error)
}

// RestClientImpl defines a http client
type RestClientImpl struct {
	timeout time.Duration
	logger  adaptlog.LevelLogger
}

// NewRestClientImpl creates a new client
func NewRestClientImpl(timeout time.Duration) *RestClientImpl {
	return &RestClientImpl{timeout, adaptlog.NewStdLevelLogger("Client")}
}

// Get sends a HTTP GET
func (c *RestClientImpl) Get(url string) (*basehttp.Response, error) {

	req, err := c.createRequest(basehttp.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Post sends a HTTP POST with a json body
func (c *RestClientImpl) Post(url string, payload interface{}) (*basehttp.Response, error) {

	req, err := c.createJSONRequest(basehttp.MethodPost, url, payload)

	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Put sends a HTTP PUT with a json body
func (c *RestClientImpl) Put(url string, payload interface{}) (*basehttp.Response, error) {

	req, err := c.createJSONRequest(basehttp.MethodPut, url, payload)

	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Delete sends a HTTP Delete
func (c *RestClientImpl) Delete(url string) (*basehttp.Response, error) {

	req, err := c.createRequest(basehttp.MethodDelete, url, nil)

	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

func (c *RestClientImpl) createJSONRequest(httpMethod string, url string, payload interface{}) (*basehttp.Request, error) {

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		c.logger.Errorf("Failed to marshal payload to JSON. %s", err)
		return nil, err
	}

	req, err := c.createRequest(httpMethod, url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *RestClientImpl) createRequest(httpMethod string, url string, body io.Reader) (*basehttp.Request, error) {

	req, err := basehttp.NewRequest(httpMethod, url, body)

	if err != nil {
		c.logger.Errorf("Failed to create request. %s", err)
		return nil, err
	}

	return req, nil
}

func (c *RestClientImpl) sendRequest(request *basehttp.Request) (*basehttp.Response, error) {

	client := &basehttp.Client{
		Timeout: c.timeout,
	}

	return client.Do(request)
}
