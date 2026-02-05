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
	UpdateMenu(updateMenu dto.UpdateMenu) (dto.ResponseUpdateMenu, error)
	GetCanteenList() ([]dto.ResponseGetCanteenList, error)
	GetCanteenInfo(canteenID uuid.UUID) (dto.ResponseGetCanteenInfo, error)
	GetMenuInfo(menuID uuid.UUID) (dto.ResponseGetMenuInfo, error)
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

func (c *CanteenUseCase) UpdateMenu(updateMenu dto.UpdateMenu) (dto.ResponseUpdateMenu, error) {
	menu := entity.Menu{
		ID:    updateMenu.ID,
		Name:  updateMenu.Name,
		Price: updateMenu.Price,
	}

	err := c.canteenRepo.UpdateMenu(&menu, updateMenu.UserID)

	return menu.ParseToDTOResponseUpdateMenu(), err
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

func (c *CanteenUseCase) SoftDeleteMenu(menuID uuid.UUID, userID uuid.UUID) error {
	menu := entity.Menu{
		ID: menuID,
	}

	err := c.canteenRepo.SoftDeleteMenu(&menu, userID)

	return err
}
