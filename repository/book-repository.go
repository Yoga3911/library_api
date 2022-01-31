package repository

import (
	"context"
	"project_restapi/models"
	"project_restapi/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BookR interface {
	GetAll(ctx context.Context) (pgx.Rows, error)
	GetByGenre(ctx context.Context, genre_id string) (pgx.Rows, error)
	InsertBook(ctx context.Context, book models.CUBook) error
	UpdateBook(ctx context.Context, book models.CUBook, id string) error
	DeleteBook(ctx context.Context, id string) error
	GetReview(ctx context.Context, book_id string) (pgx.Rows, error)
	AddReview(ctx context.Context, id float64, book_id string, review models.AddReview) error
	CheckReview(ctx context.Context, id float64, book_id string) pgx.Row
	UpdateReview(ctx context.Context, id float64, book_id string, review models.AddReview) error
}

type bookR struct {
	db *pgxpool.Pool
}

func NewBookR(db *pgxpool.Pool) BookR {
	return &bookR{
		db: db,
	}
}

func (b *bookR) GetAll(ctx context.Context) (pgx.Rows, error) {
	pg, err := b.db.Query(ctx, sql.GetAllBook)

	return pg, err
}

func (b *bookR) GetByGenre(ctx context.Context, genre_id string) (pgx.Rows, error) {
	pg, err := b.db.Query(ctx, sql.GetByGenre, genre_id)

	return pg, err
}

func (b *bookR) InsertBook(ctx context.Context, book models.CUBook) error {
	_, err := b.db.Exec(ctx, sql.AddBook, book.Title, book.Author, book.Sinopsis, book.GenreID, book.Quantity)

	return err
}

func (b *bookR) UpdateBook(ctx context.Context, book models.CUBook, id string) error {
	_, err := b.db.Exec(ctx, sql.UpdateBook, id, book.Title, book.Author, book.Sinopsis, book.GenreID, book.Quantity)

	return err
}

func (b *bookR) DeleteBook(ctx context.Context, id string) error {
	_, err := b.db.Exec(ctx, sql.DeleteBook, id)

	return err
}

func (b *bookR) GetReview(ctx context.Context, book_id string) (pgx.Rows, error) {
	pg, err := b.db.Query(ctx, sql.GetReview, book_id)

	return pg, err
}

func (b *bookR) AddReview(ctx context.Context, id float64, book_id string, review models.AddReview) error {
	_, err := b.db.Exec(ctx, sql.AddReview, id, book_id, review.Comment, review.Rating)
	
	return err
}

func (b *bookR) CheckReview(ctx context.Context, id float64, book_id string) pgx.Row {
	pg := b.db.QueryRow(ctx, sql.CheckReview, id, book_id)

	return pg
}

func (b *bookR) UpdateReview(ctx context.Context, id float64, book_id string, review models.AddReview) error {
	_, err := b.db.Exec(ctx, sql.UpdateReview, id, book_id, review.Comment, review.Rating)

	return err
}