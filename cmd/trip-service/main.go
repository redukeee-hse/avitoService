package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redukeee-hse/avitoService/internal/config"
	"github.com/redukeee-hse/avitoService/internal/database"
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
		if err != nil {
			json.NewEncoder(w).Encode(http.StatusServiceUnavailable)
		} else {
			json.NewEncoder(w).Encode(http.StatusOK)
		}
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(http.StatusOK)
	})

	router.Post("/api/v1/trips", func(w http.ResponseWriter, r *http.Request) {

	})

	router.Get("/api/v1/trips/{tripId}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		validId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			json.NewEncoder(w).Encode(http.StatusBadRequest)
			return
		}
		var name string

		err = pool.QueryRow(
			ctx,
			`SELECT name from trips WHERE id = $1`,
			validId,
		).Scan(&name)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Trip exists")
	})

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	server.ListenAndServe()

}
