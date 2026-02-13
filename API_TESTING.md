# Guía de Pruebas de API - QuickAttendance

Esta guía detalla los pasos para probar de manera integral la API de QuickAttendance, desde la creación de una agencia hasta la gestión de asistencias.

## Configuración Inicial en Postman
* **Base URL**: `http://localhost:8080/api/v1`
* **Colección Pública de Postman**: [Ver colección en Postman](https://www.postman.com/fco-gt/quickattendance/collection/32287192-871aaa97-840d-47f8-8d36-bf7e6f281611)
* **Env**: Crea una variable `token` para almacenar el JWT recibido en el login.

---

## Flujo de Pruebas Paso a Paso

### 1. Salud del Sistema
*   **Método**: `GET`
*   **URL**: `/health`
*   **Propósito**: Verificar que la API esté corriendo.

### 2. Registro de Agencia (Admin Inicial)
*   **Método**: `POST`
*   **URL**: `/agencies`
*   **Payload**:
    ```json
    {
      "name": "Mi Gran Empresa",
      "domain": "empresa.com",
      "address": "Calle Falsa 123",
      "phone": "+123456789",
      "admin_email": "admin@empresa.com",
      "password": "YOUR_PASSWORD"
    }
    ```
*   **Nota**: Este endpoint crea la agencia y al primer usuario con rol `admin`.

### 3. Login
*   **Método**: `POST`
*   **URL**: `/users/login`
*   **Payload**:
    ```json
    {
      "email": "admin@empresa.com",
      "password": "YOUR_PASSWORD"
    }
    ```
*   **Acción**: Copia el `token` de la respuesta y úsalo para las siguientes peticiones en el header `Authorization: Bearer <TOKEN>`.

### 4. Invitación de Empleado (Admin Only)
*   **Método**: `POST`
*   **URL**: `/users/invite`
*   **Header**: `Authorization: Bearer <ADMIN_TOKEN>`
*   **Payload**:
    ```json
    {
      "email": "empleado@empresa.com",
      "first_name": "Juan",
      "last_name": "Pérez"
    }
    ```
*   **Nota**: En desarrollo, el correo no se envía físicamente. Debes tener corriendo el **Worker** (`go run cmd/worker/main.go`) para ver el mensaje de invitación en la consola, o buscar el `activation_code` directamente en la base de datos de la tabla `users`.

### 5. Activación de Cuenta (Empleado)
*   **Método**: `POST`
*   **URL**: `/users/activate`
*   **Payload**:
    ```json
    {
      "activation_token": "YOUR_ACTIVATION_TOKEN",
      "password": "YOUR_NEW_PASSWORD",
      "profile": {
        "first_name": "Juan",
        "last_name": "Pérez Updated"
      }
    }
    ```

### 6. Gestión de Horarios (Admin Only)
#### Crear Horario
*   **Método**: `POST`
*   **URL**: `/schedules`
*   **Payload**:
    ```json
    {
      "name": "Turno Mañana",
      "days_of_week": [1, 2, 3, 4, 5],
      "entry_time_minutes": 540, 
      "exit_time_minutes": 1080,
      "grace_period_minutes": 15,
      "is_default": true
    }
    ```
    *(540 min = 09:00 AM, 1080 min = 18:00 PM)*

### 7. Registrar Asistencia (Empleado)
#### Marcar Entrada
*   **Método**: `POST`
*   **URL**: `/attendance/mark`
*   **Payload**:
    ```json
    {
      "type": "in",
      "method": "qr",
      "latitude": -34.6037,
      "longitude": -58.3816,
      "notes": "Llegando a la oficina"
    }
    ```

#### Marcar Salida
*   **Método**: `POST`
*   **URL**: `/attendance/mark`
*   **Payload**:
    ```json
    {
      "type": "out",
      "method": "manual",
      "notes": "Fin de jornada"
    }
    ```

#### Reglas de Negocio para Asistencia:
*   **Automático (QR/NFC)**: El empleado puede marcar su propia asistencia.
*   **Teletrabajo (is_remote)**: Si `is_remote` es true, la API valida que las coordenadas estén dentro del radio configurado en el perfil del usuario (`HomeLatitude`, `HomeLongitude`).
*   **Manual**: Solo los administradores pueden marcar asistencia manualmente para otros usuarios. Si un empleado intenta usar este método, recibirá un error.

### 8. Consultas Dinámicas (Búsqueda y Paginación)
#### Listar Usuarios con Filtros
*   **URL**: `/users/list?search=Juan&status=active&page=1&limit=10`
#### Listar Asistencias por Fecha
*   **URL**: `/attendance/list?start_date=2026-02-13&end_date=2026-02-13&status=late`

---

## Resumen de Roles y Permisos

| Endpoint | Método | Usuario | Admin |
| :--- | :--- | :---: | :---: |
| `/agencies` | POST | 🔓 Público | 🔓 |
| `/users/login` | POST | 🔓 Público | 🔓 |
| `/users/me` | GET | ✅ | ✅ |
| `/users/invite` | POST | ❌ | ✅ |
| `/schedules` | POST/PUT | ❌ | ✅ |
| `/attendance/mark`| POST | ✅ | ✅ |
| `/attendance/list`| GET | ✅ (Solo propia) | ✅ (Toda la agencia) |

---
