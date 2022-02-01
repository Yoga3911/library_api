package controllers

import (
	"log"
	"project_restapi/cache"
	"project_restapi/helper"
	"project_restapi/middleware"
	"project_restapi/models"
	"project_restapi/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookC interface {
	GetAllBook(c *fiber.Ctx) error
	GetBookByGenre(c *fiber.Ctx) error
	AddBook(c *fiber.Ctx) error
	UpdateBook(c *fiber.Ctx) error
	DeleteBook(c *fiber.Ctx) error
	GetReview(c *fiber.Ctx) error
	AddReview(c *fiber.Ctx) error
	UpdateReview(c *fiber.Ctx) error
}

type bookC struct {
	bookS services.BookS
	cache cache.Cache
}

func NewBookC(bookS services.BookS, cache cache.Cache) BookC {
	return &bookC{bookS: bookS, cache: cache}
}

func (b *bookC) GetAllBook(c *fiber.Ctx) error {
	if val := b.cache.Get("books"); val != nil {
		return helper.Response(c, fiber.StatusOK, val, "Get all book success!", true)
	}

	book, err := b.bookS.GetAll(c.Context())
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	log.Println("Book cache")
	b.cache.Set("books", book)

	return helper.Response(c, fiber.StatusOK, book, "Get all book success!", true)
}

func (b *bookC) GetBookByGenre(c *fiber.Ctx) error {
	if val := b.cache.Get("books/" + c.Params("id")); val != nil {
		return helper.Response(c, fiber.StatusOK, val, "Get book by genre success!", true)
	}

	book, err := b.bookS.GetByGenre(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	log.Println("Book cache")
	b.cache.Set("books/"+c.Params("id"), book)

	return helper.Response(c, fiber.StatusOK, book, "Get book by genre success!", true)
}

func (b *bookC) AddBook(c *fiber.Ctx) error {
	err := services.CheckRole(c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	var book models.CUBook

	err = c.BodyParser(&book)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	err = middleware.InputChecker(book.Title, book.Author)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	errors := middleware.StructValidator(book)
	if errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	err = b.bookS.AddBook(c.Context(), book)
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	b.cache.Del("books", "books/"+strconv.FormatInt(int64(book.GenreID), 10))

	return helper.Response(c, fiber.StatusOK, nil, "Add book success!", true)
}

func (b *bookC) UpdateBook(c *fiber.Ctx) error {
	err := services.CheckRole(c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	var book models.CUBook

	err = c.BodyParser(&book)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	err = middleware.InputChecker(book.Title, book.Author)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	errors := middleware.StructValidator(book)
	if errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	err = b.bookS.UpdateBook(c.Context(), book, c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	b.cache.Del("books", "books/"+strconv.FormatInt(int64(book.GenreID), 10))

	return helper.Response(c, fiber.StatusOK, nil, "Update book success!", true)
}

func (b *bookC) DeleteBook(c *fiber.Ctx) error {
	err := services.CheckRole(c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	genreID, err := b.bookS.DeleteBook(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	b.cache.Del("books", "books/"+strconv.FormatInt(int64(genreID), 10), "review/" + c.Params("id"))

	return helper.Response(c, fiber.StatusOK, nil, "Delete book success!", true)
}

func (b *bookC) GetReview(c *fiber.Ctx) error {
	if val := b.cache.Get("review/" + c.Params("id")); val != nil {
		return helper.Response(c, fiber.StatusOK, val, "Get reviews success!", true)
	}

	review, err := b.bookS.GetReview(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	log.Println("Review cache")
	b.cache.Set("review/"+c.Params("id"), review)

	return helper.Response(c, fiber.StatusOK, review, "Get reviews success!", true)
}

func (b *bookC) AddReview(c *fiber.Ctx) error {
	var review models.AddReview

	err := c.BodyParser(&review)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	errors := middleware.StructValidator(review)
	if errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	err = b.bookS.AddReview(c.Context(), review, c.Get("Authorization"), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}
	
	b.cache.Del("books", "review/" + c.Params("id"))
	
	return helper.Response(c, fiber.StatusOK, nil, "Add reviews success!", true)
}

func (b *bookC) UpdateReview(c *fiber.Ctx) error {
	var review models.AddReview

	err := c.BodyParser(&review)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}
	
	err = b.bookS.UpdateReview(c.Context(), c.Get("Authorization"), c.Params("id"), review)
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	b.cache.Del("books", "review/" + c.Params("id"))

	return helper.Response(c, fiber.StatusOK, nil, "Update reviews success!", true)
}
