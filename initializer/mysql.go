package initializer

import (
	"fmt"
	"store/config"
	"store/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MysqlInit() (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)

	return gorm.Open(mysql.Open(url), gormConfig)
}

func MysqlMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Product{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}