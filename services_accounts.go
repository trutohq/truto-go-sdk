package truto

import (
	"context"
	"net/http"
	"strings"
)

// BulkCreate creates or upserts many tenants in one call.
func (s *TenantsService) BulkCreate(ctx context.Context, body TenantBulkRequest, opts ...Option) (*TenantBulkCreateResponse, error) {
	return sendJSON[TenantBulkCreateResponse](ctx, s.c, http.MethodPost, "tenant/bulk", body, opts)
}

// LinkTokensService creates link tokens, which start the Truto Link connection
// flow in your frontend.
type LinkTokensService struct{ c *Client }

// Create creates a link token.
func (s *LinkTokensService) Create(ctx context.Context, body LinkTokenRequest, opts ...Option) (*LinkToken, error) {
	return sendJSON[LinkToken](ctx, s.c, http.MethodPost, "link-token", body, opts)
}

// SessionTokensService manages short-lived session tokens scoped to a tenant.
type SessionTokensService struct{ c *Client }

// Create mints a session token.
func (s *SessionTokensService) Create(ctx context.Context, body SessionTokenRequest, opts ...Option) (*SessionTokenResponse, error) {
	return sendJSON[SessionTokenResponse](ctx, s.c, http.MethodPost, "session-token", body, opts)
}

// Revoke revokes the session tokens that match body.
func (s *SessionTokensService) Revoke(ctx context.Context, body SessionTokenDeleteRequest, opts ...Option) (*SessionTokenDeleteResponse, error) {
	return sendJSON[SessionTokenDeleteResponse](ctx, s.c, http.MethodDelete, "session-token", body, opts)
}

// Renew renews the session token the client is authenticated with.
func (s *SessionTokensService) Renew(ctx context.Context, opts ...Option) (*SessionTokenRenewResponse, error) {
	return sendJSON[SessionTokenRenewResponse](ctx, s.c, http.MethodPost, "session-token/renew", nil, opts)
}

// Me describes the credential the client is authenticated with.
func (s *SessionTokensService) Me(ctx context.Context, opts ...Option) (*SessionTokenMeResponse, error) {
	return getJSON[SessionTokenMeResponse](ctx, s.c, "session-token/me", opts)
}

// IntegratedAccountsService manages integrated accounts: your customers'
// connected accounts in third-party apps.
type IntegratedAccountsService struct{ c *Client }

// List returns an iterator over integrated accounts. Filter with WithQuery,
// for example truto.Params{"tenant_id": "acme"}.
func (s *IntegratedAccountsService) List(ctx context.Context, opts ...Option) *Iterator[IntegratedAccountResponse] {
	return newIterator[IntegratedAccountResponse](ctx, s.c, "integrated-account", opts)
}

// Get fetches an integrated account.
func (s *IntegratedAccountsService) Get(ctx context.Context, id string, opts ...Option) (*IntegratedAccountResponse, error) {
	return getJSON[IntegratedAccountResponse](ctx, s.c, p("integrated-account", id), opts)
}

// Create creates (imports) an integrated account from existing credentials.
func (s *IntegratedAccountsService) Create(ctx context.Context, body IntegratedAccountCreate, opts ...Option) (*IntegratedAccountResponse, error) {
	return sendJSON[IntegratedAccountResponse](ctx, s.c, http.MethodPost, "integrated-account", body, opts)
}

// Update changes an integrated account.
func (s *IntegratedAccountsService) Update(ctx context.Context, id string, body IntegratedAccountUpdate, opts ...Option) (*IntegratedAccountResponse, error) {
	return sendJSON[IntegratedAccountResponse](ctx, s.c, http.MethodPatch, p("integrated-account", id), body, opts)
}

// Delete deletes an integrated account.
func (s *IntegratedAccountsService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("integrated-account", id), nil, opts)
}

// UpdateCredentials replaces the credentials stored in the account's context.
func (s *IntegratedAccountsService) UpdateCredentials(ctx context.Context, id string, body IntegratedAccountCredentialsUpdateRequest, opts ...Option) (*IntegratedAccountResponse, error) {
	return sendJSON[IntegratedAccountResponse](ctx, s.c, http.MethodPatch, p("integrated-account", id, "credentials"), body, opts)
}

// BulkDelete deletes many integrated accounts. See the API reference for the
// accepted body shapes.
func (s *IntegratedAccountsService) BulkDelete(ctx context.Context, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, "integrated-account/bulk-delete", body, opts)
}

// CreateToken creates an integrated-account token, which can call the Unified
// and Proxy APIs for this one account only.
func (s *IntegratedAccountsService) CreateToken(ctx context.Context, id string, opts ...Option) (*IntegratedAccountTokenResponse, error) {
	body := IntegratedAccountTokenRequest{IntegratedAccountID: &id}
	return sendJSON[IntegratedAccountTokenResponse](ctx, s.c, http.MethodPost, "integrated-account/token", body, opts)
}

// CreateFormToken creates a short-lived token for the account's RapidForm.
func (s *IntegratedAccountsService) CreateFormToken(ctx context.Context, id string, opts ...Option) (*IntegratedAccountFormTokenResponse, error) {
	body := IntegratedAccountFormTokenRequest{IntegratedAccountID: id}
	return sendJSON[IntegratedAccountFormTokenResponse](ctx, s.c, http.MethodPost, "integrated-account/form-token", body, opts)
}

// RefreshCredentials refreshes the account's OAuth credentials now.
func (s *IntegratedAccountsService) RefreshCredentials(ctx context.Context, id string, opts ...Option) (*IntegratedAccountRefreshCredentialsResponse, error) {
	body := IntegratedAccountRefreshCredentialsRequest{ID: id}
	return sendJSON[IntegratedAccountRefreshCredentialsResponse](ctx, s.c, http.MethodPost, "integrated-account/refresh-credentials", body, opts)
}

// RunPostInstallActions runs the integration's post-install actions again.
func (s *IntegratedAccountsService) RunPostInstallActions(ctx context.Context, id string, opts ...Option) (*IntegratedAccountRunPostInstallActionsResponse, error) {
	body := IntegratedAccountRunPostInstallActionsRequest{ID: id}
	return sendJSON[IntegratedAccountRunPostInstallActionsResponse](ctx, s.c, http.MethodPost, "integrated-account/run-post-install-actions", body, opts)
}

// RefreshMappingCache clears cached unified mappings.
func (s *IntegratedAccountsService) RefreshMappingCache(ctx context.Context, body IntegratedAccountRefreshMappingCacheRequest, opts ...Option) (*IntegratedAccountRefreshMappingCacheResponse, error) {
	return sendJSON[IntegratedAccountRefreshMappingCacheResponse](ctx, s.c, http.MethodPost, "integrated-account/refresh-mapping-cache", body, opts)
}

// Validate runs the integration's validation steps against the account.
func (s *IntegratedAccountsService) Validate(ctx context.Context, id string, opts ...Option) (*IntegratedAccountValidateResponse, error) {
	return sendJSON[IntegratedAccountValidateResponse](ctx, s.c, http.MethodPost, p("integrated-account", id, "validate"), nil, opts)
}

// ToolsFilter narrows the tools returned by IntegratedAccountsService.Tools.
type ToolsFilter struct {
	// Methods, e.g. "list", "get", "create". Empty means all.
	Methods []string
	// Tags set on resource methods. Empty means all.
	Tags []string
}

// Tools lists the tools (proxy API methods described for LLM use) available
// for the account.
func (s *IntegratedAccountsService) Tools(ctx context.Context, id string, filter ToolsFilter, opts ...Option) (*ToolsResponseSchema, error) {
	q := Params{}
	if len(filter.Methods) > 0 {
		q["methods"] = strings.Join(filter.Methods, ",")
	}
	if len(filter.Tags) > 0 {
		q["tags"] = strings.Join(filter.Tags, ",")
	}
	if len(q) > 0 {
		opts = append([]Option{WithQuery(q)}, opts...)
	}
	return getJSON[ToolsResponseSchema](ctx, s.c, p("integrated-account", id, "tools"), opts)
}

// Capabilities describes what the account can do: resources, methods, auth.
func (s *IntegratedAccountsService) Capabilities(ctx context.Context, id string, opts ...Option) (*CapabilitiesResponse, error) {
	return getJSON[CapabilitiesResponse](ctx, s.c, p("integrated-account", id, "capabilities"), opts)
}

// ListMCPServers returns an iterator over the account's MCP servers.
func (s *IntegratedAccountsService) ListMCPServers(ctx context.Context, id string, opts ...Option) *Iterator[MCP] {
	return newIterator[MCP](ctx, s.c, p("integrated-account", id, "mcp"), opts)
}

// GetMCPServer fetches one MCP server of the account.
func (s *IntegratedAccountsService) GetMCPServer(ctx context.Context, id, mcpID string, opts ...Option) (*MCP, error) {
	return getJSON[MCP](ctx, s.c, p("integrated-account", id, "mcp", mcpID), opts)
}

// CreateMCPServer creates an MCP server for the account.
func (s *IntegratedAccountsService) CreateMCPServer(ctx context.Context, id string, body IntegratedAccountMcpCreateRequest, opts ...Option) (*MCPCreateResponse, error) {
	return sendJSON[MCPCreateResponse](ctx, s.c, http.MethodPost, p("integrated-account", id, "mcp"), body, opts)
}

// UpdateMCPServer changes an MCP server of the account.
func (s *IntegratedAccountsService) UpdateMCPServer(ctx context.Context, id, mcpID string, body IntegratedAccountMcpUpdateRequest, opts ...Option) (*MCP, error) {
	return sendJSON[MCP](ctx, s.c, http.MethodPatch, p("integrated-account", id, "mcp", mcpID), body, opts)
}

// DeleteMCPServer deletes an MCP server of the account.
func (s *IntegratedAccountsService) DeleteMCPServer(ctx context.Context, id, mcpID string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("integrated-account", id, "mcp", mcpID), nil, opts)
}
