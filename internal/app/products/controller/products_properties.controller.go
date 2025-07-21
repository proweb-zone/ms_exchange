package controller

import (
	"context"
	"ms_exchange/internal/app/products/dto/dto_properties"
	"ms_exchange/internal/app/products/repository"
	"ms_exchange/internal/app/products/service"
	"ms_exchange/internal/utils"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductPropertiesController struct {
	logger                       *zerolog.Logger
	kafkaConsumer                *kafka.Consumer
	productPropertyTypesService  *service.ProductPropertyTypesService
	productPropertyValuesService *service.ProductPropertyValuesService
}

func GetProductPropertiesController(
	logger *zerolog.Logger,
	gormIns *gorm.DB,
	kafkaConsumer *kafka.Consumer,
) *ProductPropertiesController {
	productPropertyTypesRepository := repository.InitProductPropertyTypesRepository(gormIns)
	productPropertyValuesRepository := repository.InitProductPropertyValuesRepository(gormIns)
	productsRepository := repository.InitProductsRepository(gormIns)

	productPropertyTypesService := service.NewProductPropertyTypesService(productPropertyTypesRepository, productPropertyValuesRepository, productsRepository)
	productPropertyValuesService := service.NewProductPropertyValuesService(productPropertyTypesRepository, productPropertyValuesRepository, productsRepository)

	productCategoriesController := &ProductPropertiesController{
		logger:                       logger,
		kafkaConsumer:                kafkaConsumer,
		productPropertyTypesService:  productPropertyTypesService,
		productPropertyValuesService: productPropertyValuesService,
	}

	return productCategoriesController
}

func (c *ProductPropertiesController) UpdateOrCreateBatchProductPropertyTypes(data []byte) error {
	request := &dto_properties.ProductPropertiesRequestDto{}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productPropertyTypesService.CreateBatchProps(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен свойствами товаров завершен")

	return nil
}

func (c *ProductPropertiesController) CreateBatchProductPropertyValues(data []byte) error {
	request := &dto_properties.ProductPropertyValuesRequestDto{}

	err := c.productPropertyValuesService.Clear()
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productPropertyValuesService.CreateBatchPropValues(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен свойствами товаров завершен")

	return nil
}

func (c *ProductPropertiesController) UpdateOrCreateBatchProductPropertyValues(data []byte) error {
	request := &dto_properties.ProductPropertyValuesRequestDto{}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productPropertyValuesService.UpdateOrCreateBatchPropValues(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен свойствами товаров завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен свойствами товаров завершен")

	return nil
}
