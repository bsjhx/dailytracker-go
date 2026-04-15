# Project Structure

This project follows the standard Go project layout from [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

```
dailytracker-go/
├── cmd/
│   └── dailytracker/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth.go             # Authentication & user management
│   │   ├── entries.go          # List & create entries
│   │   ├── entry.go            # Get & update single entry
│   │   └── stats.go            # Weekly statistics
│   ├── middleware/              # HTTP middleware
│   │   └── auth.go             # Session authentication
│   ├── models/                  # Data structures
│   │   ├── entry.go            # Entry models
│   │   ├── stats.go            # Stats models
│   │   └── user.go             # User & session models
│   └── repository/              # Database layer
│       └── database.go         # DB connection & migrations
├── web/
│   ├── static/                  # Static assets (CSS, JS, images)
│   └── templates/               # HTML templates
│       ├── index.html          # Main application page
│       └── login.html          # Login page
├── scripts/
│   └── create-user.sh          # Helper script for user creation
├── go.mod                       # Go module definition
└── go.sum                       # Go module checksums

## Building & Running

### Development
```bash
go run ./cmd/dailytracker
```

### Production
```bash
go build -o dailytracker ./cmd/dailytracker
./dailytracker
```

## Key Principles

- **`cmd/`**: Contains application entry points. Keep minimal - just wire things together.
- **`internal/`**: Private application code. Cannot be imported by external projects.
  - `handlers/`: HTTP handlers - handle requests, call business logic
  - `middleware/`: HTTP middleware - authentication, logging, etc.
  - `models/`: Data structures shared across packages
  - `repository/`: Database access layer
- **`web/`**: Web application assets
  - `static/`: Static files (future: CSS, JS)
  - `templates/`: HTML templates
- **`scripts/`**: Build, installation, and utility scripts

## Migration from Old Structure

The old flat structure:
```
main.go
api/
  ├── auth.go
  ├── db.go
  ├── entries.go
  ├── entry.go
  ├── middleware.go
  └── stats.go
public/
  ├── index.html
  └── login.html
```

Has been refactored to follow Go best practices with proper separation of concerns.
