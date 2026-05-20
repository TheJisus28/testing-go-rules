# users

## Purpose
CRUD de cuentas y perfiles públicos de usuario.

## Layout
- handler/ — HTTP adapters
- service/ — reglas de negocio y autorización (solo propio usuario)
- repository/ — PostgreSQL
- interfaces/ — UserRepository, UserService
- models/ — DTOs

## Features

### users / profile
- **Service:** UserService — GetByID, GetProfile, List, Update, UpdateProfile, Delete
- **Repository:** UserRepository — CRUD + listado de perfiles
- **Routes:** GET /v1/users, GET /v1/users/:id, GET /v1/users/:id/profile, PUT /v1/users/:id, DELETE /v1/users/:id, PATCH /v1/me/profile
- **Dependencies:** Postgres pool

## Coverage
- service/: 97.5% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/users/

## Changelog
- 2026-05-21: tests de servicio ampliados (UpdateProfile, errores de repo)
