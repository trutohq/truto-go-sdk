# truto-go-sdk

Go client for the [Truto](https://truto.one) API. It covers the Unified, Proxy and Custom APIs, MCP, and the admin API (integrated accounts, integrations, sync jobs, webhooks, and the rest). Admin resources use typed models generated from Truto's OpenAPI spec.

- Go 1.23 or later, standard library only
- Retries with backoff, `Retry-After` support and per-request timeouts
- Cursor pagination as an iterator, including `range`-over-func
- Typed errors with the API's status and body

## Install

```sh
go get github.com/trutohq/truto-go-sdk
```

```go
import truto "github.com/trutohq/truto-go-sdk"
```

## Quick start

```go
ctx := context.Background()
client := truto.NewClient(os.Getenv("TRUTO_API_TOKEN"))

// Admin API: typed models
accounts, err := client.IntegratedAccounts.List(ctx,
	truto.WithQuery(truto.Params{"tenant_id": "acme"})).All()
if err != nil {
	log.Fatal(err)
}
for _, a := range accounts {
	fmt.Println(*a.ID, *a.Status)
}
```

Every method takes a `context.Context` first and accepts `Option`s last. Options passed to `NewClient` apply to every call. Options passed to a method apply to that call only.

| Option | What it does |
|---|---|
| `WithQuery(truto.Params{...})` | Query parameters. Nested maps and slices are encoded the way the API expects (`filter[status]=open`, `ids[0]=a`). |
| `WithTimeout(d)` | Time limit per attempt (default 30s). |
| `WithMaxRetries(n)` / `WithNoRetry()` | Number of retries (default 3), or none. |
| `WithRetryStatusCodes(codes...)` | Extra status codes to retry. |
| `WithOnRetry(fn)` / `WithOnRateLimited(fn)` | Hooks called before a retry or a rate-limit wait. |
| `WithHeader(k, v)` | Extra request header. |
| `WithBaseURL(url)`, `WithHTTPClient(c)`, `WithDefaultRetryAfter(d)` | Client setup. |

## Unified API

```go
contacts := client.Unified.List(ctx, truto.UnifiedResource{
	Model:               "crm",
	Resource:            "contacts",
	IntegratedAccountID: "ia_123",
}, truto.WithQuery(truto.Params{"limit": 50}))

for record, err := range contacts.Seq() {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(record["id"], record["first_name"])
}
```

Unified, Proxy and Custom API results depend on the integration, so they come back as `truto.Record` (`map[string]any`). Use `truto.Decode` to convert a record into your own type:

```go
type Contact struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
}

c, err := truto.Decode[Contact](record)
all, err := truto.DecodeAll[Contact](records)
```

Other unified calls: `Get`, `Create`, `Update`, `Delete`, `CustomMethod`, `Meta` and `ModelMeta`.

## Proxy API

```go
deals := truto.ProxyResource{Resource: "deals", IntegratedAccountID: "ia_123"}

deal, err := client.Proxy.Get(ctx, deals, "42")
_, err = client.Proxy.Update(ctx, deals, "42", map[string]any{"stage": "won"})
_, err = client.Proxy.CustomMethod(ctx, deals, "merge", map[string]any{"ids": []string{"1", "2"}})
```

## Custom API

A custom endpoint can return anything, so these calls return the raw response:

```go
res, err := client.Custom.Get(ctx, "ia_123", "/reports/monthly")
fmt.Println(res.ContentType, string(res.Body))
```

## Pagination

List methods return a lazy `*truto.Iterator[T]`. No request is sent until you read from it.

```go
it := client.SyncJobRuns.List(ctx, truto.WithQuery(truto.Params{"sync_job_id": "sj_1"}))

all, err := it.All()                // every item from every page
first, err := it.First()            // first item, or truto.Done if there are none
page, err := it.NextPage()          // one page: page.Items, page.NextCursor
item, err := it.Next()              // one item at a time; truto.Done at the end
for run, err := range it.Seq() {}   // range-over-func
```

## Errors

```go
_, err := client.IntegratedAccounts.Get(ctx, "missing")

var apiErr *truto.Error
switch {
case truto.IsNotFound(err):          // also IsUnauthorized, IsRateLimited
case errors.Is(err, truto.ErrTimeout):
case errors.As(err, &apiErr):
	fmt.Println(apiErr.StatusCode, apiErr.Message(), apiErr.Body)
}
```

`*truto.Error` carries the status code, the parsed body (`Body`), the raw body (`RawBody`) and the response headers.

## Retries

These statuses are retried: 408, 409, 425, 429, 502, 503 and 504, plus any added with `WithRetryStatusCodes`. Network errors are also retried. Timeouts are not.

- A 429 waits for the `Retry-After` header (10s by default when the header is missing).
- Other retryable statuses wait with exponential backoff and jitter, starting at 500ms and capped at 10s.

The request body is sent again on every attempt.

## Typed models

Admin API types such as `truto.IntegratedAccountResponse`, `truto.SyncJobCreateRequest` and `truto.WebhookSchema` are generated from `spec/openapi.json` into `models_gen.go`. Their behaviour:

- **Optional fields are pointers.** Use `truto.Ptr` to set them: `truto.TenantUpdateRequest{Name: truto.Ptr("Acme")}`.
- **Enums are open string types** with constants, such as `truto.IntegratedAccountStatus` and `truto.IntegratedAccountAuthenticationMethodOAuth2`. A value the SDK doesn't know yet still decodes unchanged.
- **Union types** (for example `truto.DatastoreConfig`) hold raw JSON. Decode a specific shape with the `As...` methods, such as `cfg.AsDatastoreS3Config()`, and build one with `truto.DatastoreConfigFrom...`.
- **Unknown JSON fields are ignored**, so new API fields don't break decoding.

Resources that the spec doesn't describe yet use `truto.Record`. These are teams, users, sandbox integrated accounts, daemons, daemon jobs and runs, and some request bodies.

Where the published spec doesn't match what the API returns, the fix lives in `spec/overlay.json`, an [OpenAPI Overlay](https://spec.openapis.org/overlay/v1.0.0) applied before generation. For example, resource-method `body` is typed as any JSON value, because some integrations send an array there. Header maps use `map[string]*string`, because a `null` value means "remove the inherited header".

To regenerate the models after updating `spec/openapi.json` or `spec/overlay.json` (needs Python 3):

```sh
go generate ./...
```

## Calling any endpoint

For an endpoint the SDK doesn't wrap yet:

```go
var out map[string]any
err := client.Do(ctx, "GET", "/some/path", nil, &out, truto.WithQuery(truto.Params{"a": 1}))
```

## Examples

See [`examples/`](examples): Unified API listing with typed decoding, the Proxy API, link tokens, and error and retry handling.

## Development

```sh
go test ./...
go vet ./...
```

The query-string encoder is tested against `testdata/qs_cases.json`. Those cases were generated with the `qs` library that the Truto API uses to parse query strings.

## License

MIT
