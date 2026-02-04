// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"

	"github.com/SyafaHadyan/freepass-2026/internal/app/user/usecase"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/mailer"
	"github.com/SyafaHadyan/freepass-2026/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Validator   *validator.Validate
	Middleware  middleware.MiddlewareItf
	UserUseCase usecase.UserUseCaseItf
	Config      *env.Env
	Mailer      mailer.MailerItf
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
	config *env.Env, mailer mailer.MailerItf,
) {
	userHandler := UserHandler{
		Validator:   validator,
		Middleware:  middleware,
		UserUseCase: userUseCase,
		Config:      config,
		Mailer:      mailer,
	}

	routerGroup = routerGroup.Group("/users")

	routerGroup.Post("/register", userHandler.Register)
	routerGroup.Post("/login", userHandler.Login)
	routerGroup.Delete("/:username", middleware.Authentication)
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	var register dto.Register

	err := ctx.BodyParser(&register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := u.UserUseCase.Register(register)
	if err != nil {
		return fiber.NewError(
			http.StatusConflict,
			"please use another email / username",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user registered",
		"payload": res,
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	var login dto.Login

	err := ctx.BodyParser(&login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, token, err := u.UserUseCase.Login(login)
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"invalid username or password",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user authenticated",
		"token":   token,
		"payload": res,
	})
}

func (u *UserHandler) SoftDelete(ctx *fiber.Ctx) error {
	targetUserName := ctx.Params("username")
	userIDTarget, err := u.UserUseCase.GetUserIDFromUsername(targetUserName)
	if err != nil {
		return fiber.NewError(
			http.StatusNotFound,
			"target user not found")
	}

	u.UserUseCase.SoftDelete(userIDTarget)

	return nil
}
