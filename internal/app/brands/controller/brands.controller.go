package controller

import (
	"context"
	"ms_exchange/internal/app/brands/dto/dto_brands"
	"ms_exchange/internal/app/brands/repository"
	"ms_exchange/internal/app/brands/service"
	"ms_exchange/internal/utils"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BrandsController struct {
	logger        *zerolog.Logger
	kafkaConsumer *kafka.Consumer
	brandsService *service.BrandsService
}

func GetBrandsController(
	logger *zerolog.Logger,
	gormIns *gorm.DB,
	kafkaConsumer *kafka.Consumer,
) *BrandsController {
	brandsRepository := repository.InitBrandsRepository(gormIns)
	brandsService := service.NewBrandsService(brandsRepository)

	productCategoriesController := &BrandsController{
		logger:        logger,
		kafkaConsumer: kafkaConsumer,
		brandsService: brandsService,
	}

	return productCategoriesController
}

func (c *BrandsController) UpdateBatchBrands(data []byte) error {
	request := &dto_brands.BrandsRequestDto{}

	start := time.Now()
	if err := utils.DecodeJson(data, request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен брендами завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.brandsService.CreateOrUpdateBatch(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен брендами завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен брендами завершен")

	return nil
}
