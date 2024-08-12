package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"http-server/store"
)

type bookRoutes struct {
	store store.BookStore
}

func (routes *bookRoutes) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := routes.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			jsonResponse(w, http.StatusNotFound, nil)
			return
		}

		log.Println(err)
		jsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	jsonResponse(w, http.StatusOK, book)
}

func (routes *bookRoutes) Find(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "Find")
}

func (routes *bookRoutes) Post(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Author string `json:"author"`
		Title  string `json:"title"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	log.Printf("Body %+v", body)

	book, err := routes.store.Create(r.Context(), body.Title, body.Author)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create book: %s", err))
		return
	}

	jsonResponse(w, http.StatusOK, book)
}

func (routes *bookRoutes) Put(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "Put")
}

func (routes *bookRoutes) Patch(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "Patch")
}

func (routes *bookRoutes) Delete(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "Delete")
}
