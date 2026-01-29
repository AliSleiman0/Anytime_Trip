# Anytime Travel Backend

Go-based backend server for the Anytime Travel application using MongoDB.

## Prerequisites

- Go 1.21 or higher
- MongoDB 4.4 or higher
- Port 8080 available

## Quick Start

### 1. Install MongoDB

**Ubuntu/Debian:**
```bash
sudo apt-get install mongodb
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

**MacOS:**
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

**Check MongoDB is running:**
```bash
mongosh --eval "db.version()"
```

### 2. Setup Environment

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` and update the values as needed:
```bash
nano .env
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Build the Server

```bash
go build -o bin/server main.go
```

### 5. Run the Server

```bash
./bin/server
```

Or run directly without building:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Project Structure

```
backend/
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go
├── core/
│   └── utils/             # Utility functions (password hashing, etc.)
├── internal/
│   ├── database/          # Database connection
│   ├── handlers/          # HTTP request handlers
│   │   ├── admin/        # Admin handlers
│   │   ├── app/          # App user handlers (signup, login)
│   │   └── superadmin/   # Super admin handlers
│   ├── middleware/        # Authentication & authorization middleware
│   ├── models/           # Data models
│   │   ├── admin/
│   │   ├── app/          # User, SignupRequest, LoginRequest
│   │   └── superadmin/
│   ├── repository/       # Database operations
│   │   ├── admin/
│   │   ├── app/          # User repository
│   │   └── superadmin/
│   └── routes/           # Route definitions
│       ├── admin/
│       ├── app/          # /api/app/* routes
│       └── superadmin/
├── static/               # Static assets
├── templates/            # HTML templates
└── bin/                  # Compiled binaries
```

## API Endpoints

### App User Endpoints

- `POST /api/app/signup` - Register new user
- `POST /api/app/login` - User authentication
- `GET /api/app/dashboard` - User dashboard (protected)
- `GET /api/app/profile` - User profile (protected)

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for detailed API documentation.

## Database

### MongoDB Collections

- `users` - App user accounts
- `admins` - Admin accounts
- `car_bookings` - Car rental bookings
- `flight_bookings` - Flight bookings
- `hotel_bookings` - Hotel bookings
- `support_tickets` - Customer support tickets
- `payments` - Payment transactions

### Create Database Indexes

```bash
mongosh anytime_travel
```

```javascript
// User indexes
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "phone_number": 1 }, { unique: true })
db.users.createIndex({ "created_at": -1 })
db.users.createIndex({ "last_login": -1 })
```

## Testing

### Test Signup

```bash
curl -X POST http://localhost:8080/api/app/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "phone_number": "+1234567890",
    "password": "password123",
    "confirm_password": "password123",
    "sex": "Male",
    "country": "United States"
  }'
```

### Test Login

```bash
curl -X POST http://localhost:8080/api/app/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Development

### Run with Auto-Reload

Install air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml`:
```bash
air init
```

Run with auto-reload:
```bash
air
```

### Run Tests

```bash
go test ./...
```

## Production Deployment

### 1. Set Production Environment Variables

```bash
export ENVIRONMENT=production
export JWT_SECRET=$(openssl rand -hex 32)
export MONGO_URI=mongodb://your-production-mongodb-uri
```

### 2. Build for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server main.go
```

### 3. Run with systemd

Create `/etc/systemd/system/anytime-travel.service`:

```ini
[Unit]
Description=Anytime Travel API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/anytime-travel/backend
ExecStart=/opt/anytime-travel/backend/bin/server
Restart=always
Environment="PORT=8080"
Environment="MONGO_URI=mongodb://localhost:27017"
Environment="MONGO_DB_NAME=anytime_travel"
Environment="JWT_SECRET=your-secret-here"
Environment="ENVIRONMENT=production"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable anytime-travel
sudo systemctl start anytime-travel
sudo systemctl status anytime-travel
```

## Troubleshooting

### MongoDB Connection Failed

Check if MongoDB is running:
```bash
sudo systemctl status mongodb
```

### Port Already in Use

Change the port in `.env`:
```bash
PORT=8081
```

### View Logs

```bash
# Systemd service logs
sudo journalctl -u anytime-travel -f

# Application logs (if running manually)
./bin/server 2>&1 | tee server.log
```

## License

Proprietary - Anytime Travel 2026
