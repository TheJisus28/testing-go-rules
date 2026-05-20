# Cursor workspace (Go services)

Reusable agent context for Go HTTP services using `cmd/`, `config/`, `internal/<domain>/`, and `pkg/`.

Per-repository bounded contexts are documented in `internal/<domain>/AGENTS.md` inside each project that adopts this folder.

## Rules

Files in `.cursor/rules/` use `alwaysApply: false` and apply when their globs match the files in context.

Read `rules/go-project-structure.mdc` first for layout and naming.

| Rule | Focus |
|------|--------|
| go-project-structure | Repository layout and struct-based design |
| go-config-wiring | DI, routers, logger, storage bootstrap |
| go-module-structure | Handler, service, repository layers |
| go-module-ai-doc | Maintaining `internal/<domain>/AGENTS.md` |
| go-logging | Logger usage by layer |
| go-errors | Domain errors and HTTP mapping |
| go-mockery | Mockery config and generated mocks |
| go-service-tests | Service unit tests (>80% coverage gate) |
| go-api-documentation | Swagger / swag annotations |
| semantic-commits | Conventional commits; confirm before committing |

## Typical Makefile targets

Projects that use this ruleset often expose:

| Target | Purpose |
|--------|---------|
| `mocks` | Regenerate interface mocks |
| `swagger` | Regenerate `docs/` OpenAPI |
| `generate` | mocks + swagger |
| `test-cover` | Coverage on `internal/*/service/` |

Exact targets live in the project `Makefile`; do not assume names if the repo differs.

## Adoption

1. Copy `.cursor/` to the root of the target Go repository.
2. Align `go.mod` module path and `.mockery.yaml` packages with that repo.
3. Add or update `internal/<domain>/AGENTS.md` for each bounded context in **that** repo only.
