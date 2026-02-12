# QuickAttendance - Go Backend

QuickAttendance is a professional employee attendance management platform designed for modern and scalable architectures. This repository contains the backend core fully migrated to **Go**, optimized for performance, multi-tenant security, and easy deployment with Docker.

---

> [!NOTE]  
> This is the English version of the documentation. For the Spanish version, please refer to [README.md](./README.md).

---

### Key Features

- **Multi-tenant Architecture**: Total data isolation between different agencies/companies.
- **Smart Schedule Management**: Shift configuration with grace periods and dynamic per-user assignments.
- **Attendance Control**: Check-in/out records with geolocation validation and multiple methods (QR, NFC, Manual, Remote).
- **Asynchronous Notifications**: Invitation emails processed asynchronously via **RabbitMQ** to ensure high availability and responsiveness.
- **Advanced Filtering**: Native pagination and smart search across all lists (Users, Attendance, Schedules).
- **Robust Security**: JWT-based authentication, bcrypt password hashing, and Role-Based Access Control (RBAC).
- **Containerization**: Production-ready with Docker and Docker Compose.

### Tech Stack

- **Language**: Go (Golang) 1.25+
- **Web Framework**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Messaging**: RabbitMQ (AMQP 0.9.1)
- **Authentication**: JWT (JSON Web Tokens)
- **Logger**: Structured with `slog`
- **Infrastructure**: Docker & Docker Compose

### Architecture Overview

The system follows a modular monolith approach with a clear separation of concerns:
- **Transport Layer**: HTTP Handlers using Gin.
- **Service Layer**: Business logic and orchestration.
- **Domain Layer**: Core entities and repository interfaces.
- **Repository Layer**: Data access implementation (GORM).
- **Messaging Layer**: Asynchronous task producers (RabbitMQ).

### Setup and Usage

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/quickattendance-go.git
    cd quickattendance-go
    ```
2.  **Environment Variables**:
    Copy `.env.example` to `.env` and adjust your credentials, including the `RABBITMQ_URL`.
3.  **Run with Docker**:
    ```bash
    docker-compose up --build
    ```
    This will start the API, the PostgreSQL database, and the RabbitMQ broker.
4.  **Run the Worker**:
    If you are running outside of Docker Compose or need to scale consumers:
    ```bash
    go run cmd/worker/main.go
    ```

The API will be reachable at `http://localhost:8080`.

### Additional Documentation

- [API Testing Guide (Step-by-Step)](./API_TESTING_ENG.md)
- [Database Schema](./docs/database_schema.md)

---
