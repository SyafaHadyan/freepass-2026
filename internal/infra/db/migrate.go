package db

import (
	"log"

	"github.com/SyafaHadyan/freepass-2026/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		entity.User{},
		entity.UserDetail{},
		entity.Canteen{},
		entity.Menu{},
		entity.Order{},
		entity.Payment{},
		entity.Feedback{},
	)
	if err != nil {
		log.Panic("database migration failed")
	}

	log.Println("database migration complete")
}
