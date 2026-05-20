# SocialNet — mini red social en Go

API REST con CRUD de usuarios, perfiles, posts públicos/privados y solicitudes de amistad.

## Requisitos

- Go 1.22+
- Docker (PostgreSQL)

## Inicio rápido

```bash
cp .env.example .env
make db-up
make deps
make generate   # mocks (mockery) + OpenAPI (swag)
make run
```

Documentación interactiva: `http://localhost:8080/swagger/index.html`

### Comandos útiles

| Comando | Descripción |
|---------|-------------|
| `make mocks` | Genera mocks desde `.mockery.yaml` (por defecto vía `go run`) |
| `make mocks MOCKERY=mockery` | Usa el binario `mockery` instalado en PATH |
| `make swagger` | Genera `docs/` con swag |
| `make generate` | Ejecuta mocks + swagger |
| `make test-cover` | Cobertura de tests en `service/` |

La API escucha en `http://localhost:8080`.

## Autenticación

Registro y login devuelven un JWT. Envíalo en las rutas protegidas:

```
Authorization: Bearer <token>
```

## Endpoints principales

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/v1/auth/register` | Registro |
| POST | `/v1/auth/login` | Login |
| GET | `/v1/users` | Listar perfiles |
| GET | `/v1/users/:id/profile` | Perfil público |
| PATCH | `/v1/me/profile` | Editar mi perfil (auth) |
| POST | `/v1/posts` | Crear post (auth) |
| GET | `/v1/posts/feed` | Feed personalizado (auth) |
| GET | `/v1/users/:userId/posts` | Posts de un usuario |
| POST | `/v1/friendships/requests` | Enviar solicitud (auth) |
| POST | `/v1/friendships/requests/:id/accept` | Aceptar (auth) |
| GET | `/v1/friendships/friends` | Lista de amigos (auth) |

## Visibilidad de posts

- **public**: visible para todos.
- **private**: solo el autor y sus amigos aceptados.

## Estructura

Sigue el layout del proyecto: `cmd/`, `config/`, `internal/<domain>/`, `pkg/`.
