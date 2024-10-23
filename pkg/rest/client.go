package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	JSONContentType string = "application/json"
	XMLContentType  string = "application/soap+xml"
)

// Client is a REST client.
type Client struct {
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new instance of Client.
func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		Client:  http.DefaultClient,
	}
}

// Response represents the response object for processing.
type Response struct {
	Body       []byte
	StatusCode int
}

// RequestOptions holds request customization options.
type RequestOptions struct {
	Headers http.Header
	Timeout time.Duration
}

// NewRequestOptions initializes RequestOptions with default values.
func NewRequestOptions() RequestOptions {
	return RequestOptions{Headers: make(http.Header)}
}

// httpClient returns the HTTP client, defaulting to http.DefaultClient if nil.
func (c *Client) httpClient() *http.Client {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	return c.Client
}

// Get performs an HTTP GET request and unmarshals the response.
func (c *Client) Get(ctx context.Context, reqURL string, _ interface{}, options *RequestOptions) (*Response, error) {
	return c.do(ctx, http.MethodGet, reqURL, nil, options)
}

// Post performs an HTTP POST request and unmarshals the response.
func (c *Client) Post(ctx context.Context, reqURL string, body interface{}, options *RequestOptions) (*Response, error) {
	res, err := c.do(ctx, http.MethodPost, reqURL, body, options)
	return res, err
}

// Put performs an HTTP PUT request and unmarshals the response.
func (c *Client) Put(ctx context.Context, reqURL string, body interface{}, options *RequestOptions) (*Response, error) {
	return c.do(ctx, http.MethodPut, reqURL, body, options)
}

// Patch performs an HTTP PATCH request and unmarshals the response.
func (c *Client) Patch(ctx context.Context, reqURL string, body interface{}, options *RequestOptions) (*Response, error) {
	return c.do(ctx, http.MethodPatch, reqURL, body, options)
}

// Delete performs an HTTP DELETE request and unmarshals the response.
func (c *Client) Delete(ctx context.Context, reqURL string, body interface{}, options *RequestOptions) (*Response, error) {
	return c.do(ctx, http.MethodDelete, reqURL, body, options)
}

// SendRequest executes an HTTP request and unmarshals the response.
func (c *Client) SendRequest(ctx context.Context, contentType string, req *http.Request, options *RequestOptions) (*Response, error) {
	var rsp Response
	httpClient := c.httpClient()

	req = req.WithContext(ctx)

	if options != nil {
		for key := range options.Headers {
			req.Header.Set(key, options.Headers.Get(key))
		}
		if options.Timeout > 0 {
			httpClient.Timeout = options.Timeout
		}
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("restclient: request failed: %w", err)
	}

	rsp.Body, err = io.ReadAll(resp.Body)
	rsp.StatusCode = resp.StatusCode
	if err != nil {
		return &rsp, fmt.Errorf("restclient: reading response body failed: %w", err)
	}

	return &rsp, nil
}

// do prepares and sends an HTTP request.
func (c *Client) do(ctx context.Context, method, reqURL string, body interface{}, options *RequestOptions) (*Response, error) {
	var reqBody []byte
	contentType := JSONContentType
	var err error

	if options != nil {
		contentType = options.Headers.Get("Content-Type")
	}

	if body != nil {
		switch contentType {
		case XMLContentType:
			reqBody, err = xml.Marshal(body)
		default:
			reqBody, err = json.Marshal(body)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("restclient: marshal failed: %w", err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, reqURL), bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("restclient: new request failed: %w", err)
	}

	return c.SendRequest(ctx, contentType, req, options)
}
