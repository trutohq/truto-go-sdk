# Changelog

## 0.1.0

First release.

- A client for the Unified, Proxy and Custom APIs and MCP, and for every admin resource in the Python and TypeScript SDKs, plus session tokens.
- Typed models generated from Truto's OpenAPI spec, with open enums and union types.
- Retries with backoff and `Retry-After` support, a per-attempt timeout, and retry and rate-limit hooks.
- Lazy cursor pagination with `Next`, `NextPage`, `All`, `First` and `Seq` (range-over-func).
- `truto.Error` with status, parsed body and headers, `ErrTimeout`, and `IsNotFound`, `IsUnauthorized` and `IsRateLimited`.
