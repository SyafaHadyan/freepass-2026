package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Admin(ctx *fiber.Ctx) error {
	authToken := ctx.GetReqHeaders()["Authorization"]

	if len(authToken) < 1 {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unaouthorized",
		)
	}

	if ctx.Locals("role").(string) != "ADMIN" {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unaouthorized",
		)
	}

	bearerToken := authToken[0]
	token := strings.Split(bearerToken, " ")

	userID, err := m.jwt.ValidateToken(token[1])
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"token invalid",
		)
	}

	ctx.Locals("userID", userID.String())

	return ctx.Next()
}
