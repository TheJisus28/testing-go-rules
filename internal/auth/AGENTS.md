# auth

## Purpose
Registro e inicio de sesión con JWT. Delega persistencia de usuarios al dominio `users`.

## Layout
- handler/ — HTTP adapters
- service/ — validación, bcrypt, generación de token
- interfaces/ — AuthService
- models/ — DTOs de registro/login

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
- 2026-05-21: tests de servicio con mocks de UserRepository
