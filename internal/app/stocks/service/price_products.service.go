package service

import (
	"context"

	"ms_exchange/internal/app/stocks/dto/dto_prices"
	"ms_exchange/internal/app/stocks/entity"
	"ms_exchange/internal/app/stocks/repository"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ProductPricesService struct {
	priceTypesRepository    *repository.PriceTypesRepository
	productPricesRepository *repository.ProductPricesRepository
}

func NewProductPricesService(
	priceTypesRepository *repository.PriceTypesRepository,
	productPricesRepository *repository.ProductPricesRepository,
) *ProductPricesService {
	return &ProductPricesService{
		priceTypesRepository:    priceTypesRepository,
		productPricesRepository: productPricesRepository,
	}
}

func (s *ProductPricesService) CreateOrUpdateBatch(ctx context.Context, request dto_prices.PricesRequestDto) (*dto_prices.UpdateBatchResponseDto, error) {
	priceTypesList, _ := s.priceTypesRepository.GetList(ctx, []string{"xml", "is_active"})
	updateBatchReport := &dto_prices.UpdateBatchResponseDto{}

	// создание типов цен
	if len(priceTypesList) == 0 {
		if rowsAffected, err := s.prepareAndCreatePriceTypesList(ctx, request.General.PriceTypes, &priceTypesList); rowsAffected == 0 || err != nil {
			return updateBatchReport, err
		}
	}

	productPricesList, _ := s.productPricesRepository.GetList(ctx, []string{"id", "product_xml", "price_xml"})
	// маппинг статусов типов цен
	existingTypes := buildLookupPriceTypesActiveMap(priceTypesList)

	// создание товаров с ценами
	if len(productPricesList) == 0 {
		rowsAffected, err := s.prepareAndCreatePriceProducts(ctx, request, &productPricesList, existingTypes)
		updateBatchReport.Created = rowsAffected

		return updateBatchReport, err
	}

	// *** обновление ***
	existingPriceTypes := s.buildProductIdsMap(productPricesList)

	// метка для дективации (true - деактивировать)
	var IsDeletionMark bool
	updatedIds := make([]uint, 0, 100)
	duplicateIds := make(map[uint]bool, 100)
	duplicateXmls := make(map[uuid.UUID]bool, 100)
	updateProductPrices := make([]entity.ProductPrices, 0, 1000)
	createProductPrices := make([]entity.ProductPrices, 0, 1000)

	// обновление типов цен
	storagesToUpdate := s.buildPriceTypes(request.General.PriceTypes, &duplicateXmls)
	if len(storagesToUpdate) > 0 {
		if _, err := s.priceTypesRepository.UpdateBatch(ctx, storagesToUpdate); err != nil {
			return updateBatchReport, err
		}
	}

	// обновление товаров
	for _, item := range request.Data {
		for _, itemPrice := range item.Prices {
			id, isExistType := existingPriceTypes[item.ProductXml.String()+itemPrice.TypePriceXml.String()]

			if isExistType && id == 0 {
				createProductPrices = append(createProductPrices, entity.ProductPrices{
					ProductXml: item.ProductXml,
					PriceXml:   itemPrice.TypePriceXml,
					Price:      itemPrice.Price,
					IsActive:   !IsDeletionMark,
				})

				continue
			}

			_, isExist := existingTypes[itemPrice.TypePriceXml.String()]
			if !isExist || duplicateIds[id] {
				continue
			}
			duplicateIds[id] = true

			if duplicateXmls[itemPrice.TypePriceXml] {
				IsDeletionMark = true
			} else {
				IsDeletionMark = false
			}

			updatedIds = append(updatedIds, id)
			updateProductPrices = append(updateProductPrices, entity.ProductPrices{
				Id:         id,
				ProductXml: item.ProductXml,
				PriceXml:   itemPrice.TypePriceXml,
				Price:      itemPrice.Price,
				IsActive:   !IsDeletionMark,
			})
		}
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		updateBatchReport.SoftDeleted, _ = s.productPricesRepository.ChangeStatus(gctx, false, updatedIds)
		return nil
	})

	if len(updateProductPrices) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Updated, err = s.productPricesRepository.UpdateBatch(gctx, updateProductPrices)
			return err
		})
	}

	if len(createProductPrices) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Created, err = s.productPricesRepository.CreateBatch(gctx, createProductPrices)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return updateBatchReport, err
	}

	return updateBatchReport, nil
}

func (s *ProductPricesService) buildProductIdsMap(productPricesList []entity.ProductPrices) map[string]uint {
	existingProductsPrice := make(map[string]uint, len(productPricesList))

	for _, item := range productPricesList {
		key := item.ProductXml.String() + item.PriceXml.String()
		existingProductsPrice[key] = item.Id
	}

	return existingProductsPrice
}

func (s *ProductPricesService) prepareAndCreatePriceTypesList(ctx context.Context, priceTypes []dto_prices.PriceTypeDto, result *[]entity.PriceTypes) (int64, error) {
	stocks := make([]entity.PriceTypes, 0, len(priceTypes))
	duplicateXmls := make(map[uuid.UUID]bool, len(priceTypes))

	for _, item := range priceTypes {
		if duplicateXmls[item.Xml] {
			continue
		}
		duplicateXmls[item.Xml] = true

		stocks = append(stocks, entity.PriceTypes{Xml: item.Xml, Name: item.Name, IsActive: !item.DeletionMark})
		*result = append(*result, entity.PriceTypes{Xml: item.Xml, IsActive: !item.DeletionMark})
	}

	return s.priceTypesRepository.CreateBatch(ctx, stocks)
}

func (s *ProductPricesService) prepareAndCreatePriceProducts(
	ctx context.Context,
	pricesRequest dto_prices.PricesRequestDto,
	result *[]entity.ProductPrices,
	existingPrice map[string]bool,
) (int64, error) {
	duplicateXmls := make(map[string]bool, len(pricesRequest.Data))

	for _, item := range pricesRequest.Data {
		if duplicateXmls[item.ProductXml.String()] {
			continue
		}
		duplicateXmls[item.ProductXml.String()] = true

		for _, itemProduct := range item.Prices {
			isActive, isExist := existingPrice[itemProduct.TypePriceXml.String()]

			if !isExist {
				continue
			}

			*result = append(*result, entity.ProductPrices{
				ProductXml: item.ProductXml,
				PriceXml:   itemProduct.TypePriceXml,
				Price:      itemProduct.Price,
				IsActive:   isActive,
			})
		}
	}

	return s.productPricesRepository.CreateBatch(ctx, *result)
}

func (s *ProductPricesService) buildPriceTypes(
	priceTypes []dto_prices.PriceTypeDto,
	duplicateXmls *map[uuid.UUID]bool,
) []entity.PriceTypes {
	priceTypesToUpdate := []entity.PriceTypes{}

	for _, priceType := range priceTypes {
		IsDeletionMark := priceType.DeletionMark

		if IsDeletionMark {
			(*duplicateXmls)[priceType.Xml] = true
		}

		priceTypesToUpdate = append(priceTypesToUpdate, entity.PriceTypes{
			Xml:      priceType.Xml,
			Name:     priceType.Name,
			IsActive: !IsDeletionMark,
		})
	}

	return priceTypesToUpdate
}

func buildLookupPriceTypesActiveMap(slice []entity.PriceTypes) map[string]bool {
	mappedSlice := make(map[string]bool, len(slice))

	for _, key := range slice {
		mappedSlice[key.Xml.String()] = key.IsActive
	}

	return mappedSlice
}
