package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"project_restapi/models"
	"project_restapi/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type BookS interface {
	GetAll(ctx context.Context) ([]*models.Book, error)
	GetByGenre(ctx context.Context, genre_id string) ([]*models.Book, error)
	AddBook(ctx context.Context, book models.CUBook) error
	UpdateBook(ctx context.Context, book models.CUBook, id string) error
	DeleteBook(ctx context.Context, id string) (int, error)
	GetReview(ctx context.Context, book_id string) ([]*models.BookReview, error)
	AddReview(ctx context.Context, review models.AddReview, token string, book_id string) error
	UpdateReview(ctx context.Context, token string, book_id string, review models.AddReview) error
}

type bookS struct {
	bookR repository.BookR
	jwtS  JWTS
}

func NewBookS(bookR repository.BookR, jwtS JWTS) BookS {
	return &bookS{bookR: bookR, jwtS: jwtS}
}

func (b *bookS) GetAll(ctx context.Context) ([]*models.Book, error) {
	var book []*models.Book

	pg, err := b.bookR.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	defer pg.Close()

	for pg.Next() {
		var b models.Book
		err = pg.Scan(&b.ID, &b.Title, &b.Author, &b.Sinopsis, &b.Genre, &b.Quantity, &b.Rating)
		if err != nil {
			log.Println(err)
		}
		b.Rating = (math.Round(b.Rating*100) / 100)
		book = append(book, &b)
	}

	return book, nil
}

func (b *bookS) GetByGenre(ctx context.Context, genre_id string) ([]*models.Book, error) {
	var book []*models.Book

	pg, err := b.bookR.GetByGenre(ctx, genre_id)
	if err != nil {
		return nil, err
	}

	defer pg.Close()

	for pg.Next() {
		var b models.Book
		err = pg.Scan(&b.ID, &b.Title, &b.Author, &b.Sinopsis, &b.Genre, &b.Quantity, &b.Rating)
		if err != nil {
			log.Println(err)
		}
		b.Rating = (math.Round(b.Rating*100) / 100)
		book = append(book, &b)
	}

	return book, nil
}

func (b *bookS) AddBook(ctx context.Context, book models.CUBook) error {
	err := b.bookR.InsertBook(ctx, book)

	return err
}

func (b *bookS) UpdateBook(ctx context.Context, book models.CUBook, id string) error {
	err := b.bookR.UpdateBook(ctx, book, id)

	return err
}

func (b *bookS) DeleteBook(ctx context.Context, id string) (int, error) {
	var genreID int
	pg, err := b.bookR.DeleteBook(ctx, id)
	if err != nil {
		return 0, err
	}
	err = pg.Scan(&genreID)

	return genreID, err
}

func (b *bookS) GetReview(ctx context.Context, book_id string) ([]*models.BookReview, error) {
	var books []*models.BookReview

	pg, err := b.bookR.GetReview(ctx, book_id)
	if err != nil {
		return nil, err
	}

	defer pg.Close()

	for pg.Next() {
		var book models.BookReview
		var dateP time.Time
		err := pg.Scan(&book.ID, &book.UserID, &book.BookID, &book.Comment, &book.Rating, &dateP)
		book.PostDate = dateP.Format("02-01-2006")
		if err != nil {
			log.Println(err.Error())
		}
		books = append(books, &book)
	}

	return books, nil
}

func (b *bookS) AddReview(ctx context.Context, review models.AddReview, token string, book_id string) error {
	t, err := b.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	var count int
	b.bookR.CheckReview(ctx, claims["id"].(float64), book_id).Scan(&count)
	if count != 0 {
		return fmt.Errorf("you can only send 1 review")
	}

	err = b.bookR.AddReview(ctx, claims["id"].(float64), book_id, review)

	return err
}

func (b *bookS) UpdateReview(ctx context.Context, token string, book_id string, review models.AddReview) error {
	t, err := b.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	err = b.bookR.UpdateReview(ctx, claims["id"].(float64), book_id, review)

	return err
}
