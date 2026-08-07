# Security policy

## Reporting a vulnerability

Do not open a public issue for credentials, token-handling bugs, request smuggling, or other security-sensitive defects. Email the maintainer at `jerrysabor@gmail.com` with a concise description, reproduction steps, affected version, and a safe contact method.

Please do not include real client secrets, bearer tokens, private keys, customer payloads, or production URLs in a report. Redact them before sending.

We will acknowledge a report as soon as practical, assess impact, and coordinate a fix or mitigation. This is a best-effort community project and does not provide a guaranteed response-time SLA.

## Credential safety

Use environment variables or the restrictive local profile file. Never commit `config.toml`, `.env` files, bearer tokens, client secrets, or downloaded private-key material. Non-GET operations preview by default for this reason.
