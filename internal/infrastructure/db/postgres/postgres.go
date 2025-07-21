package postgres

import (
	"context"
	"errors"
	"log"
	stocksEntity "ms_exchange/internal/app/stocks/entity"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Connect(gormIns *gorm.DB) error {
	sqlDB, err := gormIns.DB()

	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func Close(gormIns *gorm.DB) error {
	if gormIns == nil {
		return errors.New("экземпляр ORM не получен")
	}

	sqlDB, _ := gormIns.DB()

	return sqlDB.Close()
}

func Create(dsn string) (*gorm.DB, error) {
	conf := gormLogger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  gormLogger.Warn,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		Colorful:                  true,
	}
	newGormLogger := gormLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), conf)

	return gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			},
		),
		&gorm.Config{
			Logger:                 newGormLogger,
			SkipDefaultTransaction: true,
			PrepareStmt:            false,
		},
	)
}

/*
порядок следования моделей в автомиграторе
критичен для создания связанных таблиц
*/
func EnableAutoMigrate(gormIns *gorm.DB) error {
	return gormIns.AutoMigrate(
		// склады
		stocksEntity.Storages{},
		stocksEntity.ProductStorages{},
		stocksEntity.ProductPrices{},
		stocksEntity.PriceTypes{},

		// товары
		// brandsEntity.Brands{},
		// productsEntity.ProductCategoriesB2c{},
		// productsEntity.ProductCategoriesB2b{},
		// productsEntity.Products{},
		// productsEntity.ProductPropertyTypes{},
		// productsEntity.ProductPropertyValues{},
	)
}
