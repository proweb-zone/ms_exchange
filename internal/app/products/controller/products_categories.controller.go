package controller

import (
	"context"
	"ms_exchange/internal/app/products/dto/dto_categories"
	"ms_exchange/internal/app/products/repository"
	"ms_exchange/internal/app/products/service"
	"ms_exchange/internal/utils"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductCategoriesController struct {
	logger                   *zerolog.Logger
	kafkaConsumer            *kafka.Consumer
	productCategoriesService *service.ProductCategoriesService
}

func GetProductCategoriesController(
	logger *zerolog.Logger,
	gormIns *gorm.DB,
	kafkaConsumer *kafka.Consumer,
) *ProductCategoriesController {
	productCategoriesB2cRepository := repository.InitProductCategoriesB2cRepository(gormIns)
	productCategoriesB2bRepository := repository.InitProductCategoriesB2bRepository(gormIns)

	productCategoriesService := service.NewProductCategoriesService(productCategoriesB2cRepository, productCategoriesB2bRepository)

	productCategoriesController := &ProductCategoriesController{
		logger:                   logger,
		kafkaConsumer:            kafkaConsumer,
		productCategoriesService: productCategoriesService,
	}

	return productCategoriesController
}

func (c *ProductCategoriesController) UpdateBatchProductCategoriesB2c(data []byte) error {
	request := &dto_categories.ProductCategoriesRequestDto{}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен категориями товаров завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productCategoriesService.CreateOrUpdateBatchB2c(ctx, request)
	if err != nil {
		c.logger.Err(err)

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен категориями B2C завершен")

	return nil
}

func (c *ProductCategoriesController) UpdateBatchProductCategoriesB2b(data []byte) error {
	request := &dto_categories.ProductCategoriesRequestDto{}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен категориями товаров завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productCategoriesService.CreateOrUpdateBatchB2b(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен категориями товаров завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен категориями B2B завершен")

	return nil
}
