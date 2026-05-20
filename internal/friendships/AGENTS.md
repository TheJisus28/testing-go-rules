# friendships

## Purpose
Solicitudes de amistad (pending/accepted/rejected) y lista de amigos.

## Layout
- handler/ — HTTP adapters
- service/ — validación y resolución de conflictos
- repository/ — PostgreSQL
- interfaces/ — FriendshipRepository, FriendshipService
- models/ — DTOs

## Features

### friend requests
- **Service:** FriendshipService — SendRequest, Accept, Reject, ListPendingReceived, ListPendingSent, ListFriends
- **Repository:** FriendshipRepository — CRUD amistades, AreFriends
- **Routes:** POST /v1/friendships/requests, GET requests received/sent, POST accept/reject, GET /v1/friendships/friends
- **Dependencies:** Postgres, UserRepository

## Coverage
- service/: 83.7% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/friendships/

## Changelog
- 2026-05-21: tests de servicio (conflictos, accept/reject, listados)
