package truto

import (
	"net/http"
	"strings"
	"time"
)

// Option configures the Client (when passed to NewClient) or a single request
// (when passed to any API method). Request options override client options for
// that call only.
type Option func(*settings)

// RetryInfo is passed to the OnRetry hook before each retry.
type RetryInfo struct {
	// RetryCount is the number of the retry about to happen, starting at 1.
	RetryCount int
	// StatusCode is the status that triggered the retry, or 0 for a network error.
	StatusCode int
	// Wait is how long the client sleeps before retrying.
	Wait time.Duration
}

// RateLimitInfo is passed to the OnRateLimited hook when the API answers 429.
type RateLimitInfo struct {
	// RetryAfter is how long the client waits, from the Retry-After header.
	RetryAfter time.Duration
	// RetryCount is the number of retries already made for this request.
	RetryCount int
}

type settings struct {
	baseURL           string
	httpClient        *http.Client
	timeout           time.Duration
	maxRetries        int
	defaultRetryAfter time.Duration
	noRetry           bool
	retryStatusCodes  []int
	headers           http.Header
	query             Params
	onRetry           func(RetryInfo)
	onRateLimited     func(RateLimitInfo)
}

// Defaults used when no option overrides them.
const (
	DefaultBaseURL    = "https://api.truto.one"
	DefaultTimeout    = 30 * time.Second
	DefaultMaxRetries = 3
	DefaultRetryAfter = 10 * time.Second
	backoffBase       = 500 * time.Millisecond
	backoffCap        = 10 * time.Second
)

var baseRetryStatusCodes = []int{408, 409, 425, 429, 502, 503, 504}

func defaultSettings() settings {
	return settings{
		baseURL:           DefaultBaseURL,
		httpClient:        http.DefaultClient,
		timeout:           DefaultTimeout,
		maxRetries:        DefaultMaxRetries,
		defaultRetryAfter: DefaultRetryAfter,
	}
}

// clone returns a copy whose maps and slices can be changed without touching s.
func (s settings) clone() settings {
	c := s
	c.headers = s.headers.Clone()
	c.retryStatusCodes = append([]int(nil), s.retryStatusCodes...)
	if s.query != nil {
		c.query = make(Params, len(s.query))
		for k, v := range s.query {
			c.query[k] = v
		}
	}
	return c
}

// WithBaseURL points the client at another Truto API host.
func WithBaseURL(url string) Option {
	return func(s *settings) { s.baseURL = strings.TrimRight(url, "/") }
}

// WithHTTPClient sets the *http.Client used to send requests.
func WithHTTPClient(c *http.Client) Option {
	return func(s *settings) {
		if c != nil {
			s.httpClient = c
		}
	}
}

// WithTimeout limits each attempt of a request. A timed-out attempt is not
// retried; the call returns an error that matches ErrTimeout.
func WithTimeout(d time.Duration) Option {
	return func(s *settings) { s.timeout = d }
}

// WithMaxRetries sets how many times a failed request is retried.
func WithMaxRetries(n int) Option {
	return func(s *settings) {
		if n >= 0 {
			s.maxRetries = n
		}
	}
}

// WithDefaultRetryAfter sets the wait used for a 429 response that has no
// usable Retry-After header.
func WithDefaultRetryAfter(d time.Duration) Option {
	return func(s *settings) { s.defaultRetryAfter = d }
}

// WithNoRetry turns retries off.
func WithNoRetry() Option {
	return func(s *settings) { s.noRetry = true }
}

// WithRetryStatusCodes adds status codes that are retried, on top of
// 408, 409, 425, 429, 502, 503 and 504.
func WithRetryStatusCodes(codes ...int) Option {
	return func(s *settings) { s.retryStatusCodes = append(s.retryStatusCodes, codes...) }
}

// WithHeader adds a request header.
func WithHeader(key, value string) Option {
	return func(s *settings) {
		if s.headers == nil {
			s.headers = http.Header{}
		}
		s.headers.Add(key, value)
	}
}

// WithQuery adds query parameters. Nested maps and slices are encoded the way
// the Truto API expects (filter[status]=open, ids[0]=a). Calling it more than
// once merges the parameters; later values win.
func WithQuery(q Params) Option {
	return func(s *settings) {
		if s.query == nil {
			s.query = Params{}
		}
		for k, v := range q {
			s.query[k] = v
		}
	}
}

// WithOnRetry registers a hook called before each retry.
func WithOnRetry(fn func(RetryInfo)) Option {
	return func(s *settings) { s.onRetry = fn }
}

// WithOnRateLimited registers a hook called when the API answers 429, before
// the client waits.
func WithOnRateLimited(fn func(RateLimitInfo)) Option {
	return func(s *settings) { s.onRateLimited = fn }
}
