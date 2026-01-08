package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Elvi-2202/book-api/internal/handler"
	"github.com/Elvi-2202/book-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASS", "password")
	dbName := getEnv("DB_NAME", "book_db")

	db, err := repository.ConnectDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresBookRepository(db)
	bookHandler := handler.NewBookHandler(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := db.PingContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status": "down", "database": "unreachable"}`))
			return
		}
		w.Write([]byte(`{"status": "up", "database": "connected"}`))
	})

	r.Route("/books", func(r chi.Router) {
		r.Post("/", bookHandler.CreateBook)
		r.Get("/", bookHandler.ListBooks)
		r.Get("/{id}", bookHandler.GetBook)
		r.Put("/{id}", bookHandler.UpdateBook)
		r.Delete("/{id}", bookHandler.DeleteBook)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Serveur sur http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Arrêt du serveur...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erreur lors de la fermeture: %v", err)
	}
	log.Println("Serveur arrêté proprement")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}