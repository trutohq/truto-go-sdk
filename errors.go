package truto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrTimeout is matched (with errors.Is) by the error returned when a request
// attempt runs longer than the configured timeout.
var ErrTimeout = errors.New("truto: request timed out")

// Done is returned by Iterator.Next when there are no more items.
var Done = errors.New("truto: no more items")

// Error is returned when the Truto API answers with a non-2xx status.
type Error struct {
	// StatusCode is the HTTP status code, e.g. 404.
	StatusCode int
	// Status is the HTTP status text, e.g. "404 Not Found".
	Status string
	// Body is the parsed JSON body, or the body as a string when it is not JSON.
	Body any
	// RawBody is the response body as received.
	RawBody []byte
	// Header holds the response headers.
	Header http.Header
	// Method and Path identify the request that failed.
	Method string
	Path   string
}

func (e *Error) Error() string {
	msg := e.Message()
	if msg == "" {
		return fmt.Sprintf("truto: %s %s: %s", e.Method, e.Path, e.Status)
	}
	return fmt.Sprintf("truto: %s %s: %s: %s", e.Method, e.Path, e.Status, msg)
}

// Message returns the "error" or "message" field of a JSON error body, if any.
func (e *Error) Message() string {
	if m, ok := e.Body.(map[string]any); ok {
		for _, k := range []string{"error", "message"} {
			if s, ok := m[k].(string); ok && s != "" {
				return s
			}
		}
	}
	if s, ok := e.Body.(string); ok && len(s) <= 200 {
		return s
	}
	return ""
}

// DecodeBody decodes the JSON error body into v.
func (e *Error) DecodeBody(v any) error {
	return json.Unmarshal(e.RawBody, v)
}

// IsNotFound reports whether err is a Truto API 404.
func IsNotFound(err error) bool { return hasStatus(err, http.StatusNotFound) }

// IsUnauthorized reports whether err is a Truto API 401.
func IsUnauthorized(err error) bool { return hasStatus(err, http.StatusUnauthorized) }

// IsRateLimited reports whether err is a Truto API 429 (after retries ran out).
func IsRateLimited(err error) bool { return hasStatus(err, http.StatusTooManyRequests) }

func hasStatus(err error, code int) bool {
	var e *Error
	return errors.As(err, &e) && e.StatusCode == code
}
