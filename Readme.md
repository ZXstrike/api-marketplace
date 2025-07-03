# API Marketplace Platform

An open-source, modular API marketplace platform built with a microservices architecture using Go for the backend, Vue.js for the frontend, and Nginx as a reverse proxy. This platform allows developers to publish their APIs, manage subscriptions, and provides consumers with a central place to discover and subscribe to APIs. It includes features like robust authentication, rate limiting, and billing.

## Table of Contents

- [Architecture](#architecture)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
- [Development](#development)
- [API Documentation](#api-documentation)
- [License](#license)
- [Author](#author)

## Architecture

The platform is designed using a microservices architecture to ensure scalability and maintainability. The core components are:

-   **Marketplace App**: The main backend service responsible for user management, API listings, subscriptions, and billing logic.
-   **API Gateway**: Handles all incoming requests, providing authentication, rate limiting, and routing to the appropriate backend services.
-   **Marketplace Web**: A Vue.js single-page application for the user interface.
-   **Nginx**: Acts as a reverse proxy, directing traffic to the appropriate service.

## Features

-   **Microservices Architecture**: Scalable and maintainable project structure.
-   **User & API Management**: Endpoints for managing users, APIs, and subscriptions.
-   **JWT Authentication**: Secure authentication using ECDSA-based JWT.
-   **Rate Limiting**: Protects your APIs from abuse.
-   **Billing & Subscriptions**: A system for managing API subscription plans and billing.
-   **Containerized**: Easy to deploy and manage with Docker.

## Technologies Used

-   **Backend**: Go
-   **Frontend**: Vue.js, Vite
-   **Database**: PostgreSQL, Redis
-   **Reverse Proxy**: Nginx
-   **Containerization**: Docker, Docker Compose

## Project Structure

```
.
├── api-gateway/         # API Gateway service
├── marketplace-app/     # Main marketplace backend service
├── marketplace-web/     # Frontend Vue.js application
├── nginx/               # Nginx configuration
├── shared/              # Shared Go libraries (JWT, models, etc.)
├── project-docs/        # Project documentation
├── docker-compose.yml   # Docker orchestration
├── makefile             # Development commands
└── ...
```

## Getting Started

### Prerequisites

-   [Go](https://golang.org/dl/) (1.24+ recommended)
-   [Node.js](https://nodejs.org/) (for frontend development)
-   [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
-   [Make](https://www.gnu.org/software/make/) (optional, for convenience)

### Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/ZXstrike/api-marketplace.git
    cd api-marketplace
    ```

2.  **Configure environment variables:**
    Copy the example `.env` file and modify it if needed.
    ```sh
    cp .env.example .env
    ```

3.  **Generate ECDSA keys for JWT:**
    ```sh
    make gen-keys
    ```

4.  **Run the entire platform with Docker Compose:**
    This is the recommended way to run the project for development.
    ```sh
    docker-compose up --build
    ```
    The services will be available at:
    -   **Frontend**: `http://localhost:5173`
    -   **API Gateway**: `http://localhost:8080`
    -   **Marketplace App**: `http://localhost:8081`


## Development

### Running Services Individually

If you prefer not to use Docker Compose, you can run each service manually.

-   **Start Databases & Services:**
    ```sh
    make dev-up
    ```

-   **Stop Databases & Services:**
    ```sh
    make dev-down
    ```

-   **Run Marketplace App:**
    ```sh
    make run-market
    ```

-   **Run API Gateway:**
    ```sh
    make run-gateway
    ```

-   **Run Frontend:**
    ```sh
    cd marketplace-web
    npm install
    npm run dev
    ```

### Running in Production Mode

-   **Build production Docker images:**
    ```sh
    make build-prod
    ```

-   **Start all services in production mode:**
    ```sh
    make prod-up
    ```

-   **Stop all production services:**
    ```sh
    make prod-down
    ```

### Running Tests

-   **Run all tests for Go services:**
    ```sh
    make test
    ```

### Database Migration

-   **Apply database migrations:**
    ```sh
    make migrate
    ```

## API Documentation

The API endpoints are documented in `project-docs/endpoint.md`.

## License

This project is for educational purposes.

## Local Development with Subdomain Routing

If you want to use the API Gateway with subdomain-based routing (for example, to support multi-tenant APIs or per-API subdomains), you must access your services via a real domain or DNS entry that supports subdomains. Browsers and HTTP clients only send the correct Host header (and thus subdomain) when using a proper domain, not just localhost or an IP address.

### Setting Up Local DNS for Subdomains

1. **Edit your hosts file** to point a wildcard domain to 127.0.0.1. On Windows, open `C:\Windows\System32\drivers\etc\hosts` as Administrator and add lines like:
   ```
   127.0.0.1 app.test
   127.0.0.1 api.app.test
   # Add more subdomains as needed (wildcards are not supported in Windows hosts file)
   ```

2. **Configure the Vite dev server** to accept these domains. This project already allows `app.test` and `.app.test` in `vite.config.js`.

3. **Access your services** using the domain, for example:
   - Frontend: http://app.test:5173
   - API Gateway: http://api.app.test:8080

4. **For production**, configure your DNS provider to point your domain and wildcard subdomains to your server's IP address.

> **Note:** Subdomain-based routing will not work with just `localhost` or an IP address. You must use a real domain or DNS entry for the API Gateway to properly route requests based on subdomain.

---

**Author:** ZXstrike
