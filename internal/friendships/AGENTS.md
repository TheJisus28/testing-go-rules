# friendships

## Purpose
Friend requests (pending/accepted/rejected) and friends list.

## Layout
- handler/ — HTTP adapters
- service/ — validation and conflict resolution
- repository/ — PostgreSQL
- interfaces/ — FriendshipRepository, FriendshipService
- models/ — DTOs

## Features

### friend requests
- **Service:** FriendshipService — SendRequest, Accept, Reject, ListPendingReceived, ListPendingSent, ListFriends
- **Repository:** FriendshipRepository — friendship CRUD, AreFriends
- **Routes:** POST /v1/friendships/requests, GET requests received/sent, POST accept/reject, GET /v1/friendships/friends
- **Dependencies:** Postgres, UserRepository

## Coverage
- service/: 83.7% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/friendships/

## Changelog
- 2026-05-21: service tests (conflicts, accept/reject, listings)
