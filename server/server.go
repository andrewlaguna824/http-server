package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"http-server/store"
)

type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	port   string
	client *pgxpool.Pool
}

func NewServer(
	port string,
	client *pgxpool.Pool,
) Server {
	return &server{
		port:   port,
		client: client,
	}
}

func (s *server) Run(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// public endpoints
	r.Group(func(r chi.Router) {
		r.Get("/", handleRoot)
	})

	// private endpoints
	r.Group(func(r chi.Router) {
		books := bookRoutes{store: store.NewBookStore(s.client)}

		r.Route("/books", func(r chi.Router) {
			r.Get("/", books.Find)
			r.Post("/", books.Post)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", books.Get)
				r.Put("/", books.Put)
				r.Patch("/", books.Patch)
				r.Delete("/", books.Delete)
			})
		})

	})

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: r,
	}

	serverErrChan := make(chan error)
	go func() {
		serverErrChan <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("api server context error: %w", ctx.Err())
	case err := <-serverErrChan:
		return fmt.Errorf("api server error: %w", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	type tmp struct {
		Message string `json:"message"`
	}

	jsonResponse(w, http.StatusOK, tmp{
		Message: "Hello, World",
	})
}

func jsonResponse(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to encode json response %s", err.Error())
	}
}
