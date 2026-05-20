# SocialNet — mini social network in Go

REST API with user CRUD, profiles, public/private posts, and friend requests.

## Origin

This repository started from a single prompt:

> In this repository, build a mini social network with user CRUD, public and private posts, user profiles, friend requests, and related features. Bootstrap the Go project and proceed.

The codebase was generated and structured according to the Cursor rules under [`.cursor/rules/`](.cursor/rules/) — including project layout, domain modules, error handling, logging, DI wiring, Swagger, mockery, service tests (>80% coverage), and per-module `AGENTS.md` files.

## Requirements

- Go 1.22+
- Docker (PostgreSQL)

## Quick start

```bash
cp .env.example .env
make db-up
make deps
make generate   # mocks (mockery) + OpenAPI (swag)
make run
```

Interactive API docs: `http://localhost:8080/swagger/index.html`

### Useful commands

| Command | Description |
|---------|-------------|
| `make mocks` | Generate mocks from `.mockery.yaml` (default: local binary) |
| `make mocks MOCKERY=mockery` | Use the `mockery` binary from PATH |
| `make swagger` | Generate `docs/` with swag |
| `make generate` | Run mocks + swagger |
| `make test-cover` | Service-layer test coverage per domain |

The API listens on `http://localhost:8080`.

## Authentication

Register and login return a JWT. Send it on protected routes:

```
Authorization: Bearer <token>
```

## Main endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/auth/register` | Register |
| POST | `/v1/auth/login` | Login |
| GET | `/v1/users` | List profiles |
| GET | `/v1/users/:id/profile` | Public profile |
| PATCH | `/v1/me/profile` | Update my profile (auth) |
| POST | `/v1/posts` | Create post (auth) |
| GET | `/v1/posts/feed` | Personalized feed (auth) |
| GET | `/v1/users/:userId/posts` | User wall |
| POST | `/v1/friendships/requests` | Send friend request (auth) |
| POST | `/v1/friendships/requests/:id/accept` | Accept request (auth) |
| GET | `/v1/friendships/friends` | Friends list (auth) |

## Post visibility

- **public**: visible to everyone.
- **private**: only the author and accepted friends.

## Structure

Follows the project layout: `cmd/`, `config/`, `internal/<domain>/`, `pkg/`.
