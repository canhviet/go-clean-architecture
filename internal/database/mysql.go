package database

import (
	"fmt"

	"github.com/canhviet/go-clean-architecture/internal/config"
	"github.com/canhviet/go-clean-architecture/internal/model"
	"github.com/canhviet/go-clean-architecture/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
    )

	var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        SkipDefaultTransaction: true,
        PrepareStmt:            true,
    })
    if err != nil {
        return err
    }

    if logger.Log != nil {
        logger.Log.Info("MySQL connected successfully")
    } else {
        fmt.Println("MySQL connected successfully (logger not init yet)")
    }

    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

    if err := DB.AutoMigrate(&model.User{}); err != nil {
        logger.Log.Error("AutoMigrate failed", zap.Error(err))
    } else {
        logger.Log.Info("AutoMigrate completed successfully")
    }

    logger.Log.Info("MySQL connected successfully")
    return nil
}