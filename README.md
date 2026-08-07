# NGTS Warden

NGTS Warden is an independent, unofficial alpha proof of concept for working with Palo Alto Networks Next-Gen Trust Security (NGTS). It provides a Go CLI and MCP stdio server for the official NGTS REST API.

The command catalog is generated from a pinned copy of Palo Alto Networks' official OpenAPI specification. It makes the documented API operations discoverable, preserves a raw authenticated request escape hatch, and previews non-GET requests before execution.

> **Project status — Alpha:** NGTS Warden is experimental and not production-ready. Interfaces and behavior may change without notice. Evaluate it with a disposable test tenant and least-privileged credentials; do not use production credentials or customer data.

NGTS Warden is intended to support evaluation and API-driven workflows for this Palo Alto Networks capability. It is not an official Palo Alto Networks product or support offering and is not affiliated with or endorsed by Palo Alto Networks. Palo Alto Networks, NGTS, and related names are trademarks of their respective owners.

## Official resources

- [Next-Gen Trust Security product documentation](https://docs.paloaltonetworks.com/next-gen-trust-security)
- [NGTS API reference](https://pan.dev/scm/api/config/ngts/ngts-api/)

## Install

Build locally with Go 1.26 or newer:

```text
make install-local
ngts-warden --help
```

Future tagged releases will publish cross-platform archives and checksums. The planned distribution is a single static binary; no container runtime is required.

## Configure authentication

While the project is in alpha, use a disposable test tenant and a least-privileged service account. NGTS uses OAuth 2.0 client credentials. Create a service account with the appropriate roles, then configure a named profile:

```text
ngts-warden config init production \
  --client-id "$NGTS_CLIENT_ID" \
  --tsg-id "$NGTS_TSG_ID"
```

The command prompts for the client secret without echoing it. For automation, use environment variables:

```text
export NGTS_WARDEN_CLIENT_ID=...
export NGTS_WARDEN_CLIENT_SECRET=...
export NGTS_WARDEN_TSG_ID=...
export NGTS_WARDEN_PROFILE=production
```

`doctor` never prints credentials:

```text
ngts-warden --json doctor
ngts-warden --json auth status
```

Tokens are requested from `https://auth.apps.paloaltonetworks.com/oauth2/access_token`, cached briefly with restrictive permissions, and refreshed before expiry.

## Usage

Discover operations and exact inputs:

```text
ngts-warden operations list
ngts-warden operations describe certificates_getAll
```

Generated commands follow the resource group and operation ID:

```text
ngts-warden --json certificates get-all
ngts-warden --json certificates get-by-id --id CERTIFICATE_ID
ngts-warden --json machines get-all
```

Every operation's request parameters are generated from OpenAPI. JSON bodies use `--body`, `--body-file`, or stdin. POST, PUT, PATCH, and DELETE commands print a redacted preview until `--execute` is supplied.

Raw reads are available when a high-level command is not convenient:

```text
ngts-warden --json request get /v1/machines
ngts-warden --json request send POST /v1/activitylogsearch --body-file query.json --execute
```

Use `--out FILE` for CSV, text, or binary responses. JSON output is a stable envelope with `ok`, operation, status, data, metadata, and redacted errors.

## MCP

Run the stdio server for an MCP-capable agent:

```text
ngts-warden mcp
```

The server exposes `doctor`, `list_operations`, `describe_operation`, `preview_operation`, and `execute_operation`. It is read-only by default. Enable non-GET execution only when explicitly authorized:

```text
ngts-warden mcp --allow-write
```

## API provenance

The command catalog comes from the official [Palo Alto Networks NGTS API reference](https://pan.dev/scm/api/config/ngts/ngts-api/). The checked-in API snapshot, checksum, upstream commit, and generated operation inventory are in [`api/`](api/) and [`docs/api-coverage.md`](docs/api-coverage.md). Refresh them with:

```text
make update-openapi
```

The upstream specification is MIT licensed. See [`api/provenance.json`](api/provenance.json) and [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

## Development

```text
make check
make verify
```

See [`CONTRIBUTING.md`](CONTRIBUTING.md), [`SECURITY.md`](SECURITY.md), and the repo-local [Codex skill](.codex/skills/ngts-warden/SKILL.md) for contributor and automation guidance.

## License

NGTS Warden is released under the MIT License. See [`LICENSE.md`](LICENSE.md).
