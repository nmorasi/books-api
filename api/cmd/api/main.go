package main

import (
	"books-api/internal/auth"
	"books-api/internal/book"
	"books-api/internal/db"
	"books-api/internal/middleware"
	"books-api/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	database := db.Connect()

	userRepo := user.NewRepository(database)
	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	bookRepo := book.NewRepository(database)
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

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
