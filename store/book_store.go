package store

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	UUID   string `json:"uuid"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookUpdate struct {
	title  *string
	author *string
}

type BookStore interface {
	Create(ctx context.Context, title string, author string) (Book, error)
	Get(ctx context.Context, id string) (Book, error)
	Find(ctx context.Context, query string) ([]Book, error)
	Update(ctx context.Context, id string, update BookUpdate) (Book, error)
	Replace(ctx context.Context, id string, title string, author string) (Book, error)
	Delete(ctx context.Context, id string) error
}

type bookStore struct {
	client *pgxpool.Pool
}

func NewBookStore(client *pgxpool.Pool) BookStore {
	return &bookStore{client: client}
}

func (store bookStore) Create(ctx context.Context, title string, author string) (Book, error) {
	var b Book
	if err := store.client.QueryRow(
		ctx,
		"INSERT INTO book (title, author) VALUES ($1, $2) RETURNING (uuid, title, author)", title, author,
	).Scan(&b); err != nil {
		return Book{}, fmt.Errorf("failed to insert new book: %w", err)
	}

	return b, nil
}

func (store bookStore) Get(ctx context.Context, id string) (Book, error) {
	var b Book
	query := "SELECT uuid, author, title FROM book WHERE uuid=$1"

	if err := pgxscan.Get(ctx, store.client, &b, query, id); err != nil {
		return Book{}, fmt.Errorf("failed to get book with id %q: %w", id, err)
	}

	return b, nil
}

func (store bookStore) Find(ctx context.Context, query string) ([]Book, error) {
	panic("implement me")
}

func (store bookStore) Update(ctx context.Context, id string, update BookUpdate) (Book, error) {
	panic("implement me")
}

func (store bookStore) Replace(ctx context.Context, id string, title string, author string) (Book, error) {
	panic("implement me")
}

func (store bookStore) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
