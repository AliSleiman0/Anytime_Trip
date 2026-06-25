# JWT Authentication Middleware

This package provides JWT authentication middleware for the Travel backend API.

## Features

- **JWT Token Validation**: Validates JWT tokens from Authorization header or cookies
- **Role-Based Access Control**: Middleware for admin and superadmin roles
- **Flexible Token Sources**: Supports both Authorization header and cookie-based authentication
- **Helper Functions**: Utility functions to extract user information from context
- **Optional Authentication**: Support for routes that work with or without authentication

## Middleware Functions

### AuthMiddleware

The main authentication middleware that validates JWT tokens.

```go
func AuthMiddleware(jwtSecret string) fiber.Handler
```

**Token Sources** (checked in order):
1. `Authorization` header: `Bearer <token>`
2. `auth_token` cookie
3. `admin_token` cookie

**Response**: Returns 401 Unauthorized if token is missing or invalid

**Usage**:
```go
protected := router.Use(middleware.AuthMiddleware(cfg.JWTSecret))
protected.Get("/dashboard", handler.GetDashboard)
```

### AdminMiddleware

Checks if the authenticated user has admin or superadmin role.

```go
func AdminMiddleware() fiber.Handler
```

**Prerequisites**: Must be used after `AuthMiddleware`

**Response**: Returns 403 Forbidden if user doesn't have admin role

**Usage**:
```go
adminRoutes := router.Use(middleware.AuthMiddleware(cfg.JWTSecret))
adminRoutes.Use(middleware.AdminMiddleware())
adminRoutes.Get("/users", handler.ManageUsers)
```

### SuperAdminMiddleware

Checks if the authenticated user has superadmin role.

```go
func SuperAdminMiddleware() fiber.Handler
```

**Prerequisites**: Must be used after `AuthMiddleware`

**Response**: Returns 403 Forbidden if user doesn't have superadmin role

**Usage**:
```go
superAdminRoutes := router.Use(middleware.AuthMiddleware(cfg.JWTSecret))
superAdminRoutes.Use(middleware.SuperAdminMiddleware())
superAdminRoutes.Get("/system", handler.SystemSettings)
```

### OptionalAuthMiddleware

Attempts to authenticate but doesn't fail if token is missing. Useful for routes that behave differently for authenticated vs unauthenticated users.

```go
func OptionalAuthMiddleware(jwtSecret string) fiber.Handler
```

**Behavior**: Sets user context if valid token exists, otherwise continues without error

**Usage**:
```go
publicRoutes := router.Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
publicRoutes.Get("/products", handler.ListProducts) // Shows different content for authenticated users
```

## Helper Functions

### GetUserID

Extracts the user ID from the request context.

```go
func GetUserID(c *fiber.Ctx) string
```

**Returns**: User ID as string, or empty string if not authenticated

### GetUserRole

Extracts the user role from the request context.

```go
func GetUserRole(c *fiber.Ctx) string
```

**Returns**: User role as string, or empty string if not authenticated

### GetClaims

Extracts the full JWT claims from the request context.

```go
func GetClaims(c *fiber.Ctx) *utils.AdminClaims
```

**Returns**: AdminClaims pointer, or nil if not authenticated

### IsAuthenticated

Checks if the current request is authenticated.

```go
func IsAuthenticated(c *fiber.Ctx) bool
```

### IsAdmin

Checks if the current user has admin or superadmin role.

```go
func IsAdmin(c *fiber.Ctx) bool
```

### IsSuperAdmin

Checks if the current user has superadmin role.

```go
func IsSuperAdmin(c *fiber.Ctx) bool
```

## Usage Examples

### Example 1: Protected Admin Routes

```go
// Setup admin routes with authentication
adminGroup := app.Group("/admin")

// Public routes (no auth required)
adminGroup.Post("/login", handler.HandleLogin)
adminGroup.Post("/forgot-password", handler.HandleForgotPassword)

// Protected routes
protected := adminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
protected.Use(middleware.AdminMiddleware())

protected.Get("/dashboard", handler.GetDashboard)
protected.Get("/users", handler.ManageUsers)
protected.Post("/logout", handler.HandleLogout)
```

### Example 2: Using Helper Functions in Handlers

```go
func (h *Handler) GetProfile(c *fiber.Ctx) error {
    // Get user ID from context
    userID := middleware.GetUserID(c)
    
    // Get user role
    role := middleware.GetUserRole(c)
    
    // Check if user is admin
    if middleware.IsAdmin(c) {
        // Admin-specific logic
    }
    
    // Get full claims for more details
    claims := middleware.GetClaims(c)
    if claims != nil {
        log.Printf("User %s logged in at %v", claims.Subject, claims.IssuedAt)
    }
    
    // Your handler logic here...
    return c.JSON(fiber.Map{
        "user_id": userID,
        "role": role,
    })
}
```

### Example 3: Optional Authentication

```go
// Route that works for both authenticated and unauthenticated users
apiGroup := app.Group("/api")
apiGroup.Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))

apiGroup.Get("/products", func(c *fiber.Ctx) error {
    if middleware.IsAuthenticated(c) {
        // Show personalized products
        userID := middleware.GetUserID(c)
        return c.JSON(getPersonalizedProducts(userID))
    }
    
    // Show general products for unauthenticated users
    return c.JSON(getGeneralProducts())
})
```

### Example 4: SuperAdmin-Only Routes

```go
superAdminGroup := app.Group("/superadmin")
superAdminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
superAdminGroup.Use(middleware.SuperAdminMiddleware())

superAdminGroup.Get("/system-config", handler.GetSystemConfig)
superAdminGroup.Post("/system-config", handler.UpdateSystemConfig)
superAdminGroup.Delete("/users/:id", handler.DeleteUser)
```

## Token Generation

Tokens are generated using the `utils.GenerateJWT` function:

```go
import "travel/backend/core/utils"

token, err := utils.GenerateJWT(
    userEmail,           // Subject (user identifier)
    userRole,            // Role ("admin", "superadmin", "user")
    cfg.JWTSecret,       // Secret key from config
    24 * time.Hour,      // Token expiration time
)
```

## Setting Authentication Cookie

For web applications, set the JWT as an HTTP-only cookie:

```go
c.Cookie(&fiber.Cookie{
    Name:     "admin_token",
    Value:    token,
    Path:     "/",
    HTTPOnly: true,
    Secure:   false,          // Set to true in production with HTTPS
    SameSite: fiber.CookieSameSiteLaxMode,
    Expires:  time.Now().Add(24 * time.Hour),
})
```

## Logout Implementation

To logout, clear the authentication cookie:

```go
func (h *Handler) HandleLogout(c *fiber.Ctx) error {
    c.Cookie(&fiber.Cookie{
        Name:     "admin_token",
        Value:    "",
        Path:     "/",
        HTTPOnly: true,
        Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
    })
    
    return c.Redirect("/admin/login")
}
```

## API Authentication

For API clients (mobile apps, SPAs), use the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Configuration

Ensure your JWT secret is properly configured:

```go
// config/config.go
type Config struct {
    JWTSecret string
}

func LoadConfig() *Config {
    return &Config{
        JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
    }
}
```

**Important**: Use a strong, random secret in production and set it via environment variable:

```bash
export JWT_SECRET="your-strong-secret-key-here"
```

## Security Best Practices

1. **Use Strong Secrets**: Generate a random, long secret key for JWT signing
2. **HTTPS Only**: Set `Secure: true` for cookies in production
3. **Token Expiration**: Use reasonable token expiration times (e.g., 24 hours)
4. **HTTP-Only Cookies**: Prevent XSS attacks by using HTTP-only cookies
5. **Environment Variables**: Never hardcode secrets, use environment variables
6. **Token Refresh**: Implement token refresh mechanism for better UX
7. **Rate Limiting**: Add rate limiting to login endpoints
8. **Password Hashing**: Always use bcrypt for password hashing (already implemented)

## Error Handling

The middleware returns appropriate HTTP status codes:

- **401 Unauthorized**: Missing or invalid token
- **403 Forbidden**: Valid token but insufficient permissions

Error responses are in JSON format:

```json
{
    "error": "Missing authentication token"
}
```

## Testing

To test authentication in development:

1. Login to get a token:
```bash
curl -X POST http://localhost:8080/admin/login \
  -d "email=admin@travel.app" \
  -d "password=admin123"
```

2. Use the token in subsequent requests:
```bash
curl http://localhost:8080/admin/dashboard \
  -H "Authorization: Bearer <your-token>"
```

Or with cookie:
```bash
curl http://localhost:8080/admin/dashboard \
  --cookie "admin_token=<your-token>"
```
