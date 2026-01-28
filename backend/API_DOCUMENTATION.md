# Anytime Travel API Documentation

## Base URL
```
http://localhost:8080/api
```

## App User Authentication Endpoints

### 1. User Signup

**Endpoint:** `POST /api/app/signup`

**Description:** Register a new user account

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",
  "password": "securePassword123",
  "confirm_password": "securePassword123",
  "sex": "Male",
  "country": "United States"
}
```

**Field Validations:**
- `name`: Required, full name of the user
- `email`: Required, valid email format, must be unique
- `phone_number`: Required, must be unique
- `password`: Required, minimum 6 characters
- `confirm_password`: Required, must match password
- `sex`: Required, must be one of: "Male", "Female", "Other"
- `country`: Required, country of residence

**Success Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "sex": "Male",
    "country": "United States",
    "is_active": true,
    "created_at": "2026-01-27T10:30:00Z"
  }
}
```

**Error Responses:**

- **400 Bad Request:** Invalid request format or missing required fields
```json
{
  "error": "All fields are required"
}
```

- **400 Bad Request:** Passwords don't match
```json
{
  "error": "Passwords do not match"
}
```

- **400 Bad Request:** Password too short
```json
{
  "error": "Password must be at least 6 characters"
}
```

- **400 Bad Request:** Invalid sex value
```json
{
  "error": "Invalid sex value. Must be Male, Female, or Other"
}
```

- **409 Conflict:** Email already registered
```json
{
  "error": "Email already registered"
}
```

- **409 Conflict:** Phone number already registered
```json
{
  "error": "Phone number already registered"
}
```

---

### 2. User Login

**Endpoint:** `POST /api/app/login`

**Description:** Authenticate a user and receive JWT token

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Field Validations:**
- `email`: Required, valid email format
- `password`: Required

**Success Response (200 OK):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "sex": "Male",
    "country": "United States",
    "is_active": true,
    "created_at": "2026-01-27T10:30:00Z"
  }
}
```

**Error Responses:**

- **400 Bad Request:** Missing credentials
```json
{
  "error": "Email and password are required"
}
```

- **401 Unauthorized:** Invalid credentials
```json
{
  "error": "Invalid email or password"
}
```

- **403 Forbidden:** Account suspended
```json
{
  "error": "Your account has been suspended. Please contact support."
}
```

---

## Authentication

After successful signup or login, include the JWT token in the `Authorization` header for protected endpoints:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Token Details:**
- Expiration: 7 days
- Algorithm: HS256
- Claims: user_id, email, name, exp, iat

---

## Database Schema

### Users Collection

```javascript
{
  "_id": "123e4567-e89b-12d3-a456-426614174000",  // UUID
  "name": "John Doe",
  "email": "john.doe@example.com",  // Lowercase, unique
  "phone_number": "+1234567890",  // Unique
  "password_hash": "$2a$10$...",  // Bcrypt hash
  "sex": "Male",  // Male, Female, or Other
  "country": "United States",
  "is_active": true,
  "is_freezed": false,
  "total_bookings": 0,
  "last_booking": ISODate("2026-01-27T10:30:00Z"),  // Optional
  "created_at": ISODate("2026-01-27T10:30:00Z"),
  "last_login": ISODate("2026-01-27T10:30:00Z")
}
```

### Indexes

Recommended MongoDB indexes for optimal performance:

```javascript
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "phone_number": 1 }, { unique: true })
db.users.createIndex({ "created_at": -1 })
db.users.createIndex({ "last_login": -1 })
```

---

## Environment Variables

Required environment variables for the backend:

```bash
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=anytime_travel
JWT_SECRET=your-super-secret-jwt-key-change-in-production
ENVIRONMENT=development
```

---

## Testing with cURL

### Signup Example:
```bash
curl -X POST http://localhost:8080/api/app/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "password": "password123",
    "confirm_password": "password123",
    "sex": "Male",
    "country": "United States"
  }'
```

### Login Example:
```bash
curl -X POST http://localhost:8080/api/app/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Using Token for Protected Endpoints:
```bash
curl -X GET http://localhost:8080/api/app/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

---

## Security Considerations

1. **Password Hashing:** All passwords are hashed using bcrypt before storage
2. **Email Normalization:** Emails are converted to lowercase to prevent duplicate accounts
3. **JWT Tokens:** Use secure, random JWT secrets in production
4. **HTTPS:** Always use HTTPS in production environments
5. **Rate Limiting:** Consider implementing rate limiting for authentication endpoints
6. **Input Validation:** All inputs are validated before processing
7. **Account Freezing:** Frozen accounts cannot log in

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

HTTP Status Codes used:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input or request format
- `401 Unauthorized`: Authentication failed
- `403 Forbidden`: Access denied (e.g., frozen account)
- `409 Conflict`: Resource already exists (duplicate email/phone)
- `500 Internal Server Error`: Server-side error
