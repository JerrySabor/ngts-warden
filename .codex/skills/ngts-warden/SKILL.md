---
name: ngts-warden
description: Use NGTS Warden, an unofficial alpha Go CLI and MCP stdio server for Palo Alto Networks Next-Gen Trust Security, to inspect NGTS resources, discover pinned OpenAPI operations, preview non-GET requests, manage named OAuth profiles, and execute explicitly authorized API actions. Use for NGTS certificate, machine, endpoint, credential, monitoring, approval, plugin, account, event-log, or raw API work.
---

# NGTS Warden

Use the installed `ngts-warden` command as the primary interface. This proof-of-concept client supports API-driven Palo Alto Networks NGTS workflows and is generated from the repository's pinned official OpenAPI snapshot. It is not official Palo Alto Networks software.

## Workflow

1. Verify installation and configuration:

   ```text
   command -v ngts-warden
   ngts-warden --json doctor
   ```

2. Configure a named profile with `ngts-warden config init NAME`, or provide `NGTS_WARDEN_CLIENT_ID`, `NGTS_WARDEN_CLIENT_SECRET`, `NGTS_WARDEN_TSG_ID`, and `NGTS_WARDEN_PROFILE` through the environment. Never print or copy secrets into prompts, logs, issues, or files.

3. Discover the exact operation before using it:

   ```text
   ngts-warden operations list
   ngts-warden operations describe certificates_getAll
   ```

4. Prefer a generated resource command for reads. Use `--json` when parsing output and use `--out FILE` for CSV, text, or binary responses.

5. Treat every POST, PUT, PATCH, and DELETE as a potentially consequential action. Commands preview the redacted request by default; add `--execute` only when the user explicitly authorized the request.

6. Use `request get PATH` for a safe raw read. Use `request send METHOD PATH` only when no generated command fits, and preserve the same preview/`--execute` guard.

## MCP

Run `ngts-warden mcp` over stdio for an MCP-capable agent. The server exposes `doctor`, `list_operations`, `describe_operation`, `preview_operation`, and `execute_operation`.

MCP is read-only by default. Non-GET execution requires the server flag `--allow-write` and an input `confirm: true`; do not enable either without explicit user authorization.

## Examples

```text
ngts-warden --json certificates get-all
ngts-warden --json machines get-by-id --id MACHINE_ID
ngts-warden --json request get /v1/machines
```

Keep API details and the complete generated inventory in `docs/api-coverage.md`. Do not use live NGTS credentials in tests or include customer payloads, certificates, private keys, bearer tokens, or client secrets in repository artifacts.
