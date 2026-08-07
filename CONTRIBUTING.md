# Contributing

Thank you for improving NGTS Warden. Keep changes small, documented, and suitable for evaluation with disposable test tenants.

## Before opening a pull request

Run:

```text
make verify
```

This formats and tests the Go code, checks generated API coverage, builds the binary, runs smoke tests, and runs available static checks. Do not use live NGTS credentials in tests; use fixtures or local mock servers.

## API changes

The OpenAPI snapshot is pinned in `api/openapi.json`. Use `make update-openapi` when intentionally refreshing it, review the complete generated diff, and update provenance and coverage documentation together.

## Safety expectations

Never commit customer data, service-account secrets, bearer tokens, certificates, private keys, copied cookies, or full production responses. New write paths must preserve preview and explicit-execution safeguards.
