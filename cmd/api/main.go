package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/livin21/billmate/internal/db"
	"github.com/livin21/billmate/internal/env"
	"github.com/livin21/billmate/internal/store"
)

func main() {

	env.LoadEnv()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			connString:   env.GetString("DB_CONN_STRING", "postgres://postgres:postgres@localhost:5433/billmate?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		jwtSecret: env.GetString("JWT_SECRET", "secret"),
	}

	db, err := db.New(cfg.db.connString, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		writeJson: func(w http.ResponseWriter, status int, v interface{}) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			err := json.NewEncoder(w).Encode(v)
			if err != nil {
				panic(err)
			}
		},
		serverError: func(w http.ResponseWriter, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
		},
		unAuthorized: func(w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "unauthorized",
			})
		},
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

}
