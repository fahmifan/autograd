# Autograde

Autograd is a web app that provides code autograding for schools and university programming course.
Currently it supports c++ language only.

## Installation
__TBD__

## Design Overview
The project is divided into two main components:
- cmd
- db
- frontend
- pkg
- proto
- testdata

### cmd
The `cmd` folder contains the entry point for the cli application 

### db
The `db` folder contains the database schema and migration file & sqlc queries

### frontend
The `frontend` folder contains the frontend application written in reactjs

### pkg
The backend is written in Go under `pkg` folder and have this several components:
- config: hold the configuration & env for the application
- core: hold the core logic of the application
  - each sub-package represent the domain of the application
    - the domain creation, which part belong to which domain is up to the team and should be based on business process
  - inside sub-package there is `*_cmd` & `*_query` package that is the public API written in connectrpc framework
  - for data persistence, we interfacing using gorm & sqlc. Gorm is initially used, but will migrated to sqlc for better performance
- dbconn: hold the database connection and migration logic
- dbmodel: hold the database model, use to query using `gorm`
- httpsvc: wiring the service and http handler
- fs: hold the file system to store & retrieve the file
- jobqueue: handle jobqueue, it has outbox pattern implemented under `outbox` package
  - outbox: hold the outbox pattern implementation. It is used to do async operation
- logs: hold the logger configuration
- mailer: handle mailing feature, under it is mailing provider implementation
- pb: the generated protobuf file
- service: the wiring for the cmd & query
- xsqlc: the generated sqlc file 

### proto
The `proto` folder contains the protobuf file definition

### testdata
The `testdata` folder contains the test data for assignemnt & submission

## Development
- Prerequisites
  - Golang v1.22+
  - NodeJS v18+
  - pnpm
  - postgres v16
  - [mailhog](https://github.com/mailhog/MailHog)
  - [modd](https://github.com/cortesi/modd)
  - [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
  - [sql-migrate](https://github.com/rubenv/sql-migrate)
- Clone the project
- cd to the project directory
- copy .env file from .env.example
  - `cp .env.example .env`
- run `go mod tidy` to install the backend dependencies
- cd to the `frontend` directory run `pnpm install` to install the frontend dependencies

### Migrate Database
- run `make db-migrate-up`

### Running Backend Service
- run the postgres & mailhog service
- run `make run-server` to run the backend service

### Running the Frontend Service
- cd to the `frontend` directory
- run `pnpm dev` to run the frontend service

### Create Initial Admin
- run this command
  ```bash
  go run cmd/autograd/main.go admin create --email john@doe.com --name "john doe" --password "supersecret"
  ```
