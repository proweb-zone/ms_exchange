package service

import (
	"context"

	"ms_exchange/internal/app/stocks/dto/dto_stocks"
	"ms_exchange/internal/app/stocks/entity"
	"ms_exchange/internal/app/stocks/repository"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type StockProductsService struct {
	stocksRepository        *repository.StocksRepository
	stockProductsRepository *repository.StockProductsRepository
}

func NewStocksProductsService(
	stocksRepository *repository.StocksRepository,
	stockProductsRepository *repository.StockProductsRepository,
) *StockProductsService {
	return &StockProductsService{
		stocksRepository:        stocksRepository,
		stockProductsRepository: stockProductsRepository,
	}
}

func (s *StockProductsService) CreateOrUpdateBatch(ctx context.Context, request dto_stocks.StocksRequestDto) (*dto_stocks.UpdateBatchResponseDto, error) {
	stocksList, _ := s.stocksRepository.GetList(ctx, []string{"xml", "is_active"})
	updateBatchReport := &dto_stocks.UpdateBatchResponseDto{}

	// создание складов
	if len(stocksList) == 0 {
		if rowsAffected, err := s.prepareAndCreateStocks(ctx, request.General.Storages, &stocksList); rowsAffected == 0 || err != nil {
			return updateBatchReport, err
		}
	}

	productStoragesList, _ := s.stockProductsRepository.GetList(ctx, []string{"id", "product_xml", "storage_xml"})

	// маппинг статусов типов цен
	existingStorage := buildLookupStoragesActiveMap(stocksList)

	// создание товаров с остатками
	if len(productStoragesList) == 0 {
		_, err := s.prepareAndCreateStockProducts(ctx, request.Data, &productStoragesList, existingStorage)
		return updateBatchReport, err
	}

	// *** обновление ***
	existingProductsStorage := s.buildProductIdsMap(productStoragesList)

	// метка для дективации (true - деактивировать)
	var IsDeletionMark bool
	updatedIds := make([]uint, 0, 100)
	duplicateIds := make(map[uint]bool, 100)
	duplicateXmls := make(map[uuid.UUID]bool, 100)
	updateProductStorages := make([]entity.ProductStorages, 0, 1000)
	createProductStorages := make([]entity.ProductStorages, 0, 1000)

	// обновление складов
	storagesToUpdate := s.buildStocks(request.General.Storages, &duplicateXmls)
	if len(storagesToUpdate) > 0 {
		if _, err := s.stocksRepository.UpdateBatch(ctx, storagesToUpdate); err != nil {
			return updateBatchReport, err
		}
	}

	// обновление товаров
	for _, item := range request.Data {
		for _, itemStorage := range item.Storages {
			id := existingProductsStorage[item.ProductXml.String()+itemStorage.StorageXml.String()]

			if id == 0 {
				createProductStorages = append(createProductStorages, entity.ProductStorages{
					ProductXml: item.ProductXml,
					StorageXml: itemStorage.StorageXml,
					Quantity:   itemStorage.Quantity,
					IsActive:   !IsDeletionMark,
				})

				continue
			}

			_, isExist := existingStorage[itemStorage.StorageXml.String()]
			if !isExist || duplicateIds[id] {
				continue
			}
			duplicateIds[id] = true

			if duplicateXmls[itemStorage.StorageXml] {
				IsDeletionMark = true
			} else {
				IsDeletionMark = false
			}

			updatedIds = append(updatedIds, id)
			updateProductStorages = append(updateProductStorages, entity.ProductStorages{
				Id:         id,
				ProductXml: item.ProductXml,
				StorageXml: itemStorage.StorageXml,
				Quantity:   itemStorage.Quantity,
				IsActive:   !IsDeletionMark,
			})
		}
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		updateBatchReport.SoftDeleted, _ = s.stockProductsRepository.ChangeStatus(gctx, false, updatedIds)
		return nil
	})

	if len(updateProductStorages) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Updated, err = s.stockProductsRepository.UpdateBatch(gctx, updateProductStorages)
			return err
		})
	}

	if len(createProductStorages) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Created, err = s.stockProductsRepository.CreateBatch(gctx, createProductStorages)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return updateBatchReport, err
	}

	return updateBatchReport, nil
}

func (s *StockProductsService) buildProductIdsMap(productStoragesList []entity.ProductStorages) map[string]uint {
	existingProductsStorage := make(map[string]uint, len(productStoragesList))

	for _, item := range productStoragesList {
		key := item.ProductXml.String() + item.StorageXml.String()
		existingProductsStorage[key] = item.Id
	}

	return existingProductsStorage
}

func (s *StockProductsService) prepareAndCreateStocks(ctx context.Context, storages []dto_stocks.StorageDto, result *[]entity.Storages) (int64, error) {
	stocks := make([]entity.Storages, 0, len(storages))
	duplicateXmls := make(map[uuid.UUID]bool, len(storages))

	for _, item := range storages {
		if duplicateXmls[item.Xml] {
			continue
		}
		duplicateXmls[item.Xml] = true

		stocks = append(stocks, entity.Storages{Xml: item.Xml, Name: item.Name, IsActive: !item.DeletionMark})
		*result = append(*result, entity.Storages{Xml: item.Xml, IsActive: !item.DeletionMark})
	}

	return s.stocksRepository.CreateBatch(ctx, stocks)
}

func (s *StockProductsService) prepareAndCreateStockProducts(
	ctx context.Context,
	data []dto_stocks.DataDto,
	result *[]entity.ProductStorages,
	existingStorage map[string]bool,
) (int64, error) {
	duplicateXmls := make(map[string]bool, len(data))

	for _, item := range data {
		if duplicateXmls[item.ProductXml.String()] {
			continue
		}
		duplicateXmls[item.ProductXml.String()] = true

		for _, itemProduct := range item.Storages {
			isActive, isExist := existingStorage[itemProduct.StorageXml.String()]

			if !isExist {
				continue
			}

			*result = append(*result, entity.ProductStorages{
				ProductXml: item.ProductXml,
				StorageXml: itemProduct.StorageXml,
				Quantity:   itemProduct.Quantity,
				IsActive:   isActive,
			})
		}
	}

	return s.stockProductsRepository.CreateBatch(ctx, *result)
}

func (s *StockProductsService) buildStocks(
	storages []dto_stocks.StorageDto,
	duplicateXmls *map[uuid.UUID]bool,
) []entity.Storages {
	storagesToUpdate := []entity.Storages{}

	for _, storage := range storages {
		IsDeletionMark := storage.DeletionMark

		if IsDeletionMark {
			(*duplicateXmls)[storage.Xml] = true
		}

		storagesToUpdate = append(storagesToUpdate, entity.Storages{
			Xml:      storage.Xml,
			Name:     storage.Name,
			IsActive: !IsDeletionMark,
		})
	}

	return storagesToUpdate
}

func buildLookupStoragesActiveMap(slice []entity.Storages) map[string]bool {
	mappedSlice := make(map[string]bool, len(slice))

	for _, key := range slice {
		mappedSlice[key.Xml.String()] = key.IsActive
	}

	return mappedSlice
}
