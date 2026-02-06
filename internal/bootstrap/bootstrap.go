// Package bootstrap starts the backend service, sets the backend configuration, and connects to external services
package bootstrap

import (
	"log"
	"time"

	canteenhandler "github.com/SyafaHadyan/freepass-2026/internal/app/canteen/interface/rest"
	canteenrepository "github.com/SyafaHadyan/freepass-2026/internal/app/canteen/repository"
	canteenusecase "github.com/SyafaHadyan/freepass-2026/internal/app/canteen/usecase"
	userhandler "github.com/SyafaHadyan/freepass-2026/internal/app/user/interface/rest"
	userrepository "github.com/SyafaHadyan/freepass-2026/internal/app/user/repository"
	userusecase "github.com/SyafaHadyan/freepass-2026/internal/app/user/usecase"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/db"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	fiberapp "github.com/SyafaHadyan/freepass-2026/internal/infra/fiber"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/jwt"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/payment"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/redis"
	"github.com/SyafaHadyan/freepass-2026/internal/middleware"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Bootstrap struct {
	App       *fiberapp.Fiber
	Config    *env.Env
	Validator *validator.Validate
	Database  *gorm.DB
	Redis     *redis.Redis
	JWT       *jwt.JWT
}

func Start() *Bootstrap {
	log.Println("starting app")
	startTime := time.Now()

	config := env.New()

	validator := validator.New()

	database := db.New(config)

	redis := redis.New(config)

	jwt := jwt.New(config)

	payment := payment.New(config)

	app := fiberapp.New(config)

	middleware := middleware.NewMiddleware(*jwt)

	userRepository := userrepository.NewUserDB(database)
	canteenRepository := canteenrepository.NewCanteenDB(database)

	userUseCase := userusecase.NewUserUseCase(userRepository, jwt, redis)
	canteenUseCase := canteenusecase.NewCanteenUseCase(canteenRepository, payment, redis)

	userhandler.NewUserHandler(app.Router, validator, middleware, userUseCase, config)
	canteenhandler.NewCanteenHandler(app.Router, validator, middleware, canteenUseCase, config)

	Bootstrap := Bootstrap{
		App:       app,
		Config:    config,
		Validator: validator,
		Database:  database,
		Redis:     redis,
		JWT:       jwt,
	}

	log.Printf("startup time: %v", time.Since(startTime))

	return &Bootstrap
}
