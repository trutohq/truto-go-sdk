package truto

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// request describes one API call before options are applied.
type request struct {
	method string
	path   string // without leading slash, already escaped
	body   any    // JSON-encoded when non-nil
}

// response is a successful (2xx) response.
type response struct {
	statusCode  int
	contentType string
	body        []byte
}

func (r *response) isJSON() bool { return strings.Contains(r.contentType, "json") }

// decodeInto decodes a JSON response into out. Empty and 204 responses leave
// out untouched.
func (r *response) decodeInto(out any) error {
	if out == nil || r.statusCode == http.StatusNoContent || len(bytes.TrimSpace(r.body)) == 0 {
		return nil
	}
	if !r.isJSON() {
		return fmt.Errorf("truto: expected a JSON response, got %q", r.contentType)
	}
	if err := json.Unmarshal(r.body, out); err != nil {
		return fmt.Errorf("truto: decoding response: %w", err)
	}
	return nil
}

func (c *Client) settingsFor(opts []Option) settings {
	s := c.settings.clone()
	for _, o := range opts {
		o(&s)
	}
	return s
}

// do sends a request with retries and returns the 2xx response.
func (c *Client) do(ctx context.Context, req request, opts []Option) (*response, error) {
	s := c.settingsFor(opts)

	url := s.baseURL + "/" + req.path
	qs, err := encodeQuery(s.query)
	if err != nil {
		return nil, fmt.Errorf("truto: encoding query: %w", err)
	}
	if qs != "" {
		url += "?" + qs
	}

	var payload []byte
	if req.body != nil {
		payload, err = json.Marshal(req.body)
		if err != nil {
			return nil, fmt.Errorf("truto: encoding body: %w", err)
		}
	}

	retryCodes := append(slices.Clone(baseRetryStatusCodes), s.retryStatusCodes...)
	maxRetries := s.maxRetries
	if s.noRetry {
		maxRetries = 0
	}

	for attempt := 0; ; attempt++ {
		resp, status, header, body, err := c.attempt(ctx, s, req.method, url, payload)
		if err == nil && status >= 200 && status < 300 {
			return &response{statusCode: status, contentType: header.Get("Content-Type"), body: body}, nil
		}
		if err != nil {
			// Timeouts and cancellations are final; other network errors retry.
			if errors.Is(err, ErrTimeout) || ctx.Err() != nil || attempt >= maxRetries {
				return nil, err
			}
		} else if attempt >= maxRetries || !slices.Contains(retryCodes, status) {
			return nil, newAPIError(req, resp, body)
		}

		var wait time.Duration
		if status == http.StatusTooManyRequests {
			wait = retryAfter(header, s.defaultRetryAfter)
			if s.onRateLimited != nil {
				s.onRateLimited(RateLimitInfo{RetryAfter: wait, RetryCount: attempt})
			}
		} else {
			wait = backoff(attempt)
		}
		if s.onRetry != nil {
			s.onRetry(RetryInfo{RetryCount: attempt + 1, StatusCode: status, Wait: wait})
		}
		if err := sleep(ctx, wait); err != nil {
			return nil, err
		}
	}
}

// attempt sends one HTTP request. On a transport error status is 0.
func (c *Client) attempt(ctx context.Context, s settings, method, url string, payload []byte) (*http.Response, int, http.Header, []byte, error) {
	actx := ctx
	cancel := func() {}
	if s.timeout > 0 {
		actx, cancel = context.WithTimeout(ctx, s.timeout)
	}
	defer cancel()

	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	hreq, err := http.NewRequestWithContext(actx, method, url, body)
	if err != nil {
		return nil, 0, nil, nil, fmt.Errorf("truto: building request: %w", err)
	}
	for k, vs := range s.headers {
		for _, v := range vs {
			hreq.Header.Add(k, v)
		}
	}
	hreq.Header.Set("Authorization", "Bearer "+c.token)
	if hreq.Header.Get("Accept") == "" {
		hreq.Header.Set("Accept", "application/json")
	}
	hreq.Header.Set("User-Agent", userAgent)
	if payload != nil && hreq.Header.Get("Content-Type") == "" {
		hreq.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.httpClient.Do(hreq)
	if err == nil {
		var data []byte
		data, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			return resp, resp.StatusCode, resp.Header, data, nil
		}
	}
	if ctx.Err() != nil {
		return nil, 0, nil, nil, ctx.Err()
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, 0, nil, nil, fmt.Errorf("%w: %s %s after %s", ErrTimeout, method, url, s.timeout)
	}
	return nil, 0, nil, nil, fmt.Errorf("truto: %s %s: %w", method, url, err)
}

func newAPIError(req request, resp *http.Response, body []byte) *Error {
	e := &Error{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		RawBody:    body,
		Header:     resp.Header,
		Method:     req.method,
		Path:       "/" + req.path,
	}
	var parsed any
	if strings.Contains(resp.Header.Get("Content-Type"), "json") && json.Unmarshal(body, &parsed) == nil {
		e.Body = parsed
	} else if len(body) > 0 {
		e.Body = string(body)
	}
	return e
}

// retryAfter reads a Retry-After header given in seconds.
func retryAfter(h http.Header, fallback time.Duration) time.Duration {
	if v := strings.TrimSpace(h.Get("Retry-After")); v != "" {
		if secs, err := strconv.ParseFloat(v, 64); err == nil && secs > 0 {
			return time.Duration(secs * float64(time.Second))
		}
		if t, err := http.ParseTime(v); err == nil {
			if d := time.Until(t); d > 0 {
				return d
			}
		}
	}
	return fallback
}

// backoff returns an exponential backoff with full jitter for the given attempt.
func backoff(attempt int) time.Duration {
	ceiling := backoffBase << min(attempt, 10)
	if ceiling > backoffCap || ceiling <= 0 {
		ceiling = backoffCap
	}
	return time.Duration(rand.Int64N(int64(ceiling) + 1))
}

func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// Typed helpers used by the resource services.

func getJSON[T any](ctx context.Context, c *Client, path string, opts []Option) (*T, error) {
	return sendJSON[T](ctx, c, http.MethodGet, path, nil, opts)
}

func sendJSON[T any](ctx context.Context, c *Client, method, path string, body any, opts []Option) (*T, error) {
	resp, err := c.do(ctx, request{method: method, path: path, body: body}, opts)
	if err != nil {
		return nil, err
	}
	out := new(T)
	if err := resp.decodeInto(out); err != nil {
		return nil, err
	}
	return out, nil
}

func sendRaw(ctx context.Context, c *Client, method, path string, body any, opts []Option) (*RawResponse, error) {
	resp, err := c.do(ctx, request{method: method, path: path, body: body}, opts)
	if err != nil {
		return nil, err
	}
	return &RawResponse{StatusCode: resp.statusCode, ContentType: resp.contentType, Body: resp.body}, nil
}
