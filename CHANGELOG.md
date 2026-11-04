# Changelog

## Unreleased

- Fix: listing the full integration catalog no longer fails to decode. A resource method's `body` is now any JSON value (`any`), because some integrations (greythr) send a JSON array there.
- Header maps in integration configs are now `map[string]*string`, so a `null` header (which removes an inherited header) is kept instead of becoming `""`.
- Fixes to the published spec now live in `spec/overlay.json` (OpenAPI Overlay), applied by `go generate`.
- `examples/proxy` only reads now. It used to update a record.

## 0.1.0

First release.

- A client for the Unified, Proxy and Custom APIs and MCP, and for every admin resource in the Python and TypeScript SDKs, plus session tokens.
- Typed models generated from Truto's OpenAPI spec, with open enums and union types.
- Retries with backoff and `Retry-After` support, a per-attempt timeout, and retry and rate-limit hooks.
- Lazy cursor pagination with `Next`, `NextPage`, `All`, `First` and `Seq` (range-over-func).
- `truto.Error` with status, parsed body and headers, `ErrTimeout`, and `IsNotFound`, `IsUnauthorized` and `IsRateLimited`.
