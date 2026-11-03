package truto

import (
	"context"
	"net/http"
)

// IntegrationsService manages integration definitions.
type IntegrationsService struct{ c *Client }

// List returns an iterator over integrations.
func (s *IntegrationsService) List(ctx context.Context, opts ...Option) *Iterator[IntegrationResponse] {
	return newIterator[IntegrationResponse](ctx, s.c, "integration", opts)
}

// Get fetches an integration.
func (s *IntegrationsService) Get(ctx context.Context, id string, opts ...Option) (*IntegrationResponse, error) {
	return getJSON[IntegrationResponse](ctx, s.c, p("integration", id), opts)
}

// Create creates a private integration.
func (s *IntegrationsService) Create(ctx context.Context, body IntegrationCreateRequest, opts ...Option) (*IntegrationResponse, error) {
	return sendJSON[IntegrationResponse](ctx, s.c, http.MethodPost, "integration", body, opts)
}

// Update changes an integration.
func (s *IntegrationsService) Update(ctx context.Context, id string, body IntegrationUpdateRequest, opts ...Option) (*IntegrationResponse, error) {
	return sendJSON[IntegrationResponse](ctx, s.c, http.MethodPatch, p("integration", id), body, opts)
}

// Delete deletes an integration.
func (s *IntegrationsService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("integration", id), nil, opts)
}

// Install installs the integration in the environment the token belongs to.
func (s *IntegrationsService) Install(ctx context.Context, integrationID string, isEnabled bool, opts ...Option) (*EnvironmentIntegrationResponse, error) {
	body := Record{"integration_id": integrationID, "is_enabled": isEnabled}
	return sendJSON[EnvironmentIntegrationResponse](ctx, s.c, http.MethodPost, "environment-integration", body, opts)
}

// Uninstall removes the integration from the environment. It looks up the
// environment integrations for integrationID and deletes each one.
func (s *IntegrationsService) Uninstall(ctx context.Context, integrationID string, opts ...Option) ([]EnvironmentIntegrationDeleteResponse, error) {
	listOpts := append(append([]Option(nil), opts...), WithQuery(Params{"integration_id": integrationID}))
	installed, err := newIterator[EnvironmentIntegrationResponse](ctx, s.c, "environment-integration", listOpts).All()
	if err != nil {
		return nil, err
	}
	out := make([]EnvironmentIntegrationDeleteResponse, 0, len(installed))
	for _, ei := range installed {
		if ei.ID == nil {
			continue
		}
		res, err := sendJSON[EnvironmentIntegrationDeleteResponse](ctx, s.c, http.MethodDelete, p("environment-integration", *ei.ID), nil, opts)
		if err != nil {
			return out, err
		}
		out = append(out, *res)
	}
	return out, nil
}

// Tools lists the tools the integration exposes.
func (s *IntegrationsService) Tools(ctx context.Context, id string, opts ...Option) (*IntegrationToolsResponse, error) {
	return getJSON[IntegrationToolsResponse](ctx, s.c, p("integration", id, "tools"), opts)
}

// UnifiedAPIs lists the unified models and resources the integration supports.
func (s *IntegrationsService) UnifiedAPIs(ctx context.Context, id string, opts ...Option) (*IntegrationUnifiedApisResponse, error) {
	return getJSON[IntegrationUnifiedApisResponse](ctx, s.c, p("integration", id, "unified-apis"), opts)
}

// Capabilities describes the integration's resources, methods and auth.
func (s *IntegrationsService) Capabilities(ctx context.Context, id string, opts ...Option) (*CapabilitiesResponse, error) {
	return getJSON[CapabilitiesResponse](ctx, s.c, p("integration", id, "capabilities"), opts)
}

// UnifiedModelsService manages unified model definitions.
type UnifiedModelsService struct{ c *Client }

// List returns an iterator over unified models.
func (s *UnifiedModelsService) List(ctx context.Context, opts ...Option) *Iterator[UnifiedModelResponse] {
	return newIterator[UnifiedModelResponse](ctx, s.c, "unified-model", opts)
}

// Get fetches a unified model.
func (s *UnifiedModelsService) Get(ctx context.Context, id string, opts ...Option) (*UnifiedModelResponse, error) {
	return getJSON[UnifiedModelResponse](ctx, s.c, p("unified-model", id), opts)
}

// Create creates a unified model.
func (s *UnifiedModelsService) Create(ctx context.Context, body UnifiedModelCreateRequest, opts ...Option) (*UnifiedModelResponse, error) {
	return sendJSON[UnifiedModelResponse](ctx, s.c, http.MethodPost, "unified-model", body, opts)
}

// Update changes a unified model.
func (s *UnifiedModelsService) Update(ctx context.Context, id string, body UnifiedModelUpdateRequest, opts ...Option) (*UnifiedModelResponse, error) {
	return sendJSON[UnifiedModelResponse](ctx, s.c, http.MethodPatch, p("unified-model", id), body, opts)
}

// Delete deletes a unified model.
func (s *UnifiedModelsService) Delete(ctx context.Context, id string, opts ...Option) (*UnifiedModelDeleteResponse, error) {
	return sendJSON[UnifiedModelDeleteResponse](ctx, s.c, http.MethodDelete, p("unified-model", id), nil, opts)
}

// Install installs the unified model in the environment the token belongs to.
func (s *UnifiedModelsService) Install(ctx context.Context, unifiedModelID string, opts ...Option) (*EnvironmentUnifiedModelResponse, error) {
	body := Record{"unified_model_id": unifiedModelID}
	return sendJSON[EnvironmentUnifiedModelResponse](ctx, s.c, http.MethodPost, "environment-unified-model", body, opts)
}

// Uninstall removes the unified model from the environment. It looks up the
// environment unified models for unifiedModelID and deletes each one.
func (s *UnifiedModelsService) Uninstall(ctx context.Context, unifiedModelID string, opts ...Option) ([]EnvironmentUnifiedModelDeleteResponse, error) {
	listOpts := append(append([]Option(nil), opts...), WithQuery(Params{"unified_model_id": unifiedModelID}))
	installed, err := newIterator[EnvironmentUnifiedModelResponse](ctx, s.c, "environment-unified-model", listOpts).All()
	if err != nil {
		return nil, err
	}
	out := make([]EnvironmentUnifiedModelDeleteResponse, 0, len(installed))
	for _, eum := range installed {
		if eum.ID == nil {
			continue
		}
		res, err := sendJSON[EnvironmentUnifiedModelDeleteResponse](ctx, s.c, http.MethodDelete, p("environment-unified-model", *eum.ID), nil, opts)
		if err != nil {
			return out, err
		}
		out = append(out, *res)
	}
	return out, nil
}

// SyncJobCronTriggersService manages cron schedules for sync jobs.
type SyncJobCronTriggersService struct{ c *Client }

// List returns an iterator over sync job cron triggers.
func (s *SyncJobCronTriggersService) List(ctx context.Context, opts ...Option) *Iterator[SyncJobCronTriggerSchema] {
	return newIterator[SyncJobCronTriggerSchema](ctx, s.c, "sync-job-cron-trigger", opts)
}

// Get fetches a sync job cron trigger.
func (s *SyncJobCronTriggersService) Get(ctx context.Context, id string, opts ...Option) (*SyncJobCronTriggerSchema, error) {
	return getJSON[SyncJobCronTriggerSchema](ctx, s.c, p("sync-job-cron-trigger", id), opts)
}

// Create creates a sync job cron trigger.
func (s *SyncJobCronTriggersService) Create(ctx context.Context, body SyncJobCronTriggerCreateRequest, opts ...Option) (*SyncJobCronTriggerSchema, error) {
	return sendJSON[SyncJobCronTriggerSchema](ctx, s.c, http.MethodPost, "sync-job-cron-trigger", body, opts)
}

// Update changes a sync job cron trigger.
func (s *SyncJobCronTriggersService) Update(ctx context.Context, id string, body SyncJobCronTriggerUpdateRequest, opts ...Option) (*SyncJobCronTriggerSchema, error) {
	return sendJSON[SyncJobCronTriggerSchema](ctx, s.c, http.MethodPatch, p("sync-job-cron-trigger", id), body, opts)
}

// Delete deletes a sync job cron trigger.
func (s *SyncJobCronTriggersService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("sync-job-cron-trigger", id), nil, opts)
}

// Schedule (re)schedules the trigger's next run.
func (s *SyncJobCronTriggersService) Schedule(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, p("sync-job-cron-trigger", id, "schedule"), nil, opts)
}

// RecordService is a CRUD service for resources the API spec does not
// describe yet; items are Records. Use Decode for your own types.
type RecordService struct {
	c    *Client
	path string
}

// List returns an iterator over the items.
func (s *RecordService) List(ctx context.Context, opts ...Option) *Iterator[Record] {
	return newIterator[Record](ctx, s.c, s.path, opts)
}

// Get fetches one item by ID.
func (s *RecordService) Get(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return getJSON[Record](ctx, s.c, p(s.path, id), opts)
}

// Create creates an item.
func (s *RecordService) Create(ctx context.Context, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, s.path, body, opts)
}

// Update changes an item.
func (s *RecordService) Update(ctx context.Context, id string, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPatch, p(s.path, id), body, opts)
}

// Delete deletes an item.
func (s *RecordService) Delete(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodDelete, p(s.path, id), nil, opts)
}

// DaemonJobCronTriggersService manages cron schedules for daemon jobs.
type DaemonJobCronTriggersService struct{ RecordService }

// Schedule (re)schedules the trigger's next run.
func (s *DaemonJobCronTriggersService) Schedule(ctx context.Context, id string, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, p(s.path, id, "schedule"), nil, opts)
}

// WebhooksService manages outbound webhooks Truto sends to you.
type WebhooksService struct{ c *Client }

// List returns an iterator over webhooks.
func (s *WebhooksService) List(ctx context.Context, opts ...Option) *Iterator[WebhookSchema] {
	return newIterator[WebhookSchema](ctx, s.c, "webhook", opts)
}

// Get fetches a webhook.
func (s *WebhooksService) Get(ctx context.Context, id string, opts ...Option) (*WebhookSchema, error) {
	return getJSON[WebhookSchema](ctx, s.c, p("webhook", id), opts)
}

// Create creates a webhook. The response includes the signing secret, which is
// not returned again.
func (s *WebhooksService) Create(ctx context.Context, body WebhookCreateRequest, opts ...Option) (*WebhookCreateResponse, error) {
	return sendJSON[WebhookCreateResponse](ctx, s.c, http.MethodPost, "webhook", body, opts)
}

// Update changes a webhook.
func (s *WebhooksService) Update(ctx context.Context, id string, body WebhookUpdateRequest, opts ...Option) (*WebhookSchema, error) {
	return sendJSON[WebhookSchema](ctx, s.c, http.MethodPatch, p("webhook", id), body, opts)
}

// Delete deletes a webhook.
func (s *WebhooksService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("webhook", id), nil, opts)
}

// Test sends a test event to a webhook.
func (s *WebhooksService) Test(ctx context.Context, body WebhookTestRequest, opts ...Option) (*WebhookTestResponse, error) {
	return sendJSON[WebhookTestResponse](ctx, s.c, http.MethodPost, "webhook/test", body, opts)
}

// DatastoresService manages datastores that sync jobs can write to.
type DatastoresService struct{ c *Client }

// List returns an iterator over datastores.
func (s *DatastoresService) List(ctx context.Context, opts ...Option) *Iterator[DatastoreSchema] {
	return newIterator[DatastoreSchema](ctx, s.c, "datastore", opts)
}

// Get fetches a datastore.
func (s *DatastoresService) Get(ctx context.Context, id string, opts ...Option) (*DatastoreSchema, error) {
	return getJSON[DatastoreSchema](ctx, s.c, p("datastore", id), opts)
}

// Create creates a datastore.
func (s *DatastoresService) Create(ctx context.Context, body DatastoreCreateRequest, opts ...Option) (*DatastoreSchema, error) {
	return sendJSON[DatastoreSchema](ctx, s.c, http.MethodPost, "datastore", body, opts)
}

// Update changes a datastore.
func (s *DatastoresService) Update(ctx context.Context, id string, body DatastoreUpdateRequest, opts ...Option) (*DatastoreSchema, error) {
	return sendJSON[DatastoreSchema](ctx, s.c, http.MethodPatch, p("datastore", id), body, opts)
}

// Delete deletes a datastore.
func (s *DatastoresService) Delete(ctx context.Context, id string, opts ...Option) (*DatastoreDeleteResponse, error) {
	return sendJSON[DatastoreDeleteResponse](ctx, s.c, http.MethodDelete, p("datastore", id), nil, opts)
}

// Test runs one of the datastore's validation methods.
func (s *DatastoresService) Test(ctx context.Context, id, method string, body Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, p("datastore", id, "test", method), body, opts)
}

// ObservabilityDestinationsService manages where Truto ships logs.
type ObservabilityDestinationsService struct{ c *Client }

// List returns an iterator over observability destinations.
func (s *ObservabilityDestinationsService) List(ctx context.Context, opts ...Option) *Iterator[ObservabilityDestinationResponse] {
	return newIterator[ObservabilityDestinationResponse](ctx, s.c, "observability-destination", opts)
}

// Get fetches an observability destination.
func (s *ObservabilityDestinationsService) Get(ctx context.Context, id string, opts ...Option) (*ObservabilityDestinationResponse, error) {
	return getJSON[ObservabilityDestinationResponse](ctx, s.c, p("observability-destination", id), opts)
}

// Create creates an observability destination.
func (s *ObservabilityDestinationsService) Create(ctx context.Context, body ObservabilityDestinationCreateRequest, opts ...Option) (*ObservabilityDestinationResponse, error) {
	return sendJSON[ObservabilityDestinationResponse](ctx, s.c, http.MethodPost, "observability-destination", body, opts)
}

// Update changes an observability destination.
func (s *ObservabilityDestinationsService) Update(ctx context.Context, id string, body ObservabilityDestinationUpdateRequest, opts ...Option) (*ObservabilityDestinationResponse, error) {
	return sendJSON[ObservabilityDestinationResponse](ctx, s.c, http.MethodPatch, p("observability-destination", id), body, opts)
}

// Delete deletes an observability destination.
func (s *ObservabilityDestinationsService) Delete(ctx context.Context, id string, opts ...Option) (*DeleteResponse, error) {
	return sendJSON[DeleteResponse](ctx, s.c, http.MethodDelete, p("observability-destination", id), nil, opts)
}

// Test sends a test event to the destination, optionally with new config and
// secrets that are not saved.
func (s *ObservabilityDestinationsService) Test(ctx context.Context, id string, body ObservabilityDestinationTestRequest, opts ...Option) (*ObservabilityDestinationTestResponse, error) {
	return sendJSON[ObservabilityDestinationTestResponse](ctx, s.c, http.MethodPost, p("observability-destination", id, "test"), body, opts)
}

// AnalyticsService reads usage analytics. Pass filters such as environment_id,
// metric or range with WithQuery.
type AnalyticsService struct{ c *Client }

// Status reports whether analytics are available for the environment.
func (s *AnalyticsService) Status(ctx context.Context, opts ...Option) (*AnalyticsStatusResponse, error) {
	return getJSON[AnalyticsStatusResponse](ctx, s.c, "analytics/status", opts)
}

// Metrics lists the metrics and dimensions that can be queried.
func (s *AnalyticsService) Metrics(ctx context.Context, opts ...Option) (*AnalyticsMetricsResponse, error) {
	return getJSON[AnalyticsMetricsResponse](ctx, s.c, "analytics/metrics", opts)
}

// Dashboard returns the dashboard summary.
func (s *AnalyticsService) Dashboard(ctx context.Context, opts ...Option) (*AnalyticsDashboardResponse, error) {
	return getJSON[AnalyticsDashboardResponse](ctx, s.c, "analytics/dashboard", opts)
}

// Timeseries returns a metric over time.
func (s *AnalyticsService) Timeseries(ctx context.Context, opts ...Option) (*AnalyticsTimeseriesResponse, error) {
	return getJSON[AnalyticsTimeseriesResponse](ctx, s.c, "analytics/timeseries", opts)
}

// Breakdown returns a metric split by a dimension.
func (s *AnalyticsService) Breakdown(ctx context.Context, opts ...Option) (*AnalyticsBreakdownResponse, error) {
	return getJSON[AnalyticsBreakdownResponse](ctx, s.c, "analytics/breakdown", opts)
}

// BatchRequestsService runs several API requests in one call.
type BatchRequestsService struct{ c *Client }

// Create runs a batch request.
func (s *BatchRequestsService) Create(ctx context.Context, body BatchRequestSchema, opts ...Option) (*BatchRequestResponse, error) {
	return sendJSON[BatchRequestResponse](ctx, s.c, http.MethodPost, "batch-request", body, opts)
}

// WebhookProcessorService forwards an inbound third-party webhook payload to
// Truto for processing.
type WebhookProcessorService struct {
	c    *Client
	path string
}

// Process sends payload for the environment integration or integrated account id.
func (s *WebhookProcessorService) Process(ctx context.Context, id string, payload Record, opts ...Option) (*Record, error) {
	return sendJSON[Record](ctx, s.c, http.MethodPost, p(s.path, id), payload, opts)
}
