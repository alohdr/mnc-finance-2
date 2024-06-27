package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mnc-finance/entity"
)

var (
	DB *gorm.DB
)

func SetupDatabase() *gorm.DB {
	dsn := "user=username password=password dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Transaction{})
	DB = db
	return db
}
