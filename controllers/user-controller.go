package controllers

import (
	"log"
	"project_restapi/cache"
	"project_restapi/helper"
	"project_restapi/middleware"
	"project_restapi/models"
	"project_restapi/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserC interface {
	GetAll(c *fiber.Ctx) error
	GetByToken(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	ChangeEmail(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	TakeBook(c *fiber.Ctx) error
	GetBookById(c *fiber.Ctx) error
	DeleteBookById(c *fiber.Ctx) error
	UserBookById(c *fiber.Ctx) error
	RequestAdmin(c *fiber.Ctx) error
	RequestAnswer(c *fiber.Ctx) error
	GetRequest(c *fiber.Ctx) error
}

type userC struct {
	userS services.UserS
	cache cache.Cache
}

func NewUserC(userS services.UserS, cache cache.Cache) UserC {
	return &userC{userS: userS, cache: cache}
}

func (u *userC) GetAll(c *fiber.Ctx) error {
	err := services.CheckRole(c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	if val := u.cache.Get("users"); val != nil {
		return helper.Response(c, fiber.StatusOK, val, "Get all user success!", true)
	}

	users, err := u.userS.GetAll(c.Context(), c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	log.Println("User cache")
	u.cache.Set("users", users)

	return helper.Response(c, fiber.StatusOK, users, "Get all user success!", true)
}

func (u *userC) GetByToken(c *fiber.Ctx) error {
	user, err := u.userS.GetOne(c.Context(), c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, user, "Get user success!", true)
}

func (u *userC) UpdateUser(c *fiber.Ctx) error {
	var update models.Update

	err := c.BodyParser(&update)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if err = u.userS.CheckEmail(c.Context(), update.Email, c.Get("Authorization")); err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if err = middleware.InputChecker(update.Name, update.Email); err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if errors := middleware.StructValidator(update); errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	if update.B64Name != "-" {
		update.Image = strconv.FormatInt(time.Now().UnixMilli(), 10) + update.Email + "." + update.B64Name[11:14]
	} else {
		update.Image = "-"
	}
	
	token, err := u.userS.Update(c.Context(), update, c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	u.cache.Del("users")
	update.B64Name = "-"
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    update,
		"token":   token,
		"status":  true,
		"message": "Update user success",
	})
}

func (u *userC) ChangePassword(c *fiber.Ctx) error {
	var password models.ChangePass

	err := c.BodyParser(&password)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	token, err := u.userS.ChangePassword(c.Context(), password, c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	u.cache.Del("users")

	return helper.Response(c, fiber.StatusOK, token, "Change user password success!", true)
}

func (u *userC) ChangeEmail(c *fiber.Ctx) error {
	 c.Get("Authorization")
	var password models.ChangePass

	err := c.BodyParser(&password)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	token, err := u.userS.ChangePassword(c.Context(), password, c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	u.cache.Del("users")

	return helper.Response(c, fiber.StatusOK, token, "Change user password success!", true)
}

func (u *userC) DeleteUser(c *fiber.Ctx) error {
	if err := u.userS.Delete(c.Context(), c.Get("Authorization")); err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	u.cache.Del("users")

	return helper.Response(c, fiber.StatusOK, nil, "Delete user success!", true)
}

func (u *userC) TakeBook(c *fiber.Ctx) error {
	if err := u.userS.TakeBook(c.Context(), c.Params("id"), c.Get("Authorization")); err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, nil, "Take book by user id success!", true)
}

func (u *userC) GetBookById(c *fiber.Ctx) error {
	book, err := u.userS.GetById(c.Context(), c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, book, "Get book by user id success!", true)
}

func (u *userC) DeleteBookById(c *fiber.Ctx) error {
	if err := u.userS.DeleteBookById(c.Context(), c.Params("id"), c.Get("Authorization")); err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, nil, "Delete book by user id success!", true)
}

func (u *userC) UserBookById(c *fiber.Ctx) error {
	book, err := u.userS.UserOneBook(c.Context(), c.Get("Authorization"), c.Params("id"))
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, book, "Get book by user id success!", true)
}

func (u *userC) RequestAdmin(c *fiber.Ctx) error {
	if err := u.userS.ReqAdmin(c.Context(), c.Get("Authorization")); err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, nil, "Send your request is success!", true)
}

func (u *userC) RequestAnswer(c *fiber.Ctx) error {
	var ans models.Request

	err := c.BodyParser(&ans)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	token, err := u.userS.AnsAdmin(c.Context(), ans, c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, token, "Review request is success!", true)
}

func (u *userC) GetRequest(c *fiber.Ctx) error {
	users, err := u.userS.GetAllRequest(c.Context(), c.Get("Authorization"))
	if err != nil {
		return helper.Response(c, fiber.StatusBadRequest, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, users, "Get all request success!", true)
}
