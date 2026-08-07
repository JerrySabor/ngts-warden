# NGTS Warden

NGTS Warden is an unofficial Go CLI and MCP stdio server for exploring and automating the Palo Alto Networks Next-Gen Trust Security (NGTS) API.

It exposes every operation in the pinned official OpenAPI snapshot as a generated, discoverable command while keeping a raw authenticated request escape hatch. Non-GET requests preview by default and require explicit execution.

This project is not affiliated with or endorsed by Palo Alto Networks. Palo Alto Networks, NGTS, and related names are trademarks of their respective owners.

## Install

Build locally with Go 1.26 or newer:

```text
make install-local
ngts-warden --help
```

Future tagged releases publish cross-platform archives and checksums. The v1 distribution is a single static binary; no container runtime is required.

## Configure authentication

NGTS uses OAuth 2.0 client credentials. Create a service account with the appropriate roles, then configure a named profile:

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

The checked-in API snapshot, checksum, upstream commit, and generated operation inventory are in [`api/`](api/) and [`docs/api-coverage.md`](docs/api-coverage.md). Refresh them with:

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
