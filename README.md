# QuickAttendance - Backend Go

QuickAttendance es una plataforma profesional de gestión de asistencia de empleados diseñada para arquitecturas modernas y escalables. Este repositorio contiene el core del backend migrado íntegramente a **Go**, optimizado para el rendimiento, seguridad multi-tenant y facilidad de despliegue con Docker.

---

## 🇪🇸 Versión en Español

### Características Principales

- **Arquitectura Multi-tenant**: Aislamiento total de datos entre diferentes agencias/empresas.
- **Gestión de Horarios Inteligente**: Configuración de turnos con periodos de gracia y asignaciones dinámicas por usuario.
- **Control de Asistencia**: Registro de entradas/salidas con validación de geolocalización y múltiples métodos (QR, NFC, Manual, Teletrabajo).
- **Filtrado Avanzado**: Paginación nativa y búsqueda inteligente en todos los listados (Usuarios, Asistencias, Horarios).
- **Seguridad Robusta**: Autenticación basada en JWT, hashing de contraseñas con bcrypt y control de acceso basado en roles (RBAC).
- **Containerización**: Listo para producción con Docker y Docker Compose.

### Stack Tecnológico

- **Lenguaje**: Go (Golang) 1.25+
- **Framework Web**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Autenticación**: JWT (JSON Web Tokens)
- **Logger**: Estructurado con `slog`
- **Infraestructura**: Docker & Docker Compose

### Instalación y Uso

1.  **Clonar el repositorio**:
    ```bash
    git clone https://github.com/tu-usuario/quickattendance-go.git
    cd quickattendance-go
    ```
2.  **Configurar Variables de Entorno**:
    Copia el archivo `.env.example` a `.env` y ajusta tus credenciales.
3.  **Levantar con Docker**:
    ```bash
    docker-compose up --build
    ```
    El servidor estará disponible en `http://localhost:8080`.

---

## 🇺🇸 English Version

### Key Features

- **Multi-tenant Architecture**: Total data isolation between different agencies/companies.
- **Smart Schedule Management**: Shift configuration with grace periods and dynamic per-user assignments.
- **Attendance Control**: Check-in/out records with geolocation validation and multiple methods (QR, NFC, Manual, Telework).
- **Advanced Filtering**: Native pagination and smart search across all lists (Users, Attendance, Schedules).
- **Robust Security**: JWT-based authentication, bcrypt password hashing, and Role-Based Access Control (RBAC).
- **Containerization**: Production-ready with Docker and Docker Compose.

### Tech Stack

- **Language**: Go (Golang) 1.22+
- **Web Framework**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Authentication**: JWT (JSON Web Tokens)
- **Logger**: Structured with `slog`
- **Infrastructure**: Docker & Docker Compose

### Setup and Usage

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/quickattendance-go.git
    cd quickattendance-go
    ```
2.  **Environment Variables**:
    Copy `.env.example` to `.env` and adjust your credentials.
3.  **Run with Docker**:
    ```bash
    docker-compose up --build
    ```
    The server will be reachable at `http://localhost:8080`.

---

## Documentación Adicional

- [Guía de Pruebas de API (Step-by-Step)](./API_TESTING.md)
- [Esquema de Base de Datos](./docs/database_schema.md)

---
