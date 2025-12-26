# Go Project

A Go web server project with templates and frontend components.

## Structure

```
.
├── backend/                          # Go backend
│   ├── cmd/                          # Application entrypoints
│   │   └── main.go                   # Main application
│   ├── internal/                     # Private application code
│   │   ├── handlers/                 # HTTP handlers
│   │   │   ├── app/                  # App-level handlers
│   │   │   ├── admin/                # Admin-level handlers
│   │   │   └── superadmin/           # Super admin handlers
│   │   ├── models/                   # Data models
│   │   │   ├── app/                  # App models
│   │   │   ├── admin/                # Admin models
│   │   │   └── superadmin/           # Super admin models
│   │   ├── services/                 # Business logic
│   │   │   ├── app/                  # App services
│   │   │   ├── admin/                # Admin services
│   │   │   └── superadmin/           # Super admin services
│   │   ├── middleware/               # HTTP middleware
│   │   └── database/                 # Database layer
│   ├── pkg/                          # Public libraries
│   │   └── utils/                    # Utility functions
│   └── config/                       # Configuration
├── go.mod                            # Go module file
├── templates/                        # Go HTML templates
│   └── index.html
└── frontend/                         # Frontend files
    ├── index.html                    # Main HTML page
    ├── app.js                        # Main JavaScript file
    └── components/                   # JavaScript components
        ├── header.js
        └── footer.js
```

## Running the Project

```bash
cd backend/cmd
go run main.go
```

## API Endpoints

### App Routes
- `GET /api/app/dashboard` - App user dashboard
- `GET /api/app/profile` - User profile

### Admin Routes (requires authentication)
- `GET /api/admin/dashboard` - Admin dashboard
- `GET /api/admin/users` - User management

### Super Admin Routes (requires authentication)
- `GET /api/superadmin/dashboard` - Super admin dashboard
- `GET /api/superadmin/system` - System management

### Other Routes
- `GET /static/*` - Static frontend files
- `GET /template` - Template rendering example
```

The server will start on `http://localhost:8080`

- Main page (frontend): `http://localhost:8080`
- Template example: `http://localhost:8080/template`
