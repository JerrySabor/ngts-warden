SHELL := /bin/sh

GO ?= go
BINARY ?= bin/ngts-warden
OPENAPI_URL ?= https://raw.githubusercontent.com/PaloAltoNetworks/pan.dev/ea93c24d4ef2ff55e4efd195e6a2739cf7fb2ad0/openapi-specs/scm/config/ngts/tlsprotect-cloud.json

.DEFAULT_GOAL := help
.PHONY: help build install-local test race vet fmt-check mod-check generated-check coverage-check lint vuln smoke check verify update-openapi install-skill clean

help:
	@printf '%s\n' \
		'ngts-warden developer commands' \
		'' \
		'  make build            Build bin/ngts-warden' \
		'  make install-local    Install the binary under ~/.local/bin' \
		'  make check            Run formatting, modules, tests, race, vet, generation, and smoke checks' \
		'  make verify           Run the full local verification gate' \
		'  make update-openapi   Refresh the pinned OpenAPI snapshot and coverage document' \
		'  make install-skill    Install the repo-local Codex skill' \
		'  make clean            Remove local build output'

build:
	@mkdir -p "$(dir $(BINARY))"
	$(GO) build -trimpath -ldflags "-s -w" -o "$(BINARY)" ./cmd/ngts-warden

install-local: build
	@mkdir -p "$${XDG_BIN_HOME:-$(HOME)/.local/bin}"
	install -m 0755 "$(BINARY)" "$${XDG_BIN_HOME:-$(HOME)/.local/bin}/ngts-warden"

test:
	$(GO) test ./...

race:
	$(GO) test -race ./...

vet:
	$(GO) vet ./...

fmt-check:
	@files="$$(gofmt -l cmd internal)"; test -z "$$files" || { printf '%s\n' 'Go files need formatting:' "$$files"; exit 1; }

mod-check:
	$(GO) mod verify
	$(GO) mod tidy -diff

coverage-check:
	python3 scripts/generate_coverage.py --check

generated-check:
	cmp -s api/openapi.json internal/catalog/spec.json || { printf '%s\n' 'embedded OpenAPI snapshot is stale; run make update-openapi'; exit 1; }
	$(MAKE) coverage-check

lint:
	golangci-lint run ./...

vuln:
	@if command -v govulncheck >/dev/null 2>&1; then govulncheck ./...; else printf '%s\n' 'govulncheck not installed; skipping'; fi

smoke: build
	"$(BINARY)" --help >/dev/null
	"$(BINARY)" version >/dev/null
	"$(BINARY)" --json doctor >/dev/null
	test "$$("$(BINARY)" operations list | wc -l | tr -d ' ')" -eq 147

check: fmt-check mod-check test race vet generated-check build smoke

verify: check lint vuln

update-openapi:
	curl -fsSL "$(OPENAPI_URL)" -o api/openapi.json
	sha256sum api/openapi.json
	cp api/openapi.json internal/catalog/spec.json
	python3 scripts/generate_coverage.py

install-skill:
	@mkdir -p "$${CODEX_HOME:-$(HOME)/.codex}/skills"
	rm -rf "$${CODEX_HOME:-$(HOME)/.codex}/skills/ngts-warden"
	cp -R .codex/skills/ngts-warden "$${CODEX_HOME:-$(HOME)/.codex}/skills/ngts-warden"

clean:
	rm -rf bin coverage.out
