package truto

import (
	"context"
	"net/http"
	"strings"
)

// UnifiedResource identifies a unified API resource for one integrated account.
type UnifiedResource struct {
	// Model is the unified model, e.g. "crm".
	Model string
	// Resource is the resource in that model, e.g. "contacts".
	Resource string
	// IntegratedAccountID is the account to call. Sent as a query parameter.
	IntegratedAccountID string
}

func (r UnifiedResource) path(extra ...string) string {
	return p(append([]string{"unified", r.Model, r.Resource}, extra...)...)
}

func (r UnifiedResource) opts(opts []Option) []Option {
	return append([]Option{WithQuery(Params{"integrated_account_id": r.IntegratedAccountID})}, opts...)
}

// UnifiedAPIService calls the Unified API: one schema per category (CRM, ATS,
// HRIS, ...) across every integration in it. Results are Records; use Decode
// for your own types.
type UnifiedAPIService struct{ c *Client }

// List returns an iterator over the resource.
func (s *UnifiedAPIService) List(ctx context.Context, r UnifiedResource, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, r.path(), r.opts(opts))
}

// Get fetches one record.
func (s *UnifiedAPIService) Get(ctx context.Context, r UnifiedResource, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, r.path(id), r.opts(opts))
}

// Create creates a record.
func (s *UnifiedAPIService) Create(ctx context.Context, r UnifiedResource, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, r.path(), body, r.opts(opts))
}

// Update changes a record.
func (s *UnifiedAPIService) Update(ctx context.Context, r UnifiedResource, id string, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPatch, r.path(id), body, r.opts(opts))
}

// Delete deletes a record.
func (s *UnifiedAPIService) Delete(ctx context.Context, r UnifiedResource, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodDelete, r.path(id), nil, r.opts(opts))
}

// CustomMethod calls a non-CRUD method of the resource (POST
// /unified/{model}/{resource}/{method}). body may be nil.
func (s *UnifiedAPIService) CustomMethod(ctx context.Context, r UnifiedResource, method string, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, r.path(method), body, r.opts(opts))
}

// Meta describes one method of a unified resource for an integration: the
// query and body it accepts and its response schema.
func (s *UnifiedAPIService) Meta(ctx context.Context, model, resource, integrationName, method string, opts ...Option) (*UnifiedApiMetaResponse, error) {
	return getJSON[UnifiedApiMetaResponse](ctx, s.c, p("unified", model, resource, integrationName, "meta", method), opts)
}

// ModelMeta describes a unified model for an integration: plan, scopes and
// instructions.
func (s *UnifiedAPIService) ModelMeta(ctx context.Context, model, integrationName string, opts ...Option) (*UnifiedMetaGetResponse, error) {
	return getJSON[UnifiedMetaGetResponse](ctx, s.c, p("unified", "meta", model, integrationName), opts)
}

// ProxyResource identifies a proxy API resource for one integrated account.
type ProxyResource struct {
	// Resource is the integration's own resource name, e.g. "contacts".
	Resource string
	// IntegratedAccountID is the account to call. Sent as a query parameter.
	IntegratedAccountID string
}

func (r ProxyResource) path(extra ...string) string {
	return p(append([]string{"proxy", r.Resource}, extra...)...)
}

func (r ProxyResource) opts(opts []Option) []Option {
	return append([]Option{WithQuery(Params{"integrated_account_id": r.IntegratedAccountID})}, opts...)
}

// ProxyAPIService calls the Proxy API: the integration's own resources and
// fields, with Truto handling auth and pagination.
type ProxyAPIService struct{ c *Client }

// List returns an iterator over the resource.
func (s *ProxyAPIService) List(ctx context.Context, r ProxyResource, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, r.path(), r.opts(opts))
}

// Get fetches one record.
func (s *ProxyAPIService) Get(ctx context.Context, r ProxyResource, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, r.path(id), r.opts(opts))
}

// Create creates a record.
func (s *ProxyAPIService) Create(ctx context.Context, r ProxyResource, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, r.path(), body, r.opts(opts))
}

// Update changes a record.
func (s *ProxyAPIService) Update(ctx context.Context, r ProxyResource, id string, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPatch, r.path(id), body, r.opts(opts))
}

// Delete deletes a record.
func (s *ProxyAPIService) Delete(ctx context.Context, r ProxyResource, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodDelete, r.path(id), nil, r.opts(opts))
}

// CustomMethod calls a non-CRUD method of the resource (POST
// /proxy/{resource}/{method}). body may be nil.
func (s *ProxyAPIService) CustomMethod(ctx context.Context, r ProxyResource, method string, body any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, r.path(method), body, r.opts(opts))
}

// CustomAPIService calls Custom API endpoints you defined on an integration.
// Responses are returned raw, because a custom endpoint may answer with
// anything; RawResponse.Body is JSON when ContentType says so.
type CustomAPIService struct{ c *Client }

// customPath escapes each segment of a caller-provided path, keeping slashes.
func customPath(path string) string {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	return "custom/" + p(segs...)
}

func customOpts(integratedAccountID string, opts []Option) []Option {
	return append([]Option{WithQuery(Params{"integrated_account_id": integratedAccountID})}, opts...)
}

// Get sends GET /custom/{path}.
func (s *CustomAPIService) Get(ctx context.Context, integratedAccountID, path string, opts ...Option) (*RawResponse, error) {
	return sendRaw(ctx, s.c, http.MethodGet, customPath(path), nil, customOpts(integratedAccountID, opts))
}

// Post sends POST /custom/{path}.
func (s *CustomAPIService) Post(ctx context.Context, integratedAccountID, path string, body any, opts ...Option) (*RawResponse, error) {
	return sendRaw(ctx, s.c, http.MethodPost, customPath(path), body, customOpts(integratedAccountID, opts))
}

// Put sends PUT /custom/{path}.
func (s *CustomAPIService) Put(ctx context.Context, integratedAccountID, path string, body any, opts ...Option) (*RawResponse, error) {
	return sendRaw(ctx, s.c, http.MethodPut, customPath(path), body, customOpts(integratedAccountID, opts))
}

// Patch sends PATCH /custom/{path}.
func (s *CustomAPIService) Patch(ctx context.Context, integratedAccountID, path string, body any, opts ...Option) (*RawResponse, error) {
	return sendRaw(ctx, s.c, http.MethodPatch, customPath(path), body, customOpts(integratedAccountID, opts))
}

// Delete sends DELETE /custom/{path}.
func (s *CustomAPIService) Delete(ctx context.Context, integratedAccountID, path string, opts ...Option) (*RawResponse, error) {
	return sendRaw(ctx, s.c, http.MethodDelete, customPath(path), nil, customOpts(integratedAccountID, opts))
}

// MCPService talks JSON-RPC to a Truto MCP server, addressed by its token
// (see IntegratedAccountsService.CreateMCPServer).
type MCPService struct{ c *Client }

// MCPProtocolVersion is the protocol version sent by Initialize.
const MCPProtocolVersion = "2024-11-05"

// Call sends a raw JSON-RPC payload.
func (s *MCPService) Call(ctx context.Context, token string, payload any, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, p("mcp", token), payload, opts)
}

// Initialize starts an MCP session.
func (s *MCPService) Initialize(ctx context.Context, token, clientName, clientVersion string, opts ...Option) (*Record, error) {
	return s.Call(ctx, token, Record{
		"jsonrpc": "2.0",
		"method":  "initialize",
		"id":      0,
		"params": Record{
			"protocolVersion": MCPProtocolVersion,
			"clientInfo":      Record{"name": clientName, "version": clientVersion},
		},
	}, opts...)
}

// ListTools lists the server's tools.
func (s *MCPService) ListTools(ctx context.Context, token string, opts ...Option) (*Record, error) {
	return s.Call(ctx, token, Record{"jsonrpc": "2.0", "method": "tools/list", "id": 1}, opts...)
}

// CallTool calls one tool with arguments.
func (s *MCPService) CallTool(ctx context.Context, token, name string, args any, opts ...Option) (*Record, error) {
	if args == nil {
		args = Record{}
	}
	return s.Call(ctx, token, Record{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"id":      2,
		"params":  Record{"name": name, "arguments": args},
	}, opts...)
}
