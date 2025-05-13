# API Marketplace Platform

A modular API marketplace platform built with Go, supporting API publishing, subscription, authentication, rate limiting, and billing. The project is organized into microservices with shared libraries and containerized for easy deployment.

## Project Structure

```
.
├── api-gateway/         # API Gateway service (reverse proxy, auth, rate limiting)
├── marketplace-app/     # Main marketplace backend (users, APIs, billing, etc.)
├── shared/              # Shared libraries (JWT, secrets, etc.)
├── docs/                # Documentation
├── .env                 # Environment variables
├── docker-compose.yml   # Multi-service orchestration
├── Makefile             # Common development commands
└── ...
```

## Features

- **API Gateway**: Reverse proxy, authentication, rate limiting, and request routing.
- **Marketplace App**: User management, API publishing, subscriptions, billing, and analytics.
- **JWT Authentication**: ECDSA-based JWT for secure user and API key authentication.
- **PostgreSQL & Redis**: Robust data storage and caching.
- **Containerized**: Easily deployable with Docker Compose.
- **Extensible**: Modular codebase for adding new features.

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience)

### Setup

1. **Clone the repository**

   ```sh
   git clone https://github.com/ZXstrike/api-marketplace.git
   cd api-marketplace
   ```

2. **Copy and configure environment variables**

   ```sh
   cp .env.example .env
   # Edit .env as needed (DB credentials, ports, etc.)
   ```

3. **Generate ECDSA keys**

   ```sh
   make gen-keys
   ```

4. **Start services (DB, Redis, etc.)**

   ```sh
   make dev-up
   ```

5. **Run the Marketplace App**

   ```sh
   make run-market
   ```

6. **Run the API Gateway**

   ```sh
   make run-gateway
   ```

### Running with Docker Compose

Uncomment the relevant services in `docker-compose.yml` and run:

```sh
docker-compose up --build
```

## Development

- **Run all tests**

  ```sh
  make test
  ```

- **Database migration**

  ```sh
  make migrate
  ```

## Key Directories

- [`api-gateway`](api-gateway/) — Gateway service (see [`cmd/main.go`](api-gateway/cmd/main.go))
- [`marketplace-app`](marketplace-app/) — Main backend (see [`cmd/main.go`](marketplace-app/cmd/main.go))
- [`shared/pkg/jwt`](shared/pkg/jwt/) — JWT utilities ([`generate_token.go`](shared/pkg/jwt/generate_token.go), [`verify_token.go`](shared/pkg/jwt/verify_token.go))
- [`shared/secrets`](shared/secrets/) — ECDSA key generation scripts and keys

## Environment Variables

See [.env.example](.env.example) for all configuration options.

## License

This project is for educational purposes.

---

**Author:** ZXstrike
```
