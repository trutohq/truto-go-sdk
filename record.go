package truto

import "encoding/json"

// Params holds query parameters. Values may be strings, numbers, booleans,
// nil, slices, maps, or any value that encodes to JSON.
type Params map[string]any

// Record is a JSON object whose shape depends on the integration or resource:
// unified, proxy and custom API results, and resources the API spec does not
// describe. Use Decode to turn it into your own type.
type Record map[string]any

// Decode converts v (usually a Record) into T by round-tripping it through JSON.
//
//	type Contact struct{ ID string `json:"id"`; FirstName string `json:"first_name"` }
//	c, err := truto.Decode[Contact](record)
func Decode[T any](v any) (T, error) {
	var out T
	b, err := json.Marshal(v)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(b, &out)
	return out, err
}

// DecodeAll converts each element of records into T.
func DecodeAll[T any](records []Record) ([]T, error) {
	out := make([]T, 0, len(records))
	for _, r := range records {
		v, err := Decode[T](r)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

// Ptr returns a pointer to v. Optional fields in request types are pointers:
//
//	truto.TenantUpdateRequest{Name: truto.Ptr("Acme")}
func Ptr[T any](v T) *T { return &v }

// RawResponse is a response returned as bytes, for endpoints that may answer
// with something other than JSON (files, CSV, HTML).
type RawResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}
