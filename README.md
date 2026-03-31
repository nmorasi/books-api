# Books API — Monorepo

A Go REST API with JWT authentication, PostgreSQL, and Railway deployment.

```
books-api/
├── api/          ← Go backend
└── frontend/     ← React frontend (coming soon)
```

---

## Stack

- **Language:** Go
- **Router:** Chi v5
- **Database:** PostgreSQL (sqlx + lib/pq)
- **Auth:** JWT (golang-jwt/jwt v5) + bcrypt
- **Deployment:** Railway (Nixpacks build)

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | — | Healthcheck |
| POST | `/auth/register` | — | Register user |
| POST | `/auth/login` | — | Login, returns JWT |
| GET | `/books` | Bearer token | List your books |
| POST | `/books` | Bearer token | Create a book |
| GET | `/books/{id}` | Bearer token | Get a book |
| PUT | `/books/{id}` | Bearer token | Update a book |
| DELETE | `/books/{id}` | Bearer token | Delete a book |

---

## Local Development

### Prerequisites
- Go 1.21+
- Docker

### 1. Start PostgreSQL
```bash
cd api
docker compose up -d
```
Postgres runs on port `5434` (5432 and 5433 used by other projects locally).

### 2. Run the server
```bash
cd api
PORT=8082 go run ./cmd/api/main.go
```

### 3. Test with curl
```bash
# Register
curl -s -X POST http://localhost:8082/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Your Name","email":"you@example.com","password":"secret123"}'

# Login and grab token
TOKEN=$(curl -s -X POST http://localhost:8082/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"secret123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Create a book
curl -s -X POST http://localhost:8082/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Clean Code","author":"Robert C. Martin","year":2008}'

# List books
curl -s http://localhost:8082/books -H "Authorization: Bearer $TOKEN"
```

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | Full Postgres URL (used on Railway) | — |
| `DB_HOST` | Postgres host (local fallback) | `localhost` |
| `DB_PORT` | Postgres port (local fallback) | `5434` |
| `DB_USER` | Postgres user (local fallback) | `postgres` |
| `DB_PASSWORD` | Postgres password (local fallback) | `postgres` |
| `DB_NAME` | Postgres database name (local fallback) | `booksapi` |
| `JWT_SECRET` | Secret for signing JWTs | `dev-secret-change-in-production` |

---

## Railway Deployment

### Prerequisites — install CLIs
```bash
brew install railway gh
```

### 1. Authenticate
```bash
gh auth login       # opens browser
railway login       # opens browser
```

### 2. Push to GitHub
```bash
git init
git add .
git commit -m "initial commit"
gh repo create books-api --public --source=. --push
```

### 3. Create Railway project
```bash
cd api
railway init        # creates new project, select workspace
```

### 4. Deploy
```bash
railway up --detach
```

### 5. Add PostgreSQL
The `railway add --database postgres` CLI command may require re-login.
If it fails, do it in the dashboard:
1. Open your project: `railway open`
2. Click **+ Add** → **Database** → **PostgreSQL**

### 6. Link DATABASE_URL to your service
In the dashboard → click your service → **Variables** tab →
click **"Trying to connect a database? Add Variable"** → select **DATABASE_URL** → **Add**

Or via CLI:
```bash
railway variable --set "DATABASE_URL=\${{Postgres.DATABASE_URL}}"
```

### 7. Set JWT_SECRET
```bash
railway variable --set "JWT_SECRET=$(openssl rand -hex 32)"
```

### 8. Redeploy with new variables
```bash
railway redeploy --yes
```

### 9. Generate a public URL
```bash
railway domain
```

### Useful Railway commands
```bash
railway logs          # stream logs
railway logs --build  # build logs only
railway variable      # list all env vars
railway status        # show linked project/service
railway open          # open dashboard in browser
```

---

## Database Migrations

Migrations run automatically at startup (in `main.go`). Safe to redeploy at any time — all migrations use `IF NOT EXISTS`.

To add a new migration, add a SQL statement to the `runMigrations()` function in [api/cmd/api/main.go](api/cmd/api/main.go):

```go
`ALTER TABLE books ADD COLUMN IF NOT EXISTS rating INT NOT NULL DEFAULT 0`,
```

---

## Adding a New Resource

Follow the same pattern as `book/`:

1. `api/internal/<resource>/repository.go` — DB queries
2. `api/internal/<resource>/service.go` — business logic
3. `api/internal/<resource>/handler.go` — HTTP handlers
4. Wire in `main.go`:
   ```go
   repo := resource.NewRepository(database)
   svc  := resource.NewService(repo)
   resource.NewHandler(svc).RegisterRoutes(r)
   ```
5. Add migration to `runMigrations()` in `main.go`
6. `git push` — Railway auto-deploys on push if connected to GitHub

---

## Project Links

- **Live API:** https://thorough-clarity-production-9138.up.railway.app
- **GitHub:** https://github.com/nmorasi/books-api
- **Railway Dashboard:** https://railway.com/project/3349c895-ed2d-4ae9-b251-71b6b2e6f49a
