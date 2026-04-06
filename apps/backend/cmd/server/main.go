package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/yourusername/attends-moi/internal/handler"
	"github.com/yourusername/attends-moi/internal/repository"
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/attends_moi?sslmode=disable")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	boardRepo := repository.NewBoardRepo(pool)
	cardRepo := repository.NewCardRepo(pool)
	commentRepo := repository.NewCommentRepo(pool)

	boardHandler := handler.NewBoardHandler(boardRepo)
	cardHandler := handler.NewCardHandler(cardRepo)
	commentHandler := handler.NewCommentHandler(commentRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Health checks
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/boards", func(r chi.Router) {
			r.Post("/", boardHandler.Create)
			r.Get("/", boardHandler.GetAll)
			r.Get("/{id}", boardHandler.GetByID)
			r.Get("/{boardID}/cards", cardHandler.GetByBoardID)
			r.Post("/{boardID}/cards", cardHandler.Create)
		})

		r.Route("/cards", func(r chi.Router) {
			r.Get("/{id}", cardHandler.GetByID)
			r.Patch("/{id}", cardHandler.Update)
			r.Delete("/{id}", cardHandler.Delete)
			r.Get("/{cardID}/comments", commentHandler.GetByCardID)
			r.Post("/{cardID}/comments", commentHandler.Create)
		})
	})

	addr := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: r,
	}

	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped gracefully")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
