package db

import (
	"fmt"
	"ms_exchange/internal/config"
	"ms_exchange/internal/infrastructure/db/postgres"
	"ms_exchange/pkg/logger"

	"gorm.io/gorm"
)

type DbState struct {
	GormIns *gorm.DB
}

func (d *DbState) InitDb(сfgEnv *config.Config, logger *logger.Logger) *DbState {
	d.GormIns = initGormIns(сfgEnv, logger)

	return d
}

func initGormIns(сfgEnv *config.Config, logger *logger.Logger) *gorm.DB {
	gormIns, err := postgres.Create(buildDbConnectUrl(сfgEnv))

	if err != nil {
		logger.Fatal(err, "Экземпляр ORM не был создан", nil)
	}

	err = postgres.EnableAutoMigrate(gormIns)
	if err != nil {
		logger.Fatal(err, "Ошибка создания миграции БД", nil)
	}

	err = postgres.Connect(gormIns)

	if err != nil {
		logger.Fatal(err, "ошибка создания соединения с БД", nil)
		return nil
	}

	return gormIns
}

func buildDbConnectUrl(сfgEnv *config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		сfgEnv.Db.Driver,
		сfgEnv.Db.User,
		сfgEnv.Db.Password,
		сfgEnv.Db.Host,
		сfgEnv.Db.Port,
		сfgEnv.Db.Name,
		сfgEnv.Db.Option,
	)
}
