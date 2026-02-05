// Package usecase handles the logic for each user request
package usecase

import (
	"context"

	"github.com/SyafaHadyan/freepass-2026/internal/app/canteen/repository"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/SyafaHadyan/freepass-2026/internal/domain/entity"
	redisitf "github.com/SyafaHadyan/freepass-2026/internal/infra/redis"
	"github.com/google/uuid"
)

type CanteenUseCaseItf interface {
	CreateCanteen(createCanteen dto.CreateCanteen) (dto.ResponseCreateCanteen, error)
	CreateMenu(createMenu dto.CreateMenu) (dto.ResponseCreateMenu, error)
	CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error)
	UpdateMenu(updateMenu dto.UpdateMenu) (dto.ResponseUpdateMenu, error)
	UpdateOrder(updateOrder dto.UpdateOrder, userID uuid.UUID) (dto.ResponseUpdateOrder, error)
	GetCanteenList() ([]dto.ResponseGetCanteenList, error)
	GetCanteenInfo(canteenID uuid.UUID) (dto.ResponseGetCanteenInfo, error)
	GetMenuInfo(menuID uuid.UUID) (dto.ResponseGetMenuInfo, error)
	GetOrderInfo(getOrderInfo dto.GetOrderInfo) (dto.ResponseGetOrderInfo, error)
	GetOrderList(userID uuid.UUID) ([]dto.ResponseGetOrderList, error)
	SoftDeleteMenu(menuID uuid.UUID, userID uuid.UUID) error
}

type CanteenUseCase struct {
	canteenRepo  repository.CanteenDBItf
	redis        redisitf.RedisItf
	redisContext context.Context
}

func NewCanteenUseCase(
	canteenRepo repository.CanteenDBItf, redis redisitf.RedisItf,
) CanteenUseCaseItf {
	return &CanteenUseCase{
		canteenRepo:  canteenRepo,
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

	return canteen.ParseToDTORResponseCreateCanteen(), err
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
		parsedCanteen[i] = c.ParseToDTORResponseGetCanteenList()
	}

	return parsedCanteen, err
}

func (c *CanteenUseCase) GetCanteenInfo(canteenID uuid.UUID) (dto.ResponseGetCanteenInfo, error) {
	canteen := entity.Canteen{
		ID: canteenID,
	}

	err := c.canteenRepo.GetCanteenInfo(&canteen)

	return canteen.ParseToDTORResponseGetCanteenInfo(), err
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
