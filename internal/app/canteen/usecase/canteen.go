// Package usecase handles the logic for each user request
package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/SyafaHadyan/freepass-2026/internal/app/canteen/repository"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/entity"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/payment"
	redisitf "github.com/SyafaHadyan/freepass-2026/internal/infra/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CanteenUseCaseItf interface {
	CreateCanteen(createCanteen dto.CreateCanteen) (dto.ResponseCreateCanteen, error)
	CreateMenu(createMenu dto.CreateMenu) (dto.ResponseCreateMenu, error)
	CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error)
	CreatePayment(createPayment dto.CreatePayment) (dto.ResponseMidtransOrder, error)
	VerifyPayment(verifyPayment dto.VerifyPayment) error
	CreateFeedback(createFeedback dto.CreateFeedback) (dto.ResponseCreateFeedback, error)
	UpdateMenu(updateMenu dto.UpdateMenu) (dto.ResponseUpdateMenu, error)
	UpdateOrder(updateOrder dto.UpdateOrder, userID uuid.UUID) (dto.ResponseUpdateOrder, error)
	GetCanteenList() ([]dto.ResponseGetCanteenList, error)
	GetCanteenInfo(canteenID uuid.UUID) (dto.ResponseGetCanteenInfo, error)
	GetMenuInfo(menuID uuid.UUID) (dto.ResponseGetMenuInfo, error)
	GetOrderInfo(getOrderInfo dto.GetOrderInfo) (dto.ResponseGetOrderInfo, error)
	GetOrderList(userID uuid.UUID) ([]dto.ResponseGetOrderList, error)
	SoftDeleteMenu(menuID uuid.UUID, userID uuid.UUID) error
	SoftDeleteFeedback(feedbackID uuid.UUID, userID uuid.UUID) error
}

type CanteenUseCase struct {
	canteenRepo  repository.CanteenDBItf
	Payment      payment.PaymentItf
	Env          *env.Env
	redis        redisitf.RedisItf
	redisContext context.Context
}

func NewCanteenUseCase(
	canteenRepo repository.CanteenDBItf, payment payment.PaymentItf,
	env *env.Env, redis redisitf.RedisItf,
) CanteenUseCaseItf {
	return &CanteenUseCase{
		canteenRepo:  canteenRepo,
		Payment:      payment,
		Env:          env,
		redis:        redis,
		redisContext: context.Background(),
	}
}

func (c *CanteenUseCase) CreateCanteen(createCanteen dto.CreateCanteen) (dto.ResponseCreateCanteen, error) {
	canteen := entity.Canteen{
		ID:     uuid.New(),
		UserID: createCanteen.UserID,
		Name:   createCanteen.Name,
	}

	err := c.canteenRepo.CreateCanteen(&canteen)

	return canteen.ParseToDTOResponseCreateCanteen(), err
}

func (c *CanteenUseCase) CreateMenu(createMenu dto.CreateMenu) (dto.ResponseCreateMenu, error) {
	menu := entity.Menu{
		ID:        uuid.New(),
		CanteenID: createMenu.CanteenID,
		Name:      createMenu.Name,
		Price:     createMenu.Price,
	}

	err := c.canteenRepo.CreateMenu(&menu)

	return menu.ParseToDTOResponseCreateMenu(), err
}

func (c *CanteenUseCase) CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error) {
	menu := entity.Menu{
		ID: createOrder.MenuID,
	}

	order := entity.Order{
		ID:       uuid.New(),
		UserID:   createOrder.UserID,
		MenuID:   createOrder.MenuID,
		Quantity: createOrder.Quantity,
		Status:   "UNPAID",
	}

	err := c.canteenRepo.CreateOrder(&menu, &order)

	return order.ParseToDTOResponseCreateOrder(), err
}

func (c *CanteenUseCase) CreatePayment(createPayment dto.CreatePayment) (dto.ResponseMidtransOrder, error) {
	paymentID := uuid.New()

	createMidtransOrder := dto.CreateMidtransOrder{
		TransactionDetails: dto.TransactionDetails{
			OrderID:     paymentID.String(),
			GrossAmount: createPayment.Price,
		},
	}

	res, err := c.Payment.CreatePayment(createMidtransOrder)
	if err != nil {
		return dto.ResponseMidtransOrder{}, err
	}

	responseMidtransOrder := dto.ResponseMidtransOrder{
		Token:       res.Token,
		RedirectURL: res.RedirectURL,
	}

	payment := entity.Payment{
		ID:          paymentID,
		OrderID:     createPayment.OrderID,
		UserID:      createPayment.UserID,
		Price:       createPayment.Price,
		RedirectURL: responseMidtransOrder.RedirectURL,
	}

	err = c.canteenRepo.CreatePayment(&payment)

	return responseMidtransOrder, err
}

func (c *CanteenUseCase) VerifyPayment(verifyPayment dto.VerifyPayment) error {
	orderID, _ := uuid.Parse(verifyPayment.TransactionID)

	signatureKey := fmt.Sprintf(
		"%s,%s,%s,%s,",
		verifyPayment.TransactionID,
		verifyPayment.StatusCode,
		verifyPayment.GrossAmount,
		c.Env.MidtransServerKey,
	)

	hash := sha512.New()
	hash.Write([]byte(signatureKey))
	hashedData := hash.Sum(nil)
	hexHash := hex.EncodeToString(hashedData)

	if hexHash != verifyPayment.SignatureKey {
		return fiber.NewError(
			http.StatusUnauthorized,
			"payment could not be verified")
	}

	order := entity.Order{
		ID: orderID,
	}

	err := c.canteenRepo.VerifyPayment(&order)

	return err
}

func (c *CanteenUseCase) CreateFeedback(createFeedback dto.CreateFeedback) (dto.ResponseCreateFeedback, error) {
	feedback := entity.Feedback{
		ID:      uuid.New(),
		OrderID: createFeedback.OrderID,
		UserID:  createFeedback.UserID,
	}

	err := c.canteenRepo.CreateFeedback(&feedback)

	return feedback.ParseToDTOResponseCreateFeedback(), err
}

func (c *CanteenUseCase) UpdateMenu(updateMenu dto.UpdateMenu) (dto.ResponseUpdateMenu, error) {
	menu := entity.Menu{
		ID:    updateMenu.ID,
		Name:  updateMenu.Name,
		Price: updateMenu.Price,
	}

	err := c.canteenRepo.UpdateMenu(&menu, updateMenu.UserID)

	return menu.ParseToDTOResponseUpdateMenu(), err
}

func (c *CanteenUseCase) UpdateOrder(updateOrder dto.UpdateOrder, userID uuid.UUID) (dto.ResponseUpdateOrder, error) {
	order := entity.Order{
		ID:     updateOrder.ID,
		MenuID: updateOrder.MenuID,
		Status: updateOrder.Status,
	}

	err := c.canteenRepo.UpdateOrder(&order, userID)

	return order.ParseToDTOResponseUpdateOrder(), err
}

func (c *CanteenUseCase) GetCanteenList() ([]dto.ResponseGetCanteenList, error) {
	canteen := new([]entity.Canteen)

	err := c.canteenRepo.GetCanteenList(canteen)
	if err != nil {
		return nil, err
	}

	parsedCanteen := make([]dto.ResponseGetCanteenList, len(*canteen))

	for i, c := range *canteen {
		parsedCanteen[i] = c.ParseToDTOResponseGetCanteenList()
	}

	return parsedCanteen, err
}

func (c *CanteenUseCase) GetCanteenInfo(canteenID uuid.UUID) (dto.ResponseGetCanteenInfo, error) {
	canteen := entity.Canteen{
		ID: canteenID,
	}

	err := c.canteenRepo.GetCanteenInfo(&canteen)

	return canteen.ParseToDTOResponseGetCanteenInfo(), err
}

func (c *CanteenUseCase) GetMenuInfo(menuID uuid.UUID) (dto.ResponseGetMenuInfo, error) {
	menu := entity.Menu{
		ID: menuID,
	}

	err := c.canteenRepo.GetMenuInfo(&menu)

	return menu.ParseToDTOResponseGetMenuInfo(), err
}

func (c *CanteenUseCase) GetOrderInfo(getOrderInfo dto.GetOrderInfo) (dto.ResponseGetOrderInfo, error) {
	order := entity.Order{
		ID:     getOrderInfo.ID,
		UserID: getOrderInfo.UserID,
	}

	err := c.canteenRepo.GetOrderInfo(&order)

	return order.ParseToDTOResponseGetOrderInfo(), err
}

func (c *CanteenUseCase) GetOrderList(userID uuid.UUID) ([]dto.ResponseGetOrderList, error) {
	order := new([]entity.Order)

	err := c.canteenRepo.GetOrderList(order, userID)
	if err != nil {
		return nil, err
	}

	parsedOrder := make([]dto.ResponseGetOrderList, len(*order))

	for i, o := range *order {
		parsedOrder[i] = o.ParseToDTOResponseGetOrderList()
	}

	return parsedOrder, err
}

func (c *CanteenUseCase) SoftDeleteMenu(menuID uuid.UUID, userID uuid.UUID) error {
	menu := entity.Menu{
		ID: menuID,
	}

	err := c.canteenRepo.SoftDeleteMenu(&menu, userID)

	return err
}

func (c *CanteenUseCase) SoftDeleteFeedback(feedbackID uuid.UUID, userID uuid.UUID) error {
	feedback := entity.Feedback{
		ID: feedbackID,
	}

	err := c.canteenRepo.SoftDeleteFeedback(&feedback, userID)

	return err
}
