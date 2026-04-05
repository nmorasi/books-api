package main

import (
	"books-api/internal/annotation"
	"books-api/internal/auth"
	"books-api/internal/book"
	"books-api/internal/character"
	"books-api/internal/db"
	"books-api/internal/middleware"
	"books-api/internal/summary"
	"books-api/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func runMigrations(database *sqlx.DB) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id            BIGSERIAL PRIMARY KEY,
			name          TEXT        NOT NULL,
			email         TEXT        NOT NULL UNIQUE,
			password_hash TEXT        NOT NULL,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS books (
			id          BIGSERIAL PRIMARY KEY,
			user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title       TEXT        NOT NULL,
			author      TEXT        NOT NULL,
			description TEXT        NOT NULL DEFAULT '',
			year        INTEGER     NOT NULL DEFAULT 0,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS annotations (
			id         BIGSERIAL PRIMARY KEY,
			book_id    BIGINT      NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			body       TEXT        NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS characters (
			id         BIGSERIAL PRIMARY KEY,
			book_id    BIGINT      NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name       TEXT        NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, m := range migrations {
		if _, err := database.Exec(m); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}
	log.Println("migrations applied")
}

func main() {
	database := db.Connect()
	runMigrations(database)

	userRepo := user.NewRepository(database)
	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	bookRepo := book.NewRepository(database)
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

	annotationRepo := annotation.NewRepository(database)
	annotationService := annotation.NewService(annotationRepo)
	annotationHandler := annotation.NewHandler(annotationService)
	summaryHandler := summary.NewHandler(annotationRepo)
	characterRepo := character.NewRepository(database)
	characterHandler := character.NewHandler(characterRepo)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes
	authHandler.RegisterRoutes(r)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		bookHandler.RegisterRoutes(r)
		annotationHandler.RegisterRoutes(r)
		summaryHandler.RegisterRoutes(r)
		characterHandler.RegisterRoutes(r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
