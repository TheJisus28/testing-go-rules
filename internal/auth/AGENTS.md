# auth

## Purpose
JWT registration and login. Delegates user persistence to the `users` domain.

## Layout
- handler/ — HTTP adapters
- service/ — validation, bcrypt, token generation
- interfaces/ — AuthService
- models/ — register/login DTOs

## Features

### register / login
- **Service:** AuthService — Register, Login
- **Repository:** UserRepository (users) — Create, FindByUsername, FindByEmail, FindPasswordHash
- **Routes:** POST /v1/auth/register, POST /v1/auth/login
- **Dependencies:** JWT_SECRET, users repository

## Coverage
- service/: 90.7% (target >80%, last run: 2026-05-21)

## Wiring
- DI: config/generals/injector/di.go
- Router: config/generals/router/auth/

## Changelog
- 2026-05-21: service tests with UserRepository mocks
