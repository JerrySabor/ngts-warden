# NGTS Warden — Public Go CLI and MCP Server

## Summary

Create a new public repository from the empty current directory:

- Repository/module: `github.com/JerrySabor/ngts-warden`
- Binary: `ngts-warden`
- Language: Go 1.26
- License: MIT, copyright © 2026 Jared Sabin
- Distribution: single static binary; no container in v1
- Publication: create the public GitHub repository and push verified `main`; prepare release automation but do not tag `v0.1.0`

Use Palo Alto Networks’ official OpenAPI specification pinned at commit `ea93c24d4ef2ff55e4efd195e6a2739cf7fb2ad0` and SHA-256 `b929b69e2d75a971a28ef30c92eb9c69b35c8a35717e45465ddc2b6e998e85ea`. It describes 95 paths and 147 operations: 55 GET, 52 POST, 18 DELETE, 12 PUT, and 10 PATCH operations. [Official pinned OpenAPI specification](https://raw.githubusercontent.com/PaloAltoNetworks/pan.dev/ea93c24d4ef2ff55e4efd195e6a2739cf7fb2ad0/openapi-specs/scm/config/ngts/tlsprotect-cloud.json)

Implement OAuth 2.0 client credentials with Basic authentication, `scope=tsg_id:<id>`, 15-minute bearer tokens, and API base URL `https://api.strata.paloaltonetworks.com/ngts`. [NGTS API documentation](https://pan.dev/scm/api/config/ngts/ngts-api/), [access-token documentation](https://pan.dev/scm/docs/access-tokens/)

## Public Interfaces

### CLI surface

Provide:

```text
ngts-warden version
ngts-warden doctor [--online] [--strict]
ngts-warden config init|list|show|use|delete
ngts-warden auth status|clear-cache
ngts-warden operations list|describe
ngts-warden <resource-group> <generated-operation> [flags]
ngts-warden request get PATH
ngts-warden request send METHOD PATH [flags]
ngts-warden mcp [--allow-write]
```

Generate a dedicated Cobra leaf command for every documented operation under these resource groups:

| CLI group | Ops | CLI group | Ops |
|---|---:|---|---:|
| `certificates` | 8 | `certificate-import` | 1 |
| `tls-server-endpoints` | 4 | `certificate-discovery` | 5 |
| `private-key-import` | 2 | `certificate-requests` | 6 |
| `certificate-policies` | 6 | `credentials` | 14 |
| `machine-installations` | 7 | `machine-types` | 1 |
| `machines` | 10 | `event-logs` | 3 |
| `vsatellite` | 14 | `inventory-monitoring` | 5 |
| `renewal-monitoring` | 4 | `certificate-tags` | 12 |
| `issuer-configurations` | 5 | `sub-ca-providers` | 5 |
| `workload-policies` | 5 | `issuer-certificates` | 1 |
| `certificate-approvals` | 8 | `revocation-approvals` | 5 |
| `plugins` | 8 | `built-in-accounts` | 8 |

Normalize operation IDs to kebab-case, strip each group’s declared operation prefix, and fail generation on collisions, missing operation IDs, or coverage drift. Generate flags for OpenAPI path, query, and header parameters. JSON request bodies accept `--body`, `--body-file FILE`, or stdin with `--body-file -`.

GET operations execute immediately. POST, PUT, PATCH, and DELETE commands preview the fully redacted request by default and require `--execute` to send it. This conservative rule also applies to POST-based search/export endpoints.

CSV, text, and binary responses support `--accept` and `--out`. Binary data must not be written to an interactive terminal. JSON file-result output reports path, byte count, content type, operation ID, and HTTP status.

### Authentication and configuration

Use this precedence:

1. Safe one-off flags such as `--profile`, `--access-token`, `--base-url`, and `--auth-url`
2. `NGTS_WARDEN_*` environment variables
3. Selected named profile
4. Official endpoint defaults

Credential variables are:

- `NGTS_WARDEN_PROFILE`
- `NGTS_WARDEN_ACCESS_TOKEN`
- `NGTS_WARDEN_CLIENT_ID`
- `NGTS_WARDEN_CLIENT_SECRET`
- `NGTS_WARDEN_TSG_ID`
- `NGTS_WARDEN_BASE_URL`
- `NGTS_WARDEN_AUTH_URL`

Store profiles in the platform/XDG user config directory as `ngts-warden/config.toml`. `config init` reads secrets through a no-echo prompt and enforces `0600` permissions on Unix. Cache tokens per profile under the user cache directory with `0600` permissions, refreshing 60 seconds before expiry. Never print secrets or full tokens.

### Machine output and MCP

`--json` writes JSON only to stdout. Use a stable envelope:

- Success: `ok`, `operation`, `status`, `data`, and `meta`
- Failure: `ok: false` plus `error.kind`, `message`, optional HTTP status, and request ID
- Error kinds: `usage`, `config`, `auth`, `network`, `api`, `decode`, and `file`
- Exit codes: `0` success, `2` usage/config, `3` auth, `4` transport, `5` API rejection, `6` decode/file failure

Expose five MCP stdio tools using `github.com/modelcontextprotocol/go-sdk`:

- `doctor`
- `list_operations`
- `describe_operation`
- `preview_operation`
- `execute_operation`

MCP starts read-only. Non-GET execution requires both `ngts-warden mcp --allow-write` and `confirm: true` in the tool call. MCP responses use the same redacted JSON envelope as the CLI.

## Implementation and Repository Quality

- Use Cobra `v1.10.2`, MCP Go SDK `v1.7.0`, TOML `v2.4.3`, `golang.org/x/term` for secret input, and the standard HTTP/JSON libraries.
- Keep the Go API internal for v0.x; the supported public contracts are the binary, configuration schema, JSON envelope, generated command names, and MCP tools.
- Check in the pinned OpenAPI snapshot and provenance metadata. Add `make update-openapi` to download, checksum, regenerate commands/help/coverage docs, and show the resulting diff.
- Generate `docs/api-coverage.md` containing all 147 operation IDs, methods, paths, command names, required inputs, response types, and source links.
- Add README quick start, authentication guide, MCP setup, raw-request guidance, API provenance, troubleshooting, and an “unofficial community project” trademark disclaimer.
- Add the Deku-style public repository surface: `LICENSE.md`, `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `CHANGELOG.md`, `.editorconfig`, `.gitattributes`, `.gitignore`, CODEOWNERS, issue forms, PR template, Dependabot, and substantive security/reporting policies.
- Add GitHub Actions for formatting, generated-file drift, tests/race, vet, `golangci-lint`, `govulncheck`, build/smoke, secret scanning, and CodeQL.
- Add GoReleaser configuration for future Linux, macOS, and Windows amd64/arm64 archives and checksums. Do not publish a release initially.
- Create `.codex/skills/ngts-warden` using the skill scaffolder, with concise safe-read, preview/write, raw-request, profile, and MCP workflows. Validate it and provide `make install-skill`.
- Provide `make build`, `install-local`, `test`, `race`, `lint`, `generated-check`, `smoke`, `verify`, `update-openapi`, and `clean`.

## Test and Acceptance Plan

- Verify the generated catalog contains exactly 95 paths, 147 operations, 24 groups, and the documented per-method counts, with no duplicate CLI paths.
- Exercise help generation and argument validation for every generated command.
- Use mock OAuth and NGTS servers to test Basic auth, scope encoding, token refresh/cache expiry, profile precedence, timeouts, retries, redaction, and API error handling.
- Test path/query/header serialization, JSON bodies, previews, `--execute`, CSV exports, binary downloads, stdout/stderr separation, and every exit-code category.
- Verify overly permissive credential files are rejected or clearly reported and secrets never appear in help, logs, JSON errors, or preview output.
- Test MCP initialization, tool discovery, operation lookup, previews, structured results, and both write-safety gates.
- Smoke-test the installed binary from `/tmp`: `command -v`, `--help`, `version`, `--json doctor`, operation count, representative generated reads/writes against fixtures, and MCP startup.
- Run `make verify`, `golangci-lint`, `govulncheck`, repository secret/history checks, and a Deku Scrub prepublication scan.
- Make no live NGTS or authentication-service calls during initial development; all API acceptance uses fixtures and local mock servers.
- After a clean final tree, initialize Git with `main`, create `JerrySabor/ngts-warden` as a public GitHub repository, push `main`, and confirm CI passes. Do not create a version tag.

## Assumptions

- NGTS `application(s)` and `team(s)` compatibility fields remain omitted because the official docs state they are ignored by NGTS.
- `/serviceaccount` commands are labeled “Built-In Accounts” to match the NGTS UI.
- Configus Maximus supplies architectural reference only; NGTS Warden will not read or depend on its private configuration files.
- A TUI, container image, public Go SDK, live credential test, and initial `v0.1.0` release are intentionally deferred.
