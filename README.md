# Weaver Edge

[English](README.md) | [Русский](README.ru.md)

> [!WARNING]
> **Alpha version — under active development.**
>
> Weaver Edge is an experimental, incomplete system.
> Telegram integration is currently under development and is not yet fully functional.
> APIs, interfaces, and internal components may change without backward compatibility.
>
> This version is intended for development, testing, and architectural evaluation,
> not for production use.

**Weaver Edge** is the external integration and communication layer of the Weaver system.

It receives messages from external platforms, preserves message data, manages contextual information, applies policies, and communicates with Weaver Core for further processing.

Edge is designed to keep external integrations separate from business orchestration and decision-making.

## Architecture

Weaver Edge is responsible for:

- Receiving messages from external platforms.
- Integrating with Telegram.
- Persisting message-related data in PostgreSQL.
- Providing contextual information.
- Applying policy-related logic.
- Communicating with Weaver Core through an HTTP client.
- Recording audit events.

Weaver Core is responsible for business orchestration and decision-making. Its implementation is maintained separately and is not part of this repository.

For more details, see [architecture.md](architecture.md).

## Message Processing

The general integration flow:

1. A message arrives from an external platform.
2. Edge receives and processes the incoming message.
3. Relevant data is persisted and contextual information is managed.
4. Edge communicates with Weaver Core.
5. The result is handled by the application layer.
6. Relevant actions are recorded in the audit system.

The original message and its context are important for traceability and explainability.

## Technology Stack

- **Go** — application implementation.
- **PostgreSQL** — persistent storage.
- **Docker Compose** — local development and deployment.
- **Gorilla Mux** — HTTP routing.
- **Telegram Bot API** — Telegram integration.

## Project Structure

```text
weaver-edge/
├── CONTRIBUTING.MD
├── Dockerfile
├── LICENSE
├── README.md
├── README.ru.md
├── architecture.md
├── docker-compose.yml
├── docker-compose.prod.yml
├── go.mod
├── go.sum
├── config/
│   └── config.go
├── cmd/
│   ├── edge/
│   │   └── main.go
│   └── migrate/
│       └── main.go
├── internal/
│   ├── app/
│   ├── application/
│   ├── domain/
│   ├── infrastructure/
│   ├── ports/
│   └── transport/
├── migrations/
│   ├── 20261101032526_create_messages.sql
│   ├── 20261101032552_create_policies.sql
│   ├── 20261101032619_create_audit_events.sql
│   └── 20261101032632_create_pending_requests.sql
└── nginx/
    └── nginx.conf
```

## Planned Integrations

Telegram is the first external platform being integrated into Weaver Edge.

The architecture is intended to support additional integrations in the future,
allowing Weaver to receive and process messages from different communication platforms.

Potential future integrations may include:

- Slack
- Discord
- Microsoft Teams
- Other messaging and collaboration platforms

These integrations are planned and are not currently implemented.

## Requirements

- Go (version specified in `go.mod`).
- Docker and Docker Compose.
- PostgreSQL (when running without Docker).
- A Telegram bot token for Telegram integration.
- Access to a running Weaver Core instance.

## Configuration

Edge reads configuration from environment variables.

| Variable | Description |
|---|---|
| `HTTP_ADDR` | HTTP server listen address |
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_DB` | PostgreSQL database name |
| `POSTGRES_HOST` | PostgreSQL host |
| `POSTGRES_PORT` | PostgreSQL port |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token |
| `CORE_API_URL` | Weaver Core API base URL |
| `CORE_API_TOKEN` | Authentication token for Weaver Core |

Create a local `.env` file based on `.env.example`.

Do not commit real credentials or tokens.

## Local Development

The repository includes a Docker Compose configuration for local development.

Create the external Docker network if it does not already exist:

```bash
docker network create weaver-local
```

Configure the required environment variables in `.env`, then start the application:

```bash
docker compose up -d --build
```

The local configuration uses the `weaver-local` Docker network to communicate with other Weaver services.

For local communication with Weaver Core, configure `CORE_API_URL` to use the Core service's Docker network address.

## API

The HTTP router currently exposes the following endpoints:

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/webhooks/telegram` | Telegram webhook endpoint |
| POST | `/context` | Context-related endpoint |

The exact request and response contracts are defined by the HTTP handlers and DTOs in the source code.

## Database Migrations

Database migrations are stored in the `migrations/` directory.

The repository includes a separate migration entry point:

```text
    # Apply migrations
    docker compose --profile manual run --rm migrate

    # Roll back one migration
    docker compose --profile manual run --rm migrate down

    # Show migration status
    docker compose --profile manual run --rm migrate status
```

The local Docker Compose configuration includes a migration service that can be run using the `manual` profile.

## Testing

Run the Go test suite:

```bash
go test ./...
```

Some integration tests require an available PostgreSQL instance and correctly configured environment variables.

## Deployment

The repository includes separate Compose configurations:

- `docker-compose.yml` — local development.
- `docker-compose.prod.yml` — production deployment.

An Nginx configuration is also provided in `nginx/`.

Production deployment uses HTTPS and keeps internal services isolated from direct public access.

## Project Status

Weaver Edge is under active development.

The project is evolving alongside the broader Weaver architecture. Interfaces and implementation details may change.

## License

See [LICENSE](LICENSE).