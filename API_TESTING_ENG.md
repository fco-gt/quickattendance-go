# API Testing Guide - QuickAttendance

This guide details the steps to comprehensively test the QuickAttendance API, from agency creation to attendance management.

## Initial Postman Setup
* **Base URL**: `http://localhost:8080/api/v1`
* **Public Postman Collection**: [View collection on Postman](https://www.postman.com/fco-gt/quickattendance/collection/32287192-4c116f57-2c57-4903-b835-34a4e7911073/?action=share&creator=32287192&active-environment=32287192-04a5f77e-97db-4782-996e-24692f0b3443)
* **Env**: Create a `token` variable to store the JWT received upon login.

---

## Step-by-Step Testing Flow

### 1. System Health
*   **Method**: `GET`
*   **URL**: `/health`
*   **Purpose**: Verify that the API is running.

### 2. Agency Registration (Initial Admin)
*   **Method**: `POST`
*   **URL**: `/agencies`
*   **Payload**:
    ```json
    {
      "name": "My Great Company",
      "domain": "company.com",
      "address": "123 Fake Street",
      "phone": "+123456789",
      "admin_email": "admin@company.com",
      "password": "YOUR_PASSWORD"
    }
    ```
*   **Note**: This endpoint creates the agency and the first user with the `admin` role.

### 3. Login
*   **Method**: `POST`
*   **URL**: `/users/login`
*   **Payload**:
    ```json
    {
      "email": "admin@company.com",
      "password": "YOUR_PASSWORD"
    }
    ```
*   **Action**: Copy the `token` from the response and use it for subsequent requests in the `Authorization: Bearer <TOKEN>` header.

### 4. Employee Invitation (Admin Only)
*   **Method**: `POST`
*   **URL**: `/users/invite`
*   **Header**: `Authorization: Bearer <ADMIN_TOKEN>`
*   **Payload**:
    ```json
    {
      "email": "employee@company.com",
      "first_name": "John",
      "last_name": "Doe"
    }
    ```
*   **Note**: In development, the email is not physically sent. You must have the **Worker** running (`go run cmd/worker/main.go`) to see the invitation message in the console, or find the `activation_code` directly in the `users` table of the database.

### 5. Account Activation (Employee)
*   **Method**: `POST`
*   **URL**: `/users/activate`
*   **Payload**:
    ```json
    {
      "activation_token": "YOUR_ACTIVATION_TOKEN",
      "password": "YOUR_NEW_PASSWORD",
      "profile": {
        "first_name": "John",
        "last_name": "Doe Updated"
      }
    }
    ```

### 6. Schedule Management (Admin Only)
#### Create Schedule
*   **Method**: `POST`
*   **URL**: `/schedules`
*   **Payload**:
    ```json
    {
      "name": "Morning Shift",
      "days_of_week": [1, 2, 3, 4, 5],
      "entry_time_minutes": 540, 
      "exit_time_minutes": 1080,
      "grace_period_minutes": 15,
      "is_default": true
    }
    ```
    *(540 min = 09:00 AM, 1080 min = 18:00 PM)*

### 7. Mark Attendance (Employee)
#### Check-In
*   **Method**: `POST`
*   **URL**: `/attendance/mark`
*   **Payload**:
    ```json
    {
      "type": "in",
      "method": "qr",
      "latitude": -34.6037,
      "longitude": -58.3816,
      "notes": "Arriving at the office"
    }
    ```

#### Check-Out
*   **Method**: `POST`
*   **URL**: `/attendance/mark`
*   **Payload**:
    ```json
    {
      "type": "out",
      "method": "manual",
      "notes": "End of shift"
    }
    ```

#### Attendance Business Rules:
*   **Automatic (QR/NFC)**: Employees can mark their own attendance.
*   **Remote (is_remote)**: If `is_remote` is true, the API validates that the coordinates are within the radius configured in the user's profile (`HomeLatitude`, `HomeLongitude`).
*   **Manual**: Only administrators can mark attendance manually for other users. If an employee attempts to use this method, an error will be returned.

### 8. Dynamic Queries (Search and Pagination)
#### List Users with Filters
*   **URL**: `/users/list?search=John&status=active&page=1&limit=10`
#### List Attendance by Date
*   **URL**: `/attendance/list?start_date=2024-01-01&end_date=2024-12-31&status=late`

---

## Roles and Permissions Summary

| Endpoint | Method | User | Admin |
| :--- | :--- | :---: | :---: |
| `/agencies` | POST | 🔓 Public | 🔓 |
| `/users/login` | POST | 🔓 Public | 🔓 |
| `/users/me` | GET | ✅ | ✅ |
| `/users/invite` | POST | ❌ | ✅ |
| `/schedules` | POST/PUT | ❌ | ✅ |
| `/attendance/mark` | POST | ✅ | ✅ |
| `/attendance/list` | GET | ✅ (Own only) | ✅ (Entire agency) |

---
