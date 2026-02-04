// Package bootstrap starts the backend service, sets the backend configuration, and connects to external services
package bootstrap

import (
	"log"
	"time"

	userhandler "github.com/SyafaHadyan/freepass-2026/internal/app/user/interface/rest"
	userrepository "github.com/SyafaHadyan/freepass-2026/internal/app/user/repository"
	userusecase "github.com/SyafaHadyan/freepass-2026/internal/app/user/usecase"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/db"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	fiberapp "github.com/SyafaHadyan/freepass-2026/internal/infra/fiber"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/jwt"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/mailer"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/redis"
	"github.com/SyafaHadyan/freepass-2026/internal/infra/s3"
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
	Mailer    *mailer.Mailer
	S3        *s3.S3
}

func Start() *Bootstrap {
	log.Println("starting app")
	startTime := time.Now()

	config := env.New()

	validator := validator.New()

	database := db.New(config)

	redis := redis.New(config)

	jwt := jwt.New(config)

	mailer := mailer.New(config)

	s3 := s3.New(config)

	app := fiberapp.New(config)

	middleware := middleware.NewMiddleware(*jwt)

	userRepository := userrepository.NewUserDB(database)

	userUseCase := userusecase.NewUserUseCase(userRepository, jwt, redis)

	userhandler.NewUserHandler(app.Router, validator, middleware, userUseCase, config, mailer)

	Bootstrap := Bootstrap{
		App:       app,
		Config:    config,
		Validator: validator,
		Database:  database,
		Redis:     redis,
		JWT:       jwt,
		Mailer:    mailer,
		S3:        s3,
	}

	log.Printf("startup time: %v", time.Since(startTime))

	return &Bootstrap
}
