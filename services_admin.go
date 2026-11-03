package truto

import (
	"context"
	"net/http"
)

// APITokensService manages API tokens of the environment.
type APITokensService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *APITokensService) List(ctx context.Context, opts ...Option) *Iterator[ApiTokenResponse] {
	return newIterator[ApiTokenResponse](ctx, s.c, "api-token", opts)
}

// Get fetches one item by ID.
func (s *APITokensService) Get(ctx context.Context, id string, opts ...Option) (*ApiTokenResponse, error) {
	return getJSON[ApiTokenResponse](ctx, s.c, p("api-token", id), opts)
}

// EnvironmentsService manages environments of the team.
type EnvironmentsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *EnvironmentsService) List(ctx context.Context, opts ...Option) *Iterator[EnvironmentResponse] {
	return newIterator[EnvironmentResponse](ctx, s.c, "environment", opts)
}

// Get fetches one item by ID.
func (s *EnvironmentsService) Get(ctx context.Context, id string, opts ...Option) (*EnvironmentResponse, error) {
	return getJSON[EnvironmentResponse](ctx, s.c, p("environment", id), opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *EnvironmentsService) Update(ctx context.Context, id string, body Record, opts ...Option) (*EnvironmentResponse, error) {
	return sendJSON[EnvironmentResponse](ctx, s.c, http.MethodPatch, p("environment", id), body, opts)
}

// CategoriesService manages integration categories.
type CategoriesService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *CategoriesService) List(ctx context.Context, opts ...Option) *Iterator[CategorySchema] {
	return newIterator[CategorySchema](ctx, s.c, "category", opts)
}

// TeamsService manages teams the token can see.
type TeamsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *TeamsService) List(ctx context.Context, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, "team", opts)
}

// Get fetches one item by ID.
func (s *TeamsService) Get(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, p("team", id), opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *TeamsService) Update(ctx context.Context, id string, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPatch, p("team", id), body, opts)
}

// UsersService manages users of the team.
type UsersService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *UsersService) List(ctx context.Context, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, "user", opts)
}

// Get fetches one item by ID.
func (s *UsersService) Get(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, p("user", id), opts)
}

// TenantsService manages tenants (your customers).
type TenantsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *TenantsService) List(ctx context.Context, opts ...Option) *Iterator[TenantResponse] {
	return newIterator[TenantResponse](ctx, s.c, "tenant", opts)
}

// Get fetches one item by ID.
func (s *TenantsService) Get(ctx context.Context, id string, opts ...Option) (*TenantResponse, error) {
	return getJSON[TenantResponse](ctx, s.c, p("tenant", id), opts)
}

// Create creates an item.
func (s *TenantsService) Create(ctx context.Context, body TenantCreateRequest, opts ...Option) (*TenantResponse, error) {
	return sendJSON[TenantResponse](ctx, s.c, http.MethodPost, "tenant", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *TenantsService) Update(ctx context.Context, id string, body TenantUpdateRequest, opts ...Option) (*TenantResponse, error) {
	return sendJSON[TenantResponse](ctx, s.c, http.MethodPatch, p("tenant", id), body, opts)
}

// Delete deletes an item.
func (s *TenantsService) Delete(ctx context.Context, id string, opts ...Option) (*TenantResponse, error) {
	return sendJSON[TenantResponse](ctx, s.c, http.MethodDelete, p("tenant", id), nil, opts)
}

// SandboxIntegratedAccountsService manages sandbox integrated accounts.
type SandboxIntegratedAccountsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *SandboxIntegratedAccountsService) List(ctx context.Context, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, "sandbox-integrated-account", opts)
}

// Get fetches one item by ID.
func (s *SandboxIntegratedAccountsService) Get(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, p("sandbox-integrated-account", id), opts)
}

// Create creates an item.
func (s *SandboxIntegratedAccountsService) Create(ctx context.Context, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, "sandbox-integrated-account", body, opts)
}

// Delete deletes an item.
func (s *SandboxIntegratedAccountsService) Delete(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodDelete, p("sandbox-integrated-account", id), nil, opts)
}

// EnvironmentIntegrationsService manages integrations installed in the environment.
type EnvironmentIntegrationsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *EnvironmentIntegrationsService) List(ctx context.Context, opts ...Option) *Iterator[EnvironmentIntegrationResponse] {
	return newIterator[EnvironmentIntegrationResponse](ctx, s.c, "environment-integration", opts)
}

// Get fetches one item by ID.
func (s *EnvironmentIntegrationsService) Get(ctx context.Context, id string, opts ...Option) (*EnvironmentIntegrationResponse, error) {
	return getJSON[EnvironmentIntegrationResponse](ctx, s.c, p("environment-integration", id), opts)
}

// Create creates an item.
func (s *EnvironmentIntegrationsService) Create(ctx context.Context, body EnvironmentIntegrationCreateRequest, opts ...Option) (*EnvironmentIntegrationResponse, error) {
	return sendJSON[EnvironmentIntegrationResponse](ctx, s.c, http.MethodPost, "environment-integration", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *EnvironmentIntegrationsService) Update(ctx context.Context, id string, body EnvironmentIntegrationUpdateRequest, opts ...Option) (*EnvironmentIntegrationResponse, error) {
	return sendJSON[EnvironmentIntegrationResponse](ctx, s.c, http.MethodPatch, p("environment-integration", id), body, opts)
}

// Delete deletes an item.
func (s *EnvironmentIntegrationsService) Delete(ctx context.Context, id string, opts ...Option) (*EnvironmentIntegrationDeleteResponse, error) {
	return sendJSON[EnvironmentIntegrationDeleteResponse](ctx, s.c, http.MethodDelete, p("environment-integration", id), nil, opts)
}

// EnvironmentUnifiedModelsService manages unified models installed in the environment.
type EnvironmentUnifiedModelsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *EnvironmentUnifiedModelsService) List(ctx context.Context, opts ...Option) *Iterator[EnvironmentUnifiedModelResponse] {
	return newIterator[EnvironmentUnifiedModelResponse](ctx, s.c, "environment-unified-model", opts)
}

// Get fetches one item by ID.
func (s *EnvironmentUnifiedModelsService) Get(ctx context.Context, id string, opts ...Option) (*EnvironmentUnifiedModelResponse, error) {
	return getJSON[EnvironmentUnifiedModelResponse](ctx, s.c, p("environment-unified-model", id), opts)
}

// Create creates an item.
func (s *EnvironmentUnifiedModelsService) Create(ctx context.Context, body EnvironmentUnifiedModelCreateRequest, opts ...Option) (*EnvironmentUnifiedModelResponse, error) {
	return sendJSON[EnvironmentUnifiedModelResponse](ctx, s.c, http.MethodPost, "environment-unified-model", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *EnvironmentUnifiedModelsService) Update(ctx context.Context, id string, body EnvironmentUnifiedModelUpdateRequest, opts ...Option) (*EnvironmentUnifiedModelResponse, error) {
	return sendJSON[EnvironmentUnifiedModelResponse](ctx, s.c, http.MethodPatch, p("environment-unified-model", id), body, opts)
}

// Delete deletes an item.
func (s *EnvironmentUnifiedModelsService) Delete(ctx context.Context, id string, opts ...Option) (*EnvironmentUnifiedModelDeleteResponse, error) {
	return sendJSON[EnvironmentUnifiedModelDeleteResponse](ctx, s.c, http.MethodDelete, p("environment-unified-model", id), nil, opts)
}

// UnifiedModelResourceMethodsService manages unified model resource method mappings.
type UnifiedModelResourceMethodsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *UnifiedModelResourceMethodsService) List(ctx context.Context, opts ...Option) *Iterator[UnifiedModelResourceMethodRowResponse] {
	return newIterator[UnifiedModelResourceMethodRowResponse](ctx, s.c, "unified-model-resource-method", opts)
}

// Get fetches one item by ID.
func (s *UnifiedModelResourceMethodsService) Get(ctx context.Context, id string, opts ...Option) (*UnifiedModelResourceMethodRowResponse, error) {
	return getJSON[UnifiedModelResourceMethodRowResponse](ctx, s.c, p("unified-model-resource-method", id), opts)
}

// Create creates an item.
func (s *UnifiedModelResourceMethodsService) Create(ctx context.Context, body UnifiedModelResourceMethodCreateRequest, opts ...Option) (*UnifiedModelResourceMethodRowResponse, error) {
	return sendJSON[UnifiedModelResourceMethodRowResponse](ctx, s.c, http.MethodPost, "unified-model-resource-method", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *UnifiedModelResourceMethodsService) Update(ctx context.Context, id string, body UnifiedModelResourceMethodUpdateRequest, opts ...Option) (*UnifiedModelResourceMethodRowResponse, error) {
	return sendJSON[UnifiedModelResourceMethodRowResponse](ctx, s.c, http.MethodPatch, p("unified-model-resource-method", id), body, opts)
}

// Delete deletes an item.
func (s *UnifiedModelResourceMethodsService) Delete(ctx context.Context, id string, opts ...Option) (*UnifiedModelResourceMethodDeleteResponse, error) {
	return sendJSON[UnifiedModelResourceMethodDeleteResponse](ctx, s.c, http.MethodDelete, p("unified-model-resource-method", id), nil, opts)
}

// EnvironmentUnifiedModelResourceMethodsService manages environment-level overrides of unified model resource methods.
type EnvironmentUnifiedModelResourceMethodsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *EnvironmentUnifiedModelResourceMethodsService) List(ctx context.Context, opts ...Option) *Iterator[EnvironmentUnifiedModelResourceMethodRowResponse] {
	return newIterator[EnvironmentUnifiedModelResourceMethodRowResponse](ctx, s.c, "environment-unified-model-resource-method", opts)
}

// Get fetches one item by ID.
func (s *EnvironmentUnifiedModelResourceMethodsService) Get(ctx context.Context, id string, opts ...Option) (*EnvironmentUnifiedModelResourceMethodRowResponse, error) {
	return getJSON[EnvironmentUnifiedModelResourceMethodRowResponse](ctx, s.c, p("environment-unified-model-resource-method", id), opts)
}

// Create creates an item.
func (s *EnvironmentUnifiedModelResourceMethodsService) Create(ctx context.Context, body EnvironmentUnifiedModelResourceMethodCreateRequest, opts ...Option) (*EnvironmentUnifiedModelResourceMethodRowResponse, error) {
	return sendJSON[EnvironmentUnifiedModelResourceMethodRowResponse](ctx, s.c, http.MethodPost, "environment-unified-model-resource-method", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *EnvironmentUnifiedModelResourceMethodsService) Update(ctx context.Context, id string, body EnvironmentUnifiedModelResourceMethodUpdateRequest, opts ...Option) (*EnvironmentUnifiedModelResourceMethodRowResponse, error) {
	return sendJSON[EnvironmentUnifiedModelResourceMethodRowResponse](ctx, s.c, http.MethodPatch, p("environment-unified-model-resource-method", id), body, opts)
}

// Delete deletes an item.
func (s *EnvironmentUnifiedModelResourceMethodsService) Delete(ctx context.Context, id string, opts ...Option) (*EnvironmentUnifiedModelResourceMethodDeleteResponse, error) {
	return sendJSON[EnvironmentUnifiedModelResourceMethodDeleteResponse](ctx, s.c, http.MethodDelete, p("environment-unified-model-resource-method", id), nil, opts)
}

// SyncJobsService manages sync jobs.
type SyncJobsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *SyncJobsService) List(ctx context.Context, opts ...Option) *Iterator[SyncJobSchema] {
	return newIterator[SyncJobSchema](ctx, s.c, "sync-job", opts)
}

// Get fetches one item by ID.
func (s *SyncJobsService) Get(ctx context.Context, id string, opts ...Option) (*SyncJobSchema, error) {
	return getJSON[SyncJobSchema](ctx, s.c, p("sync-job", id), opts)
}

// Create creates an item.
func (s *SyncJobsService) Create(ctx context.Context, body SyncJobCreateRequest, opts ...Option) (*SyncJobSchema, error) {
	return sendJSON[SyncJobSchema](ctx, s.c, http.MethodPost, "sync-job", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *SyncJobsService) Update(ctx context.Context, id string, body SyncJobUpdateRequest, opts ...Option) (*SyncJobSchema, error) {
	return sendJSON[SyncJobSchema](ctx, s.c, http.MethodPatch, p("sync-job", id), body, opts)
}

// Delete deletes an item.
func (s *SyncJobsService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("sync-job", id), nil, opts)
}

// SyncJobRunsService manages sync job runs.
type SyncJobRunsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *SyncJobRunsService) List(ctx context.Context, opts ...Option) *Iterator[SyncJobRunResponse] {
	return newIterator[SyncJobRunResponse](ctx, s.c, "sync-job-run", opts)
}

// Get fetches one item by ID.
func (s *SyncJobRunsService) Get(ctx context.Context, id string, opts ...Option) (*SyncJobRunResponse, error) {
	return getJSON[SyncJobRunResponse](ctx, s.c, p("sync-job-run", id), opts)
}

// Create creates an item.
func (s *SyncJobRunsService) Create(ctx context.Context, body SyncJobRunCreateRequest, opts ...Option) (*SyncJobRunResponse, error) {
	return sendJSON[SyncJobRunResponse](ctx, s.c, http.MethodPost, "sync-job-run", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *SyncJobRunsService) Update(ctx context.Context, id string, body Record, opts ...Option) (*SyncJobRunResponse, error) {
	return sendJSON[SyncJobRunResponse](ctx, s.c, http.MethodPatch, p("sync-job-run", id), body, opts)
}

// Delete deletes an item.
func (s *SyncJobRunsService) Delete(ctx context.Context, id string, opts ...Option) (*SyncJobRunResponse, error) {
	return sendJSON[SyncJobRunResponse](ctx, s.c, http.MethodDelete, p("sync-job-run", id), nil, opts)
}

// SyncJobTemplatesService manages sync job templates.
type SyncJobTemplatesService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *SyncJobTemplatesService) List(ctx context.Context, opts ...Option) *Iterator[SyncJobTemplateSchema] {
	return newIterator[SyncJobTemplateSchema](ctx, s.c, "sync-job-template", opts)
}

// Get fetches one item by ID.
func (s *SyncJobTemplatesService) Get(ctx context.Context, id string, opts ...Option) (*SyncJobTemplateSchema, error) {
	return getJSON[SyncJobTemplateSchema](ctx, s.c, p("sync-job-template", id), opts)
}

// Create creates an item.
func (s *SyncJobTemplatesService) Create(ctx context.Context, body SyncJobTemplateCreate, opts ...Option) (*SyncJobTemplateSchema, error) {
	return sendJSON[SyncJobTemplateSchema](ctx, s.c, http.MethodPost, "sync-job-template", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *SyncJobTemplatesService) Update(ctx context.Context, id string, body SyncJobTemplatePatch, opts ...Option) (*SyncJobTemplateSchema, error) {
	return sendJSON[SyncJobTemplateSchema](ctx, s.c, http.MethodPatch, p("sync-job-template", id), body, opts)
}

// Delete deletes an item.
func (s *SyncJobTemplatesService) Delete(ctx context.Context, id string, opts ...Option) (*SyncJobTemplateDeleteResponse, error) {
	return sendJSON[SyncJobTemplateDeleteResponse](ctx, s.c, http.MethodDelete, p("sync-job-template", id), nil, opts)
}

// SyncJobRunStatesService manages state stored by sync job runs.
type SyncJobRunStatesService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *SyncJobRunStatesService) List(ctx context.Context, opts ...Option) *Iterator[SyncJobRunState] {
	return newIterator[SyncJobRunState](ctx, s.c, "sync-job-run-state", opts)
}

// Get fetches one item by ID.
func (s *SyncJobRunStatesService) Get(ctx context.Context, id string, opts ...Option) (*SyncJobRunState, error) {
	return getJSON[SyncJobRunState](ctx, s.c, p("sync-job-run-state", id), opts)
}

// Create creates an item.
func (s *SyncJobRunStatesService) Create(ctx context.Context, body SyncJobRunStateCreateRequest, opts ...Option) (*SyncJobRunState, error) {
	return sendJSON[SyncJobRunState](ctx, s.c, http.MethodPost, "sync-job-run-state", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *SyncJobRunStatesService) Update(ctx context.Context, id string, body SyncJobRunStateUpdateRequest, opts ...Option) (*SyncJobRunState, error) {
	return sendJSON[SyncJobRunState](ctx, s.c, http.MethodPatch, p("sync-job-run-state", id), body, opts)
}

// Delete deletes an item.
func (s *SyncJobRunStatesService) Delete(ctx context.Context, id string, opts ...Option) (*SyncJobRunState, error) {
	return sendJSON[SyncJobRunState](ctx, s.c, http.MethodDelete, p("sync-job-run-state", id), nil, opts)
}

// WorkflowsService manages workflows.
type WorkflowsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *WorkflowsService) List(ctx context.Context, opts ...Option) *Iterator[WorkflowSchema] {
	return newIterator[WorkflowSchema](ctx, s.c, "workflow", opts)
}

// Get fetches one item by ID.
func (s *WorkflowsService) Get(ctx context.Context, id string, opts ...Option) (*WorkflowSchema, error) {
	return getJSON[WorkflowSchema](ctx, s.c, p("workflow", id), opts)
}

// Create creates an item.
func (s *WorkflowsService) Create(ctx context.Context, body WorkflowCreateSchema, opts ...Option) (*WorkflowSchema, error) {
	return sendJSON[WorkflowSchema](ctx, s.c, http.MethodPost, "workflow", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *WorkflowsService) Update(ctx context.Context, id string, body Record, opts ...Option) (*WorkflowSchema, error) {
	return sendJSON[WorkflowSchema](ctx, s.c, http.MethodPatch, p("workflow", id), body, opts)
}

// Delete deletes an item.
func (s *WorkflowsService) Delete(ctx context.Context, id string, opts ...Option) (*WorkflowDeleteResponse, error) {
	return sendJSON[WorkflowDeleteResponse](ctx, s.c, http.MethodDelete, p("workflow", id), nil, opts)
}

// WorkflowRunsService manages workflow runs.
type WorkflowRunsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *WorkflowRunsService) List(ctx context.Context, opts ...Option) *Iterator[WorkflowRunResponse] {
	return newIterator[WorkflowRunResponse](ctx, s.c, "workflow-run", opts)
}

// Get fetches one item by ID.
func (s *WorkflowRunsService) Get(ctx context.Context, id string, opts ...Option) (*WorkflowRunResponse, error) {
	return getJSON[WorkflowRunResponse](ctx, s.c, p("workflow-run", id), opts)
}

// Create creates an item.
func (s *WorkflowRunsService) Create(ctx context.Context, body Record, opts ...Option) (*WorkflowRunResponse, error) {
	return sendJSON[WorkflowRunResponse](ctx, s.c, http.MethodPost, "workflow-run", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *WorkflowRunsService) Update(ctx context.Context, id string, body Record, opts ...Option) (*WorkflowRunResponse, error) {
	return sendJSON[WorkflowRunResponse](ctx, s.c, http.MethodPatch, p("workflow-run", id), body, opts)
}

// Delete deletes an item.
func (s *WorkflowRunsService) Delete(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodDelete, p("workflow-run", id), nil, opts)
}

// DocumentationService manages documentation records.
type DocumentationService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *DocumentationService) List(ctx context.Context, opts ...Option) *Iterator[Documentation] {
	return newIterator[Documentation](ctx, s.c, "documentation", opts)
}

// Get fetches one item by ID.
func (s *DocumentationService) Get(ctx context.Context, id string, opts ...Option) (*Documentation, error) {
	return getJSON[Documentation](ctx, s.c, p("documentation", id), opts)
}

// Create creates an item.
func (s *DocumentationService) Create(ctx context.Context, body DocumentationCreateRequest, opts ...Option) (*Documentation, error) {
	return sendJSON[Documentation](ctx, s.c, http.MethodPost, "documentation", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *DocumentationService) Update(ctx context.Context, id string, body DocumentationPatchRequest, opts ...Option) (*Documentation, error) {
	return sendJSON[Documentation](ctx, s.c, http.MethodPatch, p("documentation", id), body, opts)
}

// Delete deletes an item.
func (s *DocumentationService) Delete(ctx context.Context, id string, opts ...Option) error {
	_, err := s.c.do(ctx, request{method: http.MethodDelete, path: p("documentation", id)}, opts)
	return err
}

// LogsService manages logs. The API requires the log_type query parameter (unified_proxy_api, rapid_bridge, webhook, mcp, sync_job_cron_trigger), for example WithQuery(truto.Params{"log_type": "unified_proxy_api"}).
type LogsService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *LogsService) List(ctx context.Context, opts ...Option) *Iterator[Log] {
	return newIterator[Log](ctx, s.c, "log", opts)
}

// StaticGatesService manages static gates.
type StaticGatesService struct{ c *Client }

// List returns an iterator over the items. Filter with WithQuery.
func (s *StaticGatesService) List(ctx context.Context, opts ...Option) *Iterator[StaticGateResponse] {
	return newIterator[StaticGateResponse](ctx, s.c, "static-gate", opts)
}

// Get fetches one item by ID.
func (s *StaticGatesService) Get(ctx context.Context, id string, opts ...Option) (*StaticGateResponse, error) {
	return getJSON[StaticGateResponse](ctx, s.c, p("static-gate", id), opts)
}

// Create creates an item.
func (s *StaticGatesService) Create(ctx context.Context, body StaticGateCreateRequest, opts ...Option) (*StaticGateCreateResponse, error) {
	return sendJSON[StaticGateCreateResponse](ctx, s.c, http.MethodPost, "static-gate", body, opts)
}

// Update changes an item. Only the fields that are set are sent.
func (s *StaticGatesService) Update(ctx context.Context, id string, body StaticGateUpdateRequest, opts ...Option) (*StaticGateResponse, error) {
	return sendJSON[StaticGateResponse](ctx, s.c, http.MethodPatch, p("static-gate", id), body, opts)
}

// Delete deletes an item.
func (s *StaticGatesService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("static-gate", id), nil, opts)
}
