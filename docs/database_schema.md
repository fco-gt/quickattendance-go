# QuickAttendance - Database Schema

The system uses **PostgreSQL** with **GORM** as the ORM. The schema is designed for multi-tenancy, isolating data by `agency_id`.

## Entities

### Agency
Represents a company or organization using the platform.
- `id`: UUID (Primary Key)
- `name`: String
- `domain`: String (Unique)
- `address`: String
- `phone`: String

### User
Represents an employee or administrator within an agency.
- `id`: UUID (Primary Key)
- `agency_id`: UUID (Foreign Key)
- `first_name`: String
- `last_name`: String (Optional)
- `email`: String (Unique)
- `password_hash`: String
- `role`: Enum (admin, employee)
- `status`: Enum (invited, active, inactive)

### Schedule
Defines the working hours and assigned days for employees.
- `id`: UUID (Primary Key)
- `agency_id`: UUID (Foreign Key)
- `name`: String
- `days_of_week`: String (Comma separated integers 0-6)
- `entry_time_minutes`: Integer (Minutes from start of day)
- `exit_time_minutes`: Integer (Minutes from start of day)
- `grace_period_minutes`: Integer
- `is_default`: Boolean
- **Many-to-Many**: `assigned_users` (via `schedule_users` join table)

### Attendance
Records of employee check-ins and check-outs.
- `id`: UUID (Primary Key)
- `user_id`: UUID (Foreign Key)
- `agency_id`: UUID (Foreign Key)
- `date`: Date
- `check_in_time`: Timestamp
- `check_out_time`: Timestamp (Optional)
- `status`: Enum (present, late, absent)
- `method_in`: Enum (qr, nfc, manual, telework)
- `method_out`: Enum (qr, nfc, manual, telework)
- `latitude`: Float
- `longitude`: Float

---
