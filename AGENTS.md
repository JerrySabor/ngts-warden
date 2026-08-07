# NGTS Warden repository guidance

- Use Go 1.26 or newer and keep the CLI usable as a standalone binary.
- Treat `api/openapi.json` as the source snapshot. Refresh it only with `make update-openapi`, then review the generated coverage diff.
- Never commit client secrets, bearer tokens, certificates, private keys, customer payloads, copied cookies, or live API responses.
- Preserve preview-by-default behavior for all non-GET operations and the MCP `--allow-write` plus `confirm` gates.
- Run `make verify` before handing off a change.
