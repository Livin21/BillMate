package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/livin21/billmate/internal/store"
)

type application struct {
	config       config
	store        store.Storage
	writeJson    func(w http.ResponseWriter, status int, v interface{})
	serverError  func(w http.ResponseWriter, err error)
	unAuthorized func(w http.ResponseWriter)
}

type config struct {
	addr      string
	db        dbConfig
	jwtSecret string
}

type dbConfig struct {
	connString   string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", app.listUsersHandler)
				r.Post("/", app.createUserHandler)
			})
			r.Route("/expenses", func(r chi.Router) {
				r.Get("/", app.listExpensesHandler)
				r.Post("/", app.createExpenseHandler)
			})
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.loginHandler)
		r.Post("/signup", app.signupHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}
