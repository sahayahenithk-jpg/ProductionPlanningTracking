# Project Planning Tracking

This repository contains a simple authentication project with:
- Go backend using Gin, PostgreSQL, GORM, and JWT
- Vue 3 frontend using Vite and Axios
- PostgreSQL database named `ProjectPlanningTracking`

## Setup

### 1. Start PostgreSQL

From the project root:

```bash
docker-compose up -d
```

### 2. Backend

```bash
cd backend
go mod tidy
go run main.go
```

The backend listens on `http://localhost:8080`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend runs on `http://localhost:5173`.

## API Endpoints

- `POST /api/auth/register` - register a new user
- `POST /api/auth/login` - login and receive JWT token
- `GET /api/user` - fetch authenticated user profile

## Notes

- Use `.env.example` in `backend/` as a template for local environment variables.
- The backend is configured to allow CORS requests from `http://localhost:5173`.
