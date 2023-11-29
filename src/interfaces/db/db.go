package db

import (
	"fmt"

	"github.com/DiegoSan99/transaction-processor/src/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.AppConfig, log *zap.SugaredLogger) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort,
	)
	log.Info("Connecting to database...", zap.String("dsn", dsn))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.DPanic("failed to connect to database: " + err.Error())
	}

	DB = db
}

func Disconnect() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}

func GetDB() *gorm.DB {
	return DB
}
