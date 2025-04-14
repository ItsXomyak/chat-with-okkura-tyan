# AuthForge

![Go](https://img.shields.io/badge/Go-1.23-blue) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue) ![Docker](https://img.shields.io/badge/Docker-✓-blue)

## 📌 Description
AuthForge is a RESTful API for easy and fast integration of authentication and user management. The project provides basic functions such as registration, authorization, login and password reset. For the demonstration, we use **Go**, **PostgreSQL**, and **Docker.**

## 🚀 Functionality
- User Registration and Login
- Password Reset (forgot password workflow)
- Logging of user activities and critical events

## 📂 Project structure
```bash
authforge/
├── cmd/               # Initalize app
├── internal/
│   ├── api/           # HTTP handlers
│   │   └── router/    # Endpoints
│   ├── config/        # Configurations and environment variables
│   ├── logger/        # Custom log directory
│   ├── mailer/        # Mail sending service
│   ├── models/        # Definitions of data structures
│   ├── repository/    # Working with the database
│   └── service/       # Business process logic
├── main.go            # Main application launch
├── migrations/        # SQL scripts for DATABASE migrations
├── docs/              # Swagger-documentation
├── Dockerfile         # Instructions for container assembly
├── docker-compose.yml # Docker Compose to launch the service
├── .env.example       # Environment variables template
└── README.md          # This file
```

## 🛠️ Installation and launch
### 🔹 1. Cloning a repository
```sh
git clone https://github.com/azhaxyly/authforge.git
cd authforge
```

### 🔹 2. Setting up the environment
Copy `.env.example` to `.env` and adjust the variables:
```sh
cp .env.example .env
```

Example `.env`:
```env
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=authforge
DB_PORT=5432
```

### 🔹 3. Launching in Docker
```sh
docker-compose up --build
```

Examples of API requests:
- `POST /api/v1/auth/register` — Register a new user
- `POST /api/v1/auth/login` — Authenticate and log in a user
- `POST /api/v1/auth/confirm` — Confirm a registered account
- `POST /api/v1/auth/password-reset-request` — Request a password reset
- `POST /api/v1/auth/password-reset-confirm` — Reset the password using a confirmation token

## 📦 Development
### 🔹 Local launch without Docker
1. Install Go and PostgreSQL.
2. Create a `authforge` database.
3. Configure the `.env`.
4. Start the server:
```sh
go run main.go
```

## 📜 License
MIT License © 2025
