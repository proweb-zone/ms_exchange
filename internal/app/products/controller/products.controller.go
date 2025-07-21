package controller

import (
	"context"
	brands "ms_exchange/internal/app/brands/repository"
	"ms_exchange/internal/app/products/dto/dto_products"
	products "ms_exchange/internal/app/products/repository"
	"ms_exchange/internal/app/products/service"
	"ms_exchange/internal/utils"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductsController struct {
	logger          *zerolog.Logger
	kafkaConsumer   *kafka.Consumer
	productsService *service.ProductsService
}

func GetProductsController(
	logger *zerolog.Logger,
	gormIns *gorm.DB,
	kafkaConsumer *kafka.Consumer,
) *ProductsController {
	productsRepository := products.InitProductsRepository(gormIns)
	productCategoriesB2cRepository := products.InitProductCategoriesB2cRepository(gormIns)
	productCategoriesB2bRepository := products.InitProductCategoriesB2bRepository(gormIns)
	brandsRepository := brands.InitBrandsRepository(gormIns)

	productsService := service.NewProductsService(
		productsRepository,
		productCategoriesB2cRepository,
		productCategoriesB2bRepository,
		brandsRepository,
	)

	productsController := &ProductsController{
		logger:          logger,
		kafkaConsumer:   kafkaConsumer,
		productsService: productsService,
	}

	return productsController
}

func (c *ProductsController) UpdateBatchProducts(data []byte, baseUrl1C string) error {
	var request dto_products.ProductsRequestDto

	start := time.Now()
	if err := utils.DecodeJson(data, &request); err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен товарами завершен с ошибкой")

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// reqUrl, err := url.Parse(baseUrl1C + "/ka_ivankov/hs/toledo_ecxhange/data/group")
	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	// if err != nil {
	// 	return err
	// }

	updateBatchReport, err := c.productsService.CreateOrUpdateBatch(ctx, request)
	if err != nil {
		utils.WriteStdOutErr(c.logger, err, "обмен товарами завершен с ошибкой")

		return err
	}

	exchangeInfo := map[string]any{"Результат": *updateBatchReport, "время обработки": strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)}
	c.logger.Info().Interface("exchangeInfo", exchangeInfo).Msg("Обмен товарами завершен")

	return nil
}
