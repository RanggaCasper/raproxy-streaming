package httpclient

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0"
)

// Client wraps fasthttp client for proxy requests
type Client struct {
	client *fasthttp.Client
}

// New creates a new HTTP client with custom configuration
func New(timeout, connectTimeout time.Duration, maxRedirects int) *Client {
	return &Client{
		client: &fasthttp.Client{
			ReadTimeout:         timeout,
			WriteTimeout:        timeout,
			MaxConnWaitTimeout:  connectTimeout,
			TLSConfig:           &tls.Config{InsecureSkipVerify: true},
			DisablePathNormalizing: true,
		},
	}
}

// Get performs a GET request with custom headers
func (c *Client) Get(url, referer string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("User-Agent", defaultUserAgent)
	
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	if err := c.client.Do(req, resp); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		statusCode := resp.StatusCode()
		fasthttp.ReleaseResponse(resp)
		return nil, &HTTPError{StatusCode: statusCode}
	}

	return resp, nil
}

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d", e.StatusCode)
}

// IsHTTPError checks if error is an HTTPError
func IsHTTPError(err error) (*HTTPError, bool) {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr, true
	}
	return nil, false
}
