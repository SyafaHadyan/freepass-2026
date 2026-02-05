package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Admin(ctx *fiber.Ctx) error {
	if ctx.Locals("role").(string) != "ADMIN" {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	return ctx.Next()
}

func (m *Middleware) Canteen(ctx *fiber.Ctx) error {
	if ctx.Locals("role").(string) != "CANTEEN" {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	return ctx.Next()
}

func (m *Middleware) AdminOrCanteen(ctx *fiber.Ctx) error {
	role := ctx.Locals("role").(string)

	if role == "CANTEEN" || role == "ADMIN" {
	} else {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	return ctx.Next()
}
