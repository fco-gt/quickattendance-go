# QuickAttendance - Backend Go

QuickAttendance es una plataforma profesional de gestión de asistencia de empleados diseñada para arquitecturas modernas y escalables. Este repositorio contiene el core del backend migrado íntegramente a **Go**, optimizado para el rendimiento, seguridad multi-tenant y facilidad de despliegue con Docker.

---

> [!NOTE]  
> Esta es la versión en español de la documentación. Para la versión en inglés, consulta [README_ENG.md](./README_ENG.md).

---

### Características Principales

- **Arquitectura Multi-tenant**: Aislamiento total de datos entre diferentes agencias/empresas.
- **Gestión de Horarios Inteligente**: Configuración de turnos con periodos de gracia y asignaciones dinámicas por usuario.
- **Control de Asistencia**: Registro de entradas/salidas con validación de geolocalización y múltiples métodos (QR, NFC, Manual, Teletrabajo).
- **Procesamiento Asíncrono**: Invitaciones y notificaciones gestionadas mediante **RabbitMQ** para garantizar una respuesta rápida de la API y alta disponibilidad.
- **Filtrado Avanzado**: Paginación nativa y búsqueda inteligente en todos los listados (Usuarios, Asistencias, Horarios).
- **Seguridad Robusta**: Autenticación basada en JWT, hashing de contraseñas con bcrypt y control de acceso basado en roles (RBAC).
- **Containerización**: Listo para producción con Docker y Docker Compose.

### Stack Tecnológico

- **Lenguaje**: Go (Golang) 1.25+
- **Framework Web**: Gin Gonic
- **ORM**: GORM (PostgreSQL)
- **Mensajería**: RabbitMQ (AMQP 0.9.1)
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
    Copia el archivo `.env.example` a `.env` y ajusta tus credenciales, incluyendo la URL de `RABBITMQ_URL`.
3.  **Levantar con Docker**:
    ```bash
    docker-compose up --build
    ```
    Esto levantará la API, la base de datos PostgreSQL y el broker de RabbitMQ.
4.  **Ejecutar el Worker**:
    Para procesar las colas de mensajes (como el envío de emails de invitación):
    ```bash
    go run cmd/worker/main.go
    ```

El servidor principal estará disponible en `http://localhost:8080`.

### Documentación Adicional

- [Guía de Pruebas de API (Step-by-Step)](./API_TESTING.md)
- [Esquema de Base de Datos](./docs/database_schema.md)

---
