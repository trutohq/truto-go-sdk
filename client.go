// Package truto is a Go client for the Truto API (https://truto.one).
//
// Create a client with an API token and call the API through its services:
//
//	client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))
//	accounts, err := client.IntegratedAccounts.List(ctx,
//		truto.WithQuery(truto.Params{"tenant_id": "acme"})).All()
//
// Admin resources (integrated accounts, sync jobs, webhooks, ...) use the typed
// models generated from Truto's OpenAPI spec. The Unified, Proxy and Custom
// APIs return Record values, since their shape depends on the integration; use
// Decode to turn them into your own types.
//
// Every method accepts Options that override the client settings for that call,
// for example WithQuery, WithTimeout, WithNoRetry or WithHeader.
package truto

import (
	"context"
	"net/url"
	"strings"
)

// Version is the SDK version, sent in the User-Agent header.
const Version = "0.1.0"

const userAgent = "truto-go-sdk/" + Version

// Client is a Truto API client. It is safe for concurrent use.
type Client struct {
	token    string
	settings settings

	APITokens                              *APITokensService
	Environments                           *EnvironmentsService
	Categories                             *CategoriesService
	Teams                                  *TeamsService
	Users                                  *UsersService
	Tenants                                *TenantsService
	LinkTokens                             *LinkTokensService
	SessionTokens                          *SessionTokensService
	IntegratedAccounts                     *IntegratedAccountsService
	SandboxIntegratedAccounts              *SandboxIntegratedAccountsService
	Integrations                           *IntegrationsService
	EnvironmentIntegrations                *EnvironmentIntegrationsService
	UnifiedModels                          *UnifiedModelsService
	EnvironmentUnifiedModels               *EnvironmentUnifiedModelsService
	UnifiedModelResourceMethods            *UnifiedModelResourceMethodsService
	EnvironmentUnifiedModelResourceMethods *EnvironmentUnifiedModelResourceMethodsService
	SyncJobs                               *SyncJobsService
	SyncJobRuns                            *SyncJobRunsService
	SyncJobCronTriggers                    *SyncJobCronTriggersService
	SyncJobTemplates                       *SyncJobTemplatesService
	SyncJobRunStates                       *SyncJobRunStatesService
	Daemons                                *RecordService
	DaemonJobs                             *RecordService
	DaemonJobRuns                          *RecordService
	DaemonJobCronTriggers                  *DaemonJobCronTriggersService
	Workflows                              *WorkflowsService
	WorkflowRuns                           *WorkflowRunsService
	Webhooks                               *WebhooksService
	Datastores                             *DatastoresService
	Documentation                          *DocumentationService
	Logs                                   *LogsService
	StaticGates                            *StaticGatesService
	ObservabilityDestinations              *ObservabilityDestinationsService
	Analytics                              *AnalyticsService
	BatchRequests                          *BatchRequestsService
	EnvironmentIntegrationWebhooks         *WebhookProcessorService
	IntegratedAccountWebhooks              *WebhookProcessorService

	Unified *UnifiedAPIService
	Proxy   *ProxyAPIService
	Custom  *CustomAPIService
	MCP     *MCPService
}

// NewClient returns a client authenticated with an API token.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{token: token, settings: defaultSettings()}
	for _, o := range opts {
		o(&c.settings)
	}

	c.APITokens = &APITokensService{c}
	c.Environments = &EnvironmentsService{c}
	c.Categories = &CategoriesService{c}
	c.Teams = &TeamsService{c}
	c.Users = &UsersService{c}
	c.Tenants = &TenantsService{c}
	c.LinkTokens = &LinkTokensService{c}
	c.SessionTokens = &SessionTokensService{c}
	c.IntegratedAccounts = &IntegratedAccountsService{c}
	c.SandboxIntegratedAccounts = &SandboxIntegratedAccountsService{c}
	c.Integrations = &IntegrationsService{c}
	c.EnvironmentIntegrations = &EnvironmentIntegrationsService{c}
	c.UnifiedModels = &UnifiedModelsService{c}
	c.EnvironmentUnifiedModels = &EnvironmentUnifiedModelsService{c}
	c.UnifiedModelResourceMethods = &UnifiedModelResourceMethodsService{c}
	c.EnvironmentUnifiedModelResourceMethods = &EnvironmentUnifiedModelResourceMethodsService{c}
	c.SyncJobs = &SyncJobsService{c}
	c.SyncJobRuns = &SyncJobRunsService{c}
	c.SyncJobCronTriggers = &SyncJobCronTriggersService{c}
	c.SyncJobTemplates = &SyncJobTemplatesService{c}
	c.SyncJobRunStates = &SyncJobRunStatesService{c}
	c.Daemons = &RecordService{c, "daemon"}
	c.DaemonJobs = &RecordService{c, "daemon-job"}
	c.DaemonJobRuns = &RecordService{c, "daemon-job-run"}
	c.DaemonJobCronTriggers = &DaemonJobCronTriggersService{RecordService{c, "daemon-job-cron-trigger"}}
	c.Workflows = &WorkflowsService{c}
	c.WorkflowRuns = &WorkflowRunsService{c}
	c.Webhooks = &WebhooksService{c}
	c.Datastores = &DatastoresService{c}
	c.Documentation = &DocumentationService{c}
	c.Logs = &LogsService{c}
	c.StaticGates = &StaticGatesService{c}
	c.ObservabilityDestinations = &ObservabilityDestinationsService{c}
	c.Analytics = &AnalyticsService{c}
	c.BatchRequests = &BatchRequestsService{c}
	c.EnvironmentIntegrationWebhooks = &WebhookProcessorService{c, "environment-integration-webhook"}
	c.IntegratedAccountWebhooks = &WebhookProcessorService{c, "integrated-account-webhook"}
	c.Unified = &UnifiedAPIService{c}
	c.Proxy = &ProxyAPIService{c}
	c.Custom = &CustomAPIService{c}
	c.MCP = &MCPService{c}
	return c
}

// Ping checks that the API is reachable and the token is valid.
func (c *Client) Ping(ctx context.Context, opts ...Option) (*Ping, error) {
	return getJSON[Ping](ctx, c, "ping", opts)
}

// Do sends a request to any API path and decodes the JSON response into out
// (which may be nil). Use it for endpoints the SDK does not wrap yet.
func (c *Client) Do(ctx context.Context, method, path string, body, out any, opts ...Option) error {
	resp, err := c.do(ctx, request{method: method, path: strings.TrimLeft(path, "/"), body: body}, opts)
	if err != nil {
		return err
	}
	return resp.decodeInto(out)
}

// p joins path segments, escaping each one.
func p(segments ...string) string {
	out := make([]string, len(segments))
	for i, s := range segments {
		out[i] = url.PathEscape(s)
	}
	return strings.Join(out, "/")
}
