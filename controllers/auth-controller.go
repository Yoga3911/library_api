package controllers

import (
	"project_restapi/cache"
	"project_restapi/helper"
	"project_restapi/middleware"
	"project_restapi/models"
	"project_restapi/services"

	"github.com/gofiber/fiber/v2"
)

type AuthC interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authC struct {
	authS services.AuthS
	cache cache.Cache
}

func NewAuthC(authS services.AuthS, cache cache.Cache) AuthC {
	return &authC{authS: authS, cache: cache}
}

func (a *authC) Login(c *fiber.Ctx) error {
	var user models.Login

	err := c.BodyParser(&user)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	errors := middleware.StructValidator(user)
	if errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	usr, err := a.authS.VerifyCredential(c.Context(), user)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, usr, "Login success!", true)
}

func (a *authC) Register(c *fiber.Ctx) error {
	var user models.Register

	err := c.BodyParser(&user)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	err = middleware.InputChecker(user.Name, user.Email, user.Password)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	errors := middleware.StructValidator(user)
	if errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	token, hash, err := a.authS.CreateUser(c.Context(), user)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	user.Password = hash
	user.RoleID = 1
	user.Coin = 2
	user.Token = token

	a.cache.Del("users")

	return helper.Response(c, fiber.StatusOK, user, "Register success!", true)
}

func (a *authC) Logout(c *fiber.Ctx) error {	
	return helper.Response(c, fiber.StatusOK, nil, "Logout success!", true)
}