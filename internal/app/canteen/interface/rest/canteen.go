// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"

	"github.com/SyafaHadyan/freepass-2026/internal/app/canteen/usecase"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	"github.com/SyafaHadyan/freepass-2026/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CanteenHandler struct {
	Validator      *validator.Validate
	Middleware     middleware.MiddlewareItf
	CanteenUseCase usecase.CanteenUseCaseItf
	Config         *env.Env
}

func NewCanteenHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, canteenUseCase usecase.CanteenUseCaseItf,
	config *env.Env,
) {
	canteenHandler := CanteenHandler{
		Validator:      validator,
		Middleware:     middleware,
		CanteenUseCase: canteenUseCase,
		Config:         config,
	}

	routerGroup = routerGroup.Group("/canteen")

	routerGroup.Post("", middleware.Authentication, middleware.Canteen, canteenHandler.CreateCanteen)
	routerGroup.Post("/menu", middleware.Authentication, middleware.Canteen, canteenHandler.CreateMenu)
	routerGroup.Post("/menu/order", middleware.Authentication, canteenHandler.CreateOrder)
	routerGroup.Post("/payment", middleware.Authentication, canteenHandler.CreatePayment)
	routerGroup.Post("/menu/order/feedback", middleware.Authentication, canteenHandler.CreateFeedback)
	routerGroup.Patch("/menu/:id", middleware.Authentication, middleware.Canteen, canteenHandler.UpdateMenu)
	routerGroup.Patch("/menu/order", middleware.Authentication, middleware.Canteen, canteenHandler.UpdateOrder)
	routerGroup.Get("", middleware.Authentication, canteenHandler.GetCanteenList)
	routerGroup.Get("/:id", middleware.Authentication, canteenHandler.GetCanteenInfo)
	routerGroup.Get("/menu/:id", middleware.Authentication, canteenHandler.GetMenuInfo)
	routerGroup.Get("/menu/order", middleware.Authentication, canteenHandler.GetOrderInfo)
	routerGroup.Delete("/menu/:id", middleware.Authentication, middleware.Canteen, canteenHandler.SoftDeleteMenu)
	routerGroup.Delete("/menu/order/feedback/:id", middleware.Authentication, middleware.Canteen, canteenHandler.SoftDeleteFeedback)
}

func (c *CanteenHandler) CreateCanteen(ctx *fiber.Ctx) error {
	var createCanteen dto.CreateCanteen

	err := ctx.BodyParser(&createCanteen)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = c.Validator.Struct(createCanteen)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.CreateCanteen(createCanteen)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to create canteen")
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "canteen created",
		"payload": res,
	})
}

func (c *CanteenHandler) CreateMenu(ctx *fiber.Ctx) error {
	var createMenu dto.CreateMenu

	err := ctx.BodyParser(&createMenu)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = c.Validator.Struct(createMenu)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.CreateMenu(createMenu)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to create menu",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "menu created",
		"payload": res,
	})
}

func (c *CanteenHandler) CreateOrder(ctx *fiber.Ctx) error {
	var createOrder dto.CreateOrder

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createOrder.UserID = userID

	err = c.Validator.Struct(createOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.CreateOrder(createOrder)
	if err == gorm.ErrInvalidValue {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid quantity",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "order created",
		"payload": res,
	})
}

func (c *CanteenHandler) CreatePayment(ctx *fiber.Ctx) error {
	var createPayment dto.CreatePayment

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createPayment.UserID = userID

	err = c.Validator.Struct(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.CreatePayment(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to create payment",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "payment created",
		"payload": res,
	})
}

func (c *CanteenHandler) VerifyPayment(ctx *fiber.Ctx) error {
	var verifyPayment dto.VerifyPayment

	err := ctx.BodyParser(&verifyPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = c.CanteenUseCase.VerifyPayment(verifyPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to verify payment",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "payment verified",
	})
}

func (c *CanteenHandler) CreateFeedback(ctx *fiber.Ctx) error {
	var createFeedback dto.CreateFeedback

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createFeedback)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createFeedback.UserID = userID

	err = c.Validator.Struct(createFeedback)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.CreateFeedback(createFeedback)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to create feedback",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "feedback created",
		"payload": res,
	})
}

func (c *CanteenHandler) UpdateMenu(ctx *fiber.Ctx) error {
	var updateMenu dto.UpdateMenu

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	menuID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid menu id",
		)
	}

	err = ctx.BodyParser(&updateMenu)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	updateMenu.ID = menuID
	updateMenu.UserID = userID

	err = c.Validator.Struct(updateMenu)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.UpdateMenu(updateMenu)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"menu not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete menu",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "menu updated",
		"payload": res,
	})
}

func (c *CanteenHandler) UpdateOrder(ctx *fiber.Ctx) error {
	var updateOrder dto.UpdateOrder

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}

	err = ctx.BodyParser(&updateOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = c.Validator.Struct(updateOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.UpdateOrder(updateOrder, userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete order",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "order updated",
		"payload": res,
	})
}

func (c *CanteenHandler) GetCanteenList(ctx *fiber.Ctx) error {
	res, err := c.CanteenUseCase.GetCanteenList()
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get canteen list",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully get canteen list",
		"payload": res,
	})
}

func (c *CanteenHandler) GetCanteenInfo(ctx *fiber.Ctx) error {
	canteenID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid canteen id",
		)
	}

	res, err := c.CanteenUseCase.GetCanteenInfo(canteenID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get canteen info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully get canteen info",
		"payload": res,
	})
}

func (c *CanteenHandler) GetMenuInfo(ctx *fiber.Ctx) error {
	menuID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid menu id",
		)
	}

	res, err := c.CanteenUseCase.GetMenuInfo(menuID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get menu info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully get menu info",
		"payload": res,
	})
}

func (c *CanteenHandler) GetOrderInfo(ctx *fiber.Ctx) error {
	var getOrderInfo dto.GetOrderInfo

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}

	err = ctx.BodyParser(&getOrderInfo)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	getOrderInfo.UserID = userID

	err = c.Validator.Struct(getOrderInfo)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := c.CanteenUseCase.GetOrderInfo(getOrderInfo)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get order info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully retrieved order info",
		"payload": res,
	})
}

func (c *CanteenHandler) GetOrderList(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}

	res, err := c.CanteenUseCase.GetOrderList(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order list empty",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully retrieved order list",
		"payload": res,
	})
}

func (c *CanteenHandler) SoftDeleteMenu(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	menuID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid menu id",
		)
	}

	err = c.CanteenUseCase.SoftDeleteMenu(menuID, userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"menu not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete menu",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
}

func (c *CanteenHandler) SoftDeleteFeedback(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	feedbackID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid menu id",
		)
	}

	err = c.CanteenUseCase.SoftDeleteFeedback(feedbackID, userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"feedback not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete feedback",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
}
