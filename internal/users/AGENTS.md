# users

## Purpose
User account CRUD and public profiles.

## Layout
- handler/ — HTTP adapters
- service/ — business rules and authorization (own account only)
- repository/ — PostgreSQL
- interfaces/ — UserRepository, UserService
- models/ — DTOs

## Features

### users / profile
- **Service:** UserService — GetByID, GetProfile, List, Update, UpdateProfile, Delete
- **Repository:** UserRepository — CRUD + profile listing
- **Routes:** GET /v1/users, GET /v1/users/:id, GET /v1/users/:id/profile, PUT /v1/users/:id, DELETE /v1/users/:id, PATCH /v1/me/profile
- **Dependencies:** Postgres pool

## Coverage
- service/: 97.5% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/users/

## Changelog
- 2026-05-21: expanded service tests (UpdateProfile, repository errors)
