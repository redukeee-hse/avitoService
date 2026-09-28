package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redukeee-hse/avitoService/internal/config"
	"github.com/redukeee-hse/avitoService/internal/database"
	api "github.com/redukeee-hse/avitoService/internal/generated"
)


func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка скачивания .env файлов:", err)
	}
	ctx := context.Background()
	router := chi.NewRouter()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Ошибка конфига:", err)
	}

	pool, err := database.NewPool(ctx, cfg.DB)

	if err != nil {
		log.Fatal("Ошибка создания пула:", err)
	}

	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, cancel := context.WithTimeout(r.Context(), cfg.PingTimeout)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, api.HealthResponse{Status: api.Unavailable})
			return
		}
		writeJSON(w, http.StatusOK, api.HealthResponse{Status: api.Ok})
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, api.HealthResponse{Status: api.Ok})
	})

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	serversErr := make(chan error, 1)
	go func() {
		log.Printf("Запускаю сервер на %s", cfg.Addr)
		serversErr <- server.ListenAndServe()
	}()

	err = <-serversErr
	log.Fatalf("Ошибка сервера: %v", err)

}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("ошибка записи ответа: %v", err)
	}
}
