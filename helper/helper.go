package helper

import "github.com/gofiber/fiber/v2"

func Response(c *fiber.Ctx, r int, d interface{}, m interface{}, s bool) error {
	return c.Status(r).JSON(fiber.Map{
		"data":    d,
		"message": m,
		"status":  s,
	})
}