package controllers

import (
	"log"
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
	SendVerif(c *fiber.Ctx) error
	Verif(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authC struct {
	authS services.AuthS
	cache cache.Cache
	redis cache.RedisC
}

func NewAuthC(authS services.AuthS, cache cache.Cache, redis cache.RedisC) AuthC {
	return &authC{authS: authS, cache: cache, redis: redis}
}

func (a *authC) Login(c *fiber.Ctx) error {
	var user models.Login

	err := c.BodyParser(&user)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if errors := middleware.StructValidator(user); errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	token, usr, err := a.authS.VerifyCredential(c.Context(), user)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    usr,
		"token":   token,
		"status":  true,
		"message": "Login success!",
	})
}

func (a *authC) Register(c *fiber.Ctx) error {
	var user models.Register

	err := c.BodyParser(&user)
	if err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if err = middleware.InputChecker(user.Name, user.Email, user.Password); err != nil {
		return helper.Response(c, fiber.StatusNotAcceptable, nil, err.Error(), false)
	}

	if errors := middleware.StructValidator(user); errors != nil {
		return helper.Response(c, fiber.StatusConflict, nil, errors, false)
	}

	if err = a.authS.CreateUser(c.Context(), user); err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	return helper.Response(c, fiber.StatusOK, nil, "Register success!", true)
}

func (a *authC) Logout(c *fiber.Ctx) error {
	return helper.Response(c, fiber.StatusOK, nil, "Logout success!", true)
}

func (a *authC) SendVerif(c *fiber.Ctx) error {
	var user models.OTP

	err := c.BodyParser(&user)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	otp, err := helper.GenerateOTP(6)
	if err != nil {
		log.Println(err)
	}

	helper.SendOTP(otp, user.Email)
	a.redis.Set(user.Email, &otp)

	return helper.Response(c, fiber.StatusOK, nil, "Send OTP successful!", true)
}

func (a *authC) Verif(c *fiber.Ctx) error {
	var otp models.OTP

	err := c.BodyParser(&otp)
	if err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	code := a.redis.Get(otp.Email)

	if code != otp.Otp {
		return helper.Response(c, fiber.StatusBadRequest, nil, "Invalid OTP code!", false)
	}

	if err = a.authS.UpdateActive(c.Context(), otp.Email); err != nil {
		return helper.Response(c, fiber.StatusConflict, nil, err.Error(), false)
	}

	a.cache.Del("users")

	return helper.Response(c, fiber.StatusOK, nil, "Verify OTP code success!", true)
}
