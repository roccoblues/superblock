package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// DefaultBaseURL used for API requests. Use WithBaseURL() to change.
const DefaultBaseURL = "https://api.twitter.com/2"

// Client represents is a Twitter API client.
type Client struct {
	httpClient     HTTPClient
	logger         Logger
	rateLimitRetry bool
	baseURL        string
}

// HTTPClient models the http client interface.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Logger is a simple logger interface.
type Logger interface {
	Printf(format string, v ...interface{})
}

// HTTPError contains information about a failed http request.
type HTTPError struct {
	Response       *http.Response
	StatusCode     int
	Status         string
	RateLimitReset time.Time // time when the rate-limiting will reset and it's safe to retry requests
}

// Error returns a string representation of the HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d %s", e.Response.StatusCode, e.Response.Status)
}

// ClientOption sets some additional options on a client.
type ClientOption func(*Client)

// meta fields on API responses
type meta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// NewClient returns a new Client struct.
func NewClient(httpClient HTTPClient, options ...ClientOption) *Client {
	c := &Client{
		httpClient: httpClient,
		baseURL:    DefaultBaseURL,
		logger:     log.New(ioutil.Discard, "", 0),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// WithBaseURL sets the base URL for request.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) { c.baseURL = url }
}

// WithLogger sets the logger used in the client. By default no log output is emitted.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) { c.logger = logger }
}

// EnableRateLimitRetry enables wait and retry on 429 responses. Default is false.
func EnableRateLimitRetry() ClientOption {
	return func(c *Client) { c.rateLimitRetry = true }
}

// NewRequest returns a http.Request for the given path.
// If a payload is provided it will get JSON encoded.
func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	if payload == nil {
		return http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), nil)
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}
	return http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), buf)
}

// DoRequest makes a request to the API and unmarshales the response into v.
func (c *Client) DoRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")

	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	for {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// only 200 and 201 responses contain a body that can be unmashalled into v
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			return json.NewDecoder(resp.Body).Decode(v)
		}

		httpErr := HTTPError{
			Response:   resp,
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return &httpErr
		}

		if reset := resp.Header.Get("x-rate-limit-reset"); reset != "" {
			ts, err := strconv.ParseInt(reset, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to convert rate limit reset value: %w", err)
			}
			httpErr.RateLimitReset = time.Unix(ts, 0)
		}

		if !c.rateLimitRetry || httpErr.RateLimitReset.IsZero() {
			return &httpErr
		}

		c.logger.Printf("%s %s: Rate limit exceeded! Waiting until %s before retrying request.\n",
			req.Method,
			req.URL.Path,
			httpErr.RateLimitReset,
		)
		time.Sleep(time.Until(httpErr.RateLimitReset))
	}
}
