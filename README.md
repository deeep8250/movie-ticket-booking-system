# Movie Ticket Booking System

A backend API for a movie ticket booking platform built with Go, Gin, PostgreSQL, Redis, Docker, and GitHub Actions.

The project supports user authentication, movie/show browsing, seat locking, seat booking, booking cancellation, and booking history.

---

## Features

- User signup and login
- JWT-based authentication
- Public movie, theater, show, and seat APIs
- Seat availability checking
- Redis-based temporary seat locking
- Seat unlocking
- PostgreSQL transaction-based booking flow
- Booking cancellation
- User booking history
- Booking details by booking ID
- Dockerized backend
- Integration tests for core flows
- GitHub Actions CI/CD pipeline

---

## Tech Stack

- Go
- Gin
- PostgreSQL
- Redis
- sqlx
- golang-migrate
- JWT
- Docker
- Docker Compose
- GitHub Actions
- GitHub Container Registry

---

## Architecture

This project follows a feature-based three-layer architecture.

Each main feature has its own:

```text
Handler → Service → Repository
```

### Responsibility

```text
Handler     → Handles HTTP request and response
Service     → Handles business logic
Repository  → Handles database/Redis operations
```

### Main Features

```text
auth
movies
theaters
shows
seats
bookings
```

---

## Project Structure

```text
movie-ticket-booking-system/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── db/
│   ├── dto/
│   ├── features/
│   │   ├── auth/
│   │   ├── movies/
│   │   ├── theaters/
│   │   ├── shows/
│   │   ├── seats/
│   │   └── bookings/
│   ├── middleware/
│   └── routes/
├── migrations/
├── tests/
│   └── integration/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## Core Booking Flow

```text
User logs in
↓
User selects show and seats
↓
System locks selected seats in Redis
↓
User creates booking
↓
System validates seat lock ownership
↓
System creates booking inside a PostgreSQL transaction
↓
System stores booked seats
↓
System removes Redis seat locks
```

---

## Seat Status Logic

Seats can have these statuses:

```text
available
locked
booked
disabled
```

Priority:

```text
disabled > booked > locked > available
```

---

## API Endpoints

### Health Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |

---

### Public Routes

| Method | Endpoint | Description |
|---|---|---|
| POST | `/public/signup` | Create user account |
| POST | `/public/login` | Login user |
| GET | `/public/movies` | Get all movies |
| GET | `/public/movies/:id` | Get movie by ID |
| GET | `/public/movies/:id/shows` | Get shows by movie |
| GET | `/public/theaters` | Get all theaters |
| GET | `/public/theaters/shows/:id` | Get shows by theater |
| GET | `/public/theaters/shows/:id/seats` | Get seats for a show |

---

### Private Routes

These routes require JWT authentication.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/private/user` | Get logged-in user |
| POST | `/private/shows/:id/seats/lock` | Lock seats for a show |
| POST | `/private/shows/:id/seats/unlock` | Unlock seats for a show |
| POST | `/private/bookings` | Create booking |
| GET | `/private/bookings/:id/details` | Get booking details |
| GET | `/private/users/me/bookings` | Get user booking history |
| PATCH | `/private/bookings/:id/cancel` | Cancel booking |

---

## Environment Variables

Create a `.env` file:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=movie_booking

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your-secret-key
```

When running with Docker Compose, use service names:

```env
DB_HOST=db
REDIS_HOST=redis
```

---

## Run Locally

Start PostgreSQL and Redis first.

Run database migrations:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/movie_booking?sslmode=disable" up
```

Run the application:

```bash
go run ./cmd/api
```

Server starts on:

```text
http://localhost:8080
```

---

## Run with Docker

Build the Docker image:

```bash
docker build -t movie-ticket-booking-system:local .
```

Run using Docker Compose:

```bash
docker compose up --build
```

Check health:

```bash
curl http://localhost:8080/health
```

---

## Run Tests

Run all Go tests:

```bash
go test ./... -v
```

Run integration tests:

```powershell
.\run-integration-tests.ps1
```

Integration tests currently cover:

```text
Router setup
Health route
Protected route without token
Protected route with token
Signup and login flow
Seat lock flow
Seat unlock flow
Booking flow
Booking cancellation flow
```

---

## CI/CD

This project uses GitHub Actions for continuous integration.

The pipeline checks:

```text
Go dependency consistency
Code formatting
go vet
Go tests
Docker image build
Docker image push to GitHub Container Registry
```

---

## Current Status

Version 1 is completed.

Current version includes:

```text
Authentication
Movie/theater/show APIs
Seat locking with Redis
Booking with PostgreSQL transactions
Booking cancellation
Docker setup
CI/CD
Integration tests
```

---

## Future Improvements

- Pending booking flow
- Simulated payment confirmation
- Background worker to expire unpaid bookings
- Admin APIs for movies, theaters, halls, shows, and seats
- Role-based authorization
- Swagger/OpenAPI documentation
- More edge-case integration tests
- Production deployment

---

## Author

Developed by [deeep8250](https://github.com/deeep8250)
