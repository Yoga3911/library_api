package controllers

import (
	"fmt"
	"log"
	"project_restapi/cache"
	"project_restapi/helper"
	"project_restapi/middleware"
	"project_restapi/models"
	"project_restapi/services"
	"time"

	"github.com/gofiber/fiber/v2"
	cac "github.com/patrickmn/go-cache"
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
	ca := cac.New(5*time.Minute, 10*time.Minute)
	ca.LoadFile("ok")

	bk, ok := ca.Get("books")
	if ok {
		log.Println("1")
		return helper.Response(c, fiber.StatusOK, &bk, "Get all book success!", true)
	}
	// data := b.cache.GetCacheBook("books")
	// if data != nil {
	// 	log.Println("2")
	// 	return helper.Response(c, fiber.StatusOK, &data, "Get all book success!", true)
	// }

	book, err := b.bookS.GetAll(c.Context())
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}
	
	log.Println("Book cache")
	ca.SetDefault("books", &book)
	ca.SaveFile("ok")
	// err = b.cache.SetCache("books", &book, time.Minute*10)
	// if err != nil {
	// 	return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	// }

	return helper.Response(c, fiber.StatusOK, book, "Get all book success!", true)
}

func (b *bookC) GetBookByGenre(c *fiber.Ctx) error {
	data := b.cache.GetCacheBook("books/" + c.Params("id"))
	if data != nil {
		return helper.Response(c, fiber.StatusOK, data, "Get book by genre success!", true)
	}

	book, err := b.bookS.GetByGenre(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	log.Println("Book cache")
	err = b.cache.SetCache("books/"+c.Params("id"), book, time.Minute*10)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

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

	ca := cac.New(5*time.Minute, 10*time.Minute)
	ca.LoadFile("ok")
	ca.Delete("books")
	ca.SaveFile("ok")
	// b.cache.DestroyCache("books", "books/"+fmt.Sprintf("%v", book.GenreID))

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

	b.cache.DestroyCache("books", "books/"+fmt.Sprintf("%v", book.GenreID))

	return helper.Response(c, fiber.StatusOK, nil, "Update book success!", true)
}

func (b *bookC) DeleteBook(c *fiber.Ctx) error {
	err := services.CheckRole(c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	err = b.bookS.DeleteBook(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	// b.cache.DestroyCache("books", )

	return helper.Response(c, fiber.StatusOK, nil, "Delete book success!", true)
}

func (b *bookC) GetReview(c *fiber.Ctx) error {
	data := b.cache.GetCacheReview("review/" + c.Params("id"))
	if data != nil {
		return helper.Response(c, fiber.StatusOK, data, "Get reviews success!", true)
	}

	book, err := b.bookS.GetReview(c.Context(), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	log.Println("Book cache")
	err = b.cache.SetCache("review/"+c.Params("id"), book, time.Minute)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, book, "Get reviews success!", true)
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

	b.cache.DestroyCache("review/" + fmt.Sprintf("%v", c.Params("id")))

	return helper.Response(c, fiber.StatusOK, nil, "Add reviews success!", true)
}

func (b *bookC) UpdateReview(c *fiber.Ctx) error {
	var review models.AddReview

	err := c.BodyParser(&review)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	b.bookS.UpdateReview(c.Context(), c.Get("Authorization"), c.Params("id"), review)

	b.cache.DestroyCache("review/" + fmt.Sprintf("%v", c.Params("id")))

	return helper.Response(c, fiber.StatusOK, nil, "Update reviews success!", true)
}
