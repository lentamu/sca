# SCA

## Prerequisites

- Go (1.24+)
- Docker (with Docker Compose)
- golang-migrate (`go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)

## Start the Application

- `Apply migrations:`

```bash
migrate -database "mysql://user:password@tcp(host:port)/database" -path ./migrations up
```

- `Start app:`

```bash
docker compose up --build
```
