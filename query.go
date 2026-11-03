package truto

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// encodeQuery encodes params the way the qs library does (the Truto API parses
// query strings with qs): keys sorted, nested objects as a[b]=v, arrays with
// indices as a[0]=v, empty arrays and objects left out, nil as "a=".
//
// Values are first normalised through JSON, so structs with json tags work too.
func encodeQuery(params Params) (string, error) {
	if len(params) == 0 {
		return "", nil
	}
	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v map[string]any
	if err := dec.Decode(&v); err != nil {
		return "", err
	}
	var parts []string
	appendQuery(&parts, "", v)
	return strings.Join(parts, "&"), nil
}

func appendQuery(parts *[]string, prefix string, v any) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			key := k
			if prefix != "" {
				key = prefix + "[" + k + "]"
			}
			appendQuery(parts, key, t[k])
		}
	case []any:
		for i, item := range t {
			appendQuery(parts, prefix+"["+strconv.Itoa(i)+"]", item)
		}
	case nil:
		*parts = append(*parts, escape(prefix)+"=")
	case string:
		*parts = append(*parts, escape(prefix)+"="+escape(t))
	case json.Number:
		*parts = append(*parts, escape(prefix)+"="+escape(t.String()))
	case bool:
		s := "false"
		if t {
			s = "true"
		}
		*parts = append(*parts, escape(prefix)+"="+s)
	}
}

const upperhex = "0123456789ABCDEF"

// escape percent-encodes everything outside the RFC 3986 unreserved set,
// matching qs (and encodeURIComponent, which also encodes !'()*).
func escape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&15])
	}
	return b.String()
}
