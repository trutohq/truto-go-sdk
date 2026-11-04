package truto_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	truto "github.com/trutohq/truto-go-sdk"
)

// recorded is one request the mock server received.
type recorded struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

type mock struct {
	t        *testing.T
	srv      *httptest.Server
	mu       sync.Mutex
	requests []recorded
}

// newMock starts a server whose handler gets the call index (0-based).
func newMock(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, call int)) (*mock, *truto.Client) {
	t.Helper()
	m := &mock{t: t}
	m.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		m.mu.Lock()
		call := len(m.requests)
		m.requests = append(m.requests, recorded{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Clone(), string(body)})
		m.mu.Unlock()
		handler(w, r, call)
	}))
	t.Cleanup(m.srv.Close)
	return m, truto.NewClient("test-token", truto.WithBaseURL(m.srv.URL+"/"))
}

func (m *mock) calls() []recorded {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]recorded(nil), m.requests...)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

var ctx = context.Background()

func TestHeadersAndGet(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"id": "ia_1", "tenant_id": "acme"})
	})
	acct, err := c.IntegratedAccounts.Get(ctx, "ia/1", truto.WithQuery(truto.Params{"x": "y"}))
	if err != nil {
		t.Fatal(err)
	}
	if *acct.ID != "ia_1" || *acct.TenantID != "acme" {
		t.Errorf("decoded %+v", acct)
	}
	got := m.calls()[0]
	if got.Method != "GET" || got.Path != "/integrated-account/ia/1" && got.Path != "/integrated-account/ia%2F1" {
		t.Errorf("request %s %s", got.Method, got.Path)
	}
	if got.Query != "x=y" {
		t.Errorf("query %q", got.Query)
	}
	if h := got.Header.Get("Authorization"); h != "Bearer test-token" {
		t.Errorf("auth %q", h)
	}
	if h := got.Header.Get("User-Agent"); h != "truto-go-sdk/"+truto.Version {
		t.Errorf("user agent %q", h)
	}
	if h := got.Header.Get("Accept"); h != "application/json" {
		t.Errorf("accept %q", h)
	}
}

func TestPathSegmentsAreEscaped(t *testing.T) {
	var rawPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath = r.URL.EscapedPath()
		writeJSON(w, 200, map[string]any{})
	}))
	defer srv.Close()
	c := truto.NewClient("t", truto.WithBaseURL(srv.URL))
	if _, err := c.Tenants.Get(ctx, "a/b c"); err != nil {
		t.Fatal(err)
	}
	if rawPath != "/tenant/a%2Fb%20c" {
		t.Errorf("path %q", rawPath)
	}
}

func TestPostJSONBody(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 201, map[string]any{"link_token": "lt_1"})
	})
	lt, err := c.LinkTokens.Create(ctx, truto.LinkTokenRequest{TenantID: truto.Ptr("acme")})
	if err != nil {
		t.Fatal(err)
	}
	if *lt.LinkToken != "lt_1" {
		t.Errorf("link token %v", lt.LinkToken)
	}
	got := m.calls()[0]
	if got.Method != "POST" || got.Path != "/link-token" {
		t.Errorf("request %s %s", got.Method, got.Path)
	}
	if got.Header.Get("Content-Type") != "application/json" {
		t.Errorf("content type %q", got.Header.Get("Content-Type"))
	}
	if strings.TrimSpace(got.Body) != `{"tenant_id":"acme"}` {
		t.Errorf("body %s", got.Body)
	}
}

func TestNoContent(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		w.WriteHeader(204)
	})
	if err := c.Documentation.Delete(ctx, "doc_1"); err != nil {
		t.Fatal(err)
	}
	rec, err := c.Daemons.Delete(ctx, "d_1")
	if err != nil {
		t.Fatal(err)
	}
	if rec == nil || len(*rec) != 0 {
		t.Errorf("want empty record, got %v", rec)
	}
}

func TestCustomAPIRawBody(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("a,b\n1,2\n"))
	})
	res, err := c.Custom.Get(ctx, "ia_1", "/reports/monthly")
	if err != nil {
		t.Fatal(err)
	}
	if res.ContentType != "text/csv" || string(res.Body) != "a,b\n1,2\n" {
		t.Errorf("got %q %q", res.ContentType, res.Body)
	}
	got := m.calls()[0]
	if got.Path != "/custom/reports/monthly" || got.Query != "integrated_account_id=ia_1" {
		t.Errorf("request %s?%s", got.Path, got.Query)
	}
}

func TestAPIError(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 404, map[string]any{"error": "Integrated account not found"})
	})
	_, err := c.IntegratedAccounts.Get(ctx, "missing")
	var apiErr *truto.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *truto.Error, got %T %v", err, err)
	}
	if apiErr.StatusCode != 404 || apiErr.Message() != "Integrated account not found" {
		t.Errorf("got %d %q", apiErr.StatusCode, apiErr.Message())
	}
	if !truto.IsNotFound(err) || truto.IsRateLimited(err) {
		t.Error("status helpers")
	}
	if !strings.Contains(err.Error(), "Integrated account not found") {
		t.Errorf("message %q", err.Error())
	}
}

func TestRetryOn503(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		if call == 0 {
			writeJSON(w, 503, map[string]any{"error": "unavailable"})
			return
		}
		writeJSON(w, 200, map[string]any{"result": []any{}, "next_cursor": nil})
	})
	var retries []truto.RetryInfo
	_, err := c.Tenants.List(ctx, truto.WithOnRetry(func(i truto.RetryInfo) { retries = append(retries, i) })).All()
	if err != nil {
		t.Fatal(err)
	}
	if len(m.calls()) != 2 || len(retries) != 1 || retries[0].StatusCode != 503 || retries[0].RetryCount != 1 {
		t.Errorf("calls=%d retries=%+v", len(m.calls()), retries)
	}
}

func TestRetryPostResendsBody(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		if call == 0 {
			writeJSON(w, 502, nil)
			return
		}
		writeJSON(w, 200, map[string]any{"id": "t1"})
	})
	if _, err := c.Tenants.Create(ctx, truto.TenantCreateRequest{ID: "t1"}); err != nil {
		t.Fatal(err)
	}
	calls := m.calls()
	if len(calls) != 2 || calls[0].Body != calls[1].Body || calls[1].Body == "" {
		t.Errorf("bodies %q %q", calls[0].Body, calls[1].Body)
	}
}

func TestRateLimitRetryAfter(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		if call == 0 {
			w.Header().Set("Retry-After", "1")
			writeJSON(w, 429, map[string]any{"error": "slow down"})
			return
		}
		writeJSON(w, 200, map[string]any{"status": "ok"})
	})
	var limited []truto.RateLimitInfo
	start := time.Now()
	_, err := c.Analytics.Status(ctx, truto.WithOnRateLimited(func(i truto.RateLimitInfo) { limited = append(limited, i) }))
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed < 900*time.Millisecond {
		t.Errorf("did not wait for Retry-After: %s", elapsed)
	}
	if len(m.calls()) != 2 || len(limited) != 1 || limited[0].RetryAfter != time.Second || limited[0].RetryCount != 0 {
		t.Errorf("calls=%d limited=%+v", len(m.calls()), limited)
	}
}

func TestRetriesExhausted(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 503, map[string]any{"error": "down"})
	})
	_, err := c.Ping(ctx, truto.WithMaxRetries(2))
	var apiErr *truto.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 {
		t.Fatalf("got %v", err)
	}
	if n := len(m.calls()); n != 3 {
		t.Errorf("want 3 attempts, got %d", n)
	}
}

func TestNoRetry(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 503, nil)
	})
	if _, err := c.Ping(ctx, truto.WithNoRetry()); err == nil {
		t.Fatal("want error")
	}
	if n := len(m.calls()); n != 1 {
		t.Errorf("want 1 attempt, got %d", n)
	}
}

func TestNonRetryableStatusIsNotRetried(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 400, map[string]any{"error": "bad"})
	})
	if _, err := c.Ping(ctx); err == nil {
		t.Fatal("want error")
	}
	if n := len(m.calls()); n != 1 {
		t.Errorf("want 1 attempt, got %d", n)
	}
}

func TestTimeout(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	})
	_, err := c.Ping(ctx, truto.WithTimeout(100*time.Millisecond))
	if !errors.Is(err, truto.ErrTimeout) {
		t.Fatalf("want ErrTimeout, got %v", err)
	}
	if n := len(m.calls()); n != 1 {
		t.Errorf("timeouts must not retry, got %d attempts", n)
	}
}

func TestPagination(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		switch r.URL.Query().Get("next_cursor") {
		case "":
			writeJSON(w, 200, map[string]any{"result": []any{map[string]any{"id": "1"}, map[string]any{"id": "2"}}, "next_cursor": "c2"})
		case "c2":
			// An empty page that still has a cursor must not end the listing.
			writeJSON(w, 200, map[string]any{"result": []any{}, "next_cursor": "c3"})
		case "c3":
			// Some endpoints are documented with "results".
			writeJSON(w, 200, map[string]any{"results": []any{map[string]any{"id": "3"}}, "next_cursor": nil})
		}
	})
	q := truto.Params{"tenant_id": "acme"}
	accounts, err := c.IntegratedAccounts.List(ctx, truto.WithQuery(q)).All()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, a := range accounts {
		ids = append(ids, *a.ID)
	}
	if strings.Join(ids, ",") != "1,2,3" {
		t.Errorf("ids %v", ids)
	}
	calls := m.calls()
	if len(calls) != 3 {
		t.Fatalf("calls %d", len(calls))
	}
	if calls[0].Query != "tenant_id=acme" || calls[1].Query != "next_cursor=c2&tenant_id=acme" || calls[2].Query != "next_cursor=c3&tenant_id=acme" {
		t.Errorf("queries %q %q %q", calls[0].Query, calls[1].Query, calls[2].Query)
	}
	if len(q) != 1 {
		t.Errorf("caller query was mutated: %v", q)
	}
}

func TestPaginationHelpers(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		if r.URL.Query().Get("next_cursor") == "" {
			writeJSON(w, 200, map[string]any{"result": []any{map[string]any{"id": "a"}}, "next_cursor": "n"})
			return
		}
		writeJSON(w, 200, map[string]any{"result": []any{map[string]any{"id": "b"}}})
	})

	first, err := c.Webhooks.List(ctx).First()
	if err != nil || *first.ID != "a" {
		t.Fatalf("first %v %v", first, err)
	}

	it := c.Webhooks.List(ctx)
	page, err := it.NextPage()
	if err != nil || len(page.Items) != 1 || page.NextCursor != "n" {
		t.Fatalf("page %+v %v", page, err)
	}
	page, err = it.NextPage()
	if err != nil || page.NextCursor != "" || *page.Items[0].ID != "b" {
		t.Fatalf("page 2 %+v %v", page, err)
	}
	if _, err := it.NextPage(); err != truto.Done {
		t.Fatalf("want Done, got %v", err)
	}

	var seen []string
	for w, err := range c.Webhooks.List(ctx).Seq() {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, *w.ID)
	}
	if strings.Join(seen, "") != "ab" {
		t.Errorf("seq %v", seen)
	}
}

// Truto error bodies carry the status text in "error" and the detail in
// "message"; Message must return the detail.
func TestAPIErrorPrefersMessage(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 400, map[string]any{"message": "state_key is required", "statusCode": 400, "error": "Bad Request"})
	})
	_, err := c.SyncJobRunStates.List(ctx).First()
	var apiErr *truto.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *truto.Error, got %T %v", err, err)
	}
	if apiErr.Message() != "state_key is required" {
		t.Errorf("Message() = %q", apiErr.Message())
	}
	if !strings.Contains(err.Error(), "state_key is required") {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestPaginationRejectsNonJSON(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`{"result":[{"id":"t1"}]}`))
	})
	_, err := c.Tenants.List(ctx).NextPage()
	if err == nil || !strings.Contains(err.Error(), "expected a JSON response") {
		t.Fatalf("want a content-type error, got %v", err)
	}
}

func TestPaginationErrorStopsIteration(t *testing.T) {
	_, c := newMock(t, func(w http.ResponseWriter, r *http.Request, call int) {
		writeJSON(w, 401, map[string]any{"error": "bad token"})
	})
	it := c.Tenants.List(ctx)
	if _, err := it.Next(); !truto.IsUnauthorized(err) {
		t.Fatalf("got %v", err)
	}
	if _, err := it.Next(); !truto.IsUnauthorized(err) {
		t.Fatalf("error must stick, got %v", err)
	}
}

// A realistic payload with a field and an enum value the SDK does not know.
const integratedAccountJSON = `{
  "id": "6b7e2a52-ef1c-4bba-9a3c-1f6c2a0d6d0e",
  "tenant_id": "acme",
  "environment_integration_id": "ei_1",
  "status": "paused_by_vendor",
  "is_sandbox": false,
  "context": {"instance_url": "https://acme.example.com"},
  "authentication_method": "oauth2",
  "brand_new_field": {"nested": true},
  "integration": {"id": "int_1", "name": "hubspot", "category": "crm"},
  "created_at": "2026-10-01T10:00:00.000Z"
}`

func TestTypedDecodeIsLenient(t *testing.T) {
	var acct truto.IntegratedAccountResponse
	if err := json.Unmarshal([]byte(integratedAccountJSON), &acct); err != nil {
		t.Fatal(err)
	}
	if *acct.Status != "paused_by_vendor" {
		t.Errorf("unknown enum value lost: %v", *acct.Status)
	}
	if *acct.AuthenticationMethod != truto.IntegratedAccountAuthenticationMethodOAuth2 {
		t.Errorf("auth method %v", *acct.AuthenticationMethod)
	}
	if acct.Integration == nil || *acct.Integration.Name != "hubspot" {
		t.Errorf("integration %+v", acct.Integration)
	}
	if acct.Context["instance_url"] != "https://acme.example.com" {
		t.Errorf("context %v", acct.Context)
	}
	out, err := json.Marshal(truto.IntegratedAccountUpdate{Context: map[string]any{"a": 1}})
	if err != nil || string(out) != `{"context":{"a":1}}` {
		t.Errorf("update body %s %v", out, err)
	}
}

func TestUnionAccessors(t *testing.T) {
	var cfg truto.DatastoreConfig
	if err := json.Unmarshal([]byte(`{"bucket":"b1","region":"eu-west-1"}`), &cfg); err != nil {
		t.Fatal(err)
	}
	s3, err := cfg.AsDatastoreS3Config()
	if err != nil {
		t.Fatal(err)
	}
	if s3.Bucket != "b1" {
		t.Errorf("s3 %+v", s3)
	}
	wrapped, err := truto.DatastoreConfigFromDatastoreS3Config(*s3)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(wrapped)
	if !strings.Contains(string(b), `"bucket":"b1"`) {
		t.Errorf("round trip %s", b)
	}
}

func TestUnifiedList(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{
			"result":      []any{map[string]any{"id": "c1", "first_name": "Ada"}},
			"next_cursor": nil,
		})
	})
	r := truto.UnifiedResource{Model: "crm", Resource: "contacts", IntegratedAccountID: "ia_1"}
	records, err := c.Unified.List(ctx, r, truto.WithQuery(truto.Params{"limit": 10})).All()
	if err != nil {
		t.Fatal(err)
	}
	type contact struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
	}
	contacts, err := truto.DecodeAll[contact](records)
	if err != nil || len(contacts) != 1 || contacts[0].FirstName != "Ada" {
		t.Fatalf("contacts %+v %v", contacts, err)
	}
	got := m.calls()[0]
	if got.Path != "/unified/crm/contacts" || got.Query != "integrated_account_id=ia_1&limit=10" {
		t.Errorf("request %s?%s", got.Path, got.Query)
	}
}

func TestProxyCustomMethod(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"ok": true})
	})
	r := truto.ProxyResource{Resource: "deals", IntegratedAccountID: "ia_9"}
	if _, err := c.Proxy.CustomMethod(ctx, r, "merge", map[string]any{"ids": []string{"1", "2"}}); err != nil {
		t.Fatal(err)
	}
	got := m.calls()[0]
	if got.Method != "POST" || got.Path != "/proxy/deals/merge" || got.Query != "integrated_account_id=ia_9" {
		t.Errorf("request %s %s?%s", got.Method, got.Path, got.Query)
	}
}

func TestMCPPayloads(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "result": map[string]any{}})
	})
	if _, err := c.MCP.Initialize(ctx, "tok", "my-agent", "1.0"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.MCP.ListTools(ctx, "tok"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.MCP.CallTool(ctx, "tok", "list_contacts", map[string]any{"limit": 5}); err != nil {
		t.Fatal(err)
	}
	calls := m.calls()
	want := []string{
		`{"id":0,"jsonrpc":"2.0","method":"initialize","params":{"clientInfo":{"name":"my-agent","version":"1.0"},"protocolVersion":"2024-11-05"}}`,
		`{"id":1,"jsonrpc":"2.0","method":"tools/list"}`,
		`{"id":2,"jsonrpc":"2.0","method":"tools/call","params":{"arguments":{"limit":5},"name":"list_contacts"}}`,
	}
	for i, w := range want {
		if calls[i].Path != "/mcp/tok" || strings.TrimSpace(calls[i].Body) != w {
			t.Errorf("call %d: %s %s", i, calls[i].Path, calls[i].Body)
		}
	}
}

func TestIntegrationUninstall(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		if r.Method == "GET" {
			writeJSON(w, 200, map[string]any{"result": []any{map[string]any{"id": "ei_1"}, map[string]any{"id": "ei_2"}}})
			return
		}
		writeJSON(w, 200, map[string]any{"id": strings.TrimPrefix(r.URL.Path, "/environment-integration/")})
	})
	res, err := c.Integrations.Uninstall(ctx, "int_1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("results %+v", res)
	}
	calls := m.calls()
	if calls[0].Path != "/environment-integration" || calls[0].Query != "integration_id=int_1" ||
		calls[1].Method != "DELETE" || calls[1].Path != "/environment-integration/ei_1" ||
		calls[2].Path != "/environment-integration/ei_2" {
		t.Errorf("calls %+v", calls)
	}
}

func TestIntegratedAccountTools(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"result": []any{}})
	})
	_, err := c.IntegratedAccounts.Tools(ctx, "ia_1", truto.ToolsFilter{Methods: []string{"list", "get"}})
	if err != nil {
		t.Fatal(err)
	}
	if got := m.calls()[0]; got.Path != "/integrated-account/ia_1/tools" || got.Query != "methods%5B0%5D=list&methods%5B1%5D=get" {
		t.Errorf("request %s?%s", got.Path, got.Query)
	}
}

func TestClientDo(t *testing.T) {
	m, c := newMock(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"value": 42})
	})
	var out struct{ Value int }
	if err := c.Do(ctx, "GET", "/schema/x", nil, &out); err != nil || out.Value != 42 {
		t.Fatalf("%v %+v", err, out)
	}
	if m.calls()[0].Path != "/schema/x" {
		t.Errorf("path %s", m.calls()[0].Path)
	}
}
