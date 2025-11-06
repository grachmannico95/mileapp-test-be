# MileApp Test Backend

This project is a production-ready RESTful API built with the Golang Gin framework, designed for secure user authentication and task management (CRUD) operations. It adheres to Clean Architecture, ensuring scalability, maintainability, and testability

## Overview
A Golang-based REST API that provides:
- User authentication using JWT tokens and secure cookies
- Task management with advanced query options (filtering, sorting, pagination)
- Security features like CSRF protection and password hashing
- Comprehensive testing with mocks using Testify and Mockery

## Design Decision
- Clean Architecture: Provides clear separation of concerns—controllers, services, repositories, and models are decoupled for better maintainability
- JWT + CSRF Double-Submit Cookie Pattern: Ensures both stateless authentication and strong CSRF protection
- Repository Pattern: Abstracts data access, allowing future changes to databases without affecting the business logic
- Dependency Injection: Enables loose coupling between components and improves unit testability

## Strengths of the Module
- Security-Focused: Includes JWT authentication, CSRF validation, bcrypt password hashing, secure cookies and HTTP security headers
- Scalable & Maintainable Architecture: Clean Architecture make it easy to add features or change data sources without breaking existing modules
- Robust Error Handling & Validation: All APIs return consistent error formats with detailed validation feedback

## Databse Indexes
- collection `users`
  - `{ email: 1 }`: Speeds up queries for filtering user based on email, which are common for authentication or account lookups
  - `{ unique: true }`: To prevents duplicate user accounts with the same email
- collecttion `tasks`
  - `{ status: 1 }`: Speeds up queries for filtering task based on status in ascending order
  - `{ priority: 1 }`: Speeds up queries for filtering task based on priority in ascending order
  - `{ due_date: 1 }`: Speeds up queries for filtering task based on due_date in ascending order
  - `{ created_at: -1 }`: Speeds up queries for filtering task based on created_at in descending order

### Setup
- install package
  ```
  go mod tidy
  ```
- create .env file
  ```
  cp .env.example .env
  ```

### How to run
```
go run cmd/main.go
```

### How to build
```
go build -o bin/api ./cmd/api
```

### How to generate mock
```
mockery
```

### How to run unit test
```
go test ./internal/service/... -v
```
