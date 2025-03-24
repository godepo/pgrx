# PGRx - Go PostgreSQL Integration Test Extension for groat testing suite

[![codecov](https://codecov.io/gh/godepo/pgrx/graph/badge.svg?token=AC48R7RNVX)](https://codecov.io/gh/godepo/pgrx)
[![Go Report Card](https://goreportcard.com/badge/godepo/pgrx)](https://goreportcard.com/report/godepo/pgrx)
[![License](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/godepo/pgrx/blob/main/LICENSE)

PGRx is a lightweight extension for [groat](https://github.com/godepo/groat) that makes integration testing with 
PostgreSQL in Go clean, simple, and efficient. It leverages 
[testcontainers-go](https://github.com/testcontainers/testcontainers-go) to provide isolated PostgreSQL instances for 
your tests.

## Features
- 🚀 **Simple API** - Bootstrap a PostgreSQL container with minimal code
- 🧪 **Test Isolation** - Each test gets its own database, allowing parallel test execution
- 🔄 **Migrations Support** - Run SQL migrations automatically before tests
- 🔧 **Highly Configurable** - Customize everything from container image to connection pool parameters
- 💉 **Dependency Injection** - Easily inject database connections into your test structures

## Installation
``` bash
go get github.com/godepo/pgrx
```
## Quick Start

For example, we have system under test something like this: 

```go

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
```

Make file main_test.go and add this code in file for support groat conventions:

```go
type State struct{
	// groat states variables
}

// This struct contains fixtured tests dependencies
type Deps struct {
    DB *pgxpool.Pool `groat:"pgxpool"` // inject connection pool to this field after creation
	Config *pgxpool.Config `groat:"pgxconfig"` // inject pgx config to this field after creation
}
```

for work instantiate test add in main_test.go file this code for run postgresql container:

```go
// global state for testing package control all runed integrated containers by groat test suite convention
var suite *integration.Container[Deps, State, *Repostitory]

// By groat convention contain fixture constructor for each isolated test. For each test constructed self pool, database 
// and run migrations. 
func mainProvider(t *testing.T) *groat.Case[Deps, State, *Instance] {
    tcs := groat.New[Deps, State, *Instance](t, func(t *testing.T, deps Deps) *Repository {
        return New(deps.DB) // take pool from injected dependencies field and pass to SUT constructor
    })
    return tcs
}

// Instantiate package suite with single postgresql container for all test.
func TestMain(m *testing.M) {
    suite = integration.New[Deps, State, *Repository](m, mainProvider,
        pgrx.New[Deps](
            pgrx.WithContainerImage("docker.io/postgres:16"),
            pgrx.WithMigrationsPath("./sql"),
        ),
    )
    os.Exit(suite.Go())
}
```

Make file with your tests cases, by example file cases_test.go, and start writing integration tests:

```go
func TestRepository_Create(t *testing.T) {
    tcs := suite.Case(t) // take new instantiated groat test case with constructed SUT (Repository) and worked *Pool.
	
	// write your tests case and use groat and pgrx like fixture, or write groat-like tests below. 
}
```

## How It Works
Under the hood, PGRx:
1. Starts a PostgreSQL container using testcontainers-go
2. Creates a base database for management operations
3. For each test call:
    - Creates a new database with a unique name in postgresql.
    - Runs migrations on the new database.
    - Creates a connection pool for the new database.
    - Injects the connection pool into your test structure.
    - Injects the configuration structure into your test dependencies.

This approach ensures that each test has its own isolated database, making tests reliable and allowing them to run in parallel.

## Requirements
- Go 1.18+ (for generics support)
- Docker (for running the PostgreSQL containers)

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
