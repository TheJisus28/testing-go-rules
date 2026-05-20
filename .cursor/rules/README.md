# Cursor rules (`.mdc`)

Copy `*.mdc` to the service repo: `.cursor/rules/`

| Rule | Scope |
|------|--------|
| [go-project-structure](./go-project-structure.mdc) | Global layout |
| [go-config-wiring](./go-config-wiring.mdc) | `config/`, `cmd/` |
| [go-module-structure](./go-module-structure.mdc) | `internal/<domain>/` |
| [go-module-ai-doc](./go-module-ai-doc.mdc) | `AGENTS.md` |
| [go-logging](./go-logging.mdc) | Logger |
| [go-errors](./go-errors.mdc) | Errors |
| [go-mockery](./go-mockery.mdc) | Mockery |
| [go-service-tests](./go-service-tests.mdc) | Service tests |
| [go-api-documentation](./go-api-documentation.mdc) | Swagger |

**Claude Code:** use [../claude/](../claude/) — condensed `.md` rules with `paths` frontmatter (lower token use).

Template: [../templates/AGENTS.md](../templates/AGENTS.md)
