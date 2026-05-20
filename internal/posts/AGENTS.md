# posts

## Purpose
Posts con visibilidad public (todos) o private (autor y amigos aceptados).

## Layout
- handler/ — HTTP adapters
- service/ — validación, control de visibilidad vía friendships
- repository/ — PostgreSQL con filtros por amistad
- interfaces/ — PostRepository, PostService
- models/ — DTOs

## Features

### posts / feed
- **Service:** PostService — Create, GetByID, Update, Delete, ListByUser, Feed
- **Repository:** PostRepository — CRUD + feed y muro
- **Routes:** POST /v1/posts, GET /v1/posts/:id, PUT/DELETE /v1/posts/:id, GET /v1/posts/feed, GET /v1/users/:userId/posts
- **Dependencies:** Postgres, FriendshipRepository (AreFriends)

## Coverage
- service/: 88.7% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/posts/

## Changelog
- 2026-05-21: tests de servicio (visibilidad, feed, errores de repo)
