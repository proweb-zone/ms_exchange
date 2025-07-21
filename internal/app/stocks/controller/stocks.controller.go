package controller

import (
	"context"

	"ms_exchange/internal/app/stocks/dto/dto_prices"
	"ms_exchange/internal/app/stocks/dto/dto_stocks"
	"ms_exchange/internal/app/stocks/repository"
	"ms_exchange/internal/app/stocks/service"
	"ms_exchange/internal/utils"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type StocksController struct {
	logger               *zerolog.Logger
	kafkaConsumer        *kafka.Consumer
	stockProductsService *service.StockProductsService
	productPricesService *service.ProductPricesService
}

func GetStocksController(
	logger *zerolog.Logger,
	gormIns *gorm.DB,
	kafkaConsumer *kafka.Consumer,
) *StocksController {
	stocksRepository := repository.InitStocksRepository(gormIns)
	stockProductsRepository := repository.InitStockProductsRepository(gormIns)
	priceTypesRepository := repository.InitPriceTypesRepository(gormIns)
	productPrices := repository.InitProductPricesRepository(gormIns)

	stockProductsService := service.NewStocksProductsService(stocksRepository, stockProductsRepository)
	productPricesService := service.NewProductPricesService(priceTypesRepository, productPrices)

	stocksController := &StocksController{
		logger:               logger,
		kafkaConsumer:        kafkaConsumer,
		stockProductsService: stockProductsService,
		productPricesService: productPricesService,
	}

	return stocksController
}

func (c *StocksController) UpdateBatchRemainingStock(data []byte) error {
	var request dto_stocks.StocksRequestDto

	start := time.Now()
	if err := utils.DecodeJson(data, &request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен остатками завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.stockProductsService.CreateOrUpdateBatch(ctx, request)
	if err != nil {
		c.logger.Err(err)

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен остатками завершен")

	return nil
}

func (c *StocksController) UpdateBatchPrices(data []byte) error {
	var request dto_prices.PricesRequestDto

	start := time.Now()
	if err := utils.DecodeJson(data, &request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен ценами завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	updateBatchReport, err := c.productPricesService.CreateOrUpdateBatch(ctx, request)
	if err != nil {
		c.logger.Err(err)

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен ценами завершен")

	return nil
}
