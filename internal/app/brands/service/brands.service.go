package service

import (
	"context"
	"ms_exchange/internal/app/brands/dto/dto_brands"
	"ms_exchange/internal/app/brands/entity"
	"ms_exchange/internal/app/brands/repository"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type BrandsService struct {
	brandsRepository *repository.BrandsRepository
}

func NewBrandsService(
	brandsRepository *repository.BrandsRepository,
) *BrandsService {
	return &BrandsService{
		brandsRepository: brandsRepository,
	}
}

func (s *BrandsService) CreateOrUpdateBatch(ctx context.Context, request *dto_brands.BrandsRequestDto) (*dto_brands.UpdateBatchResponseDto, error) {
	brandsList, _ := s.brandsRepository.GetList(ctx, []string{"id", "xml"})
	updateBatchReport := &dto_brands.UpdateBatchResponseDto{}

	if len(brandsList) == 0 {
		rowsAffected, err := s.prepareAndCreateBrandsList(ctx, request.Data, &brandsList)
		if err != nil {
			return nil, err
		}

		updateBatchReport.Created = rowsAffected
		return updateBatchReport, err
	}

	existingBrands := s.buildProductIdsMap(brandsList)

	updatedIds := make([]uint, 0, 100)
	duplicateXmls := make(map[string]bool, 100)
	updateBrands := make([]entity.Brands, 0, 1000)
	createBrands := make([]entity.Brands, 0, 1000)

	// подготовка брендов
	for _, item := range request.Data {
		if duplicateXmls[item.Xml] {
			updateBatchReport.Duplicated = append(updateBatchReport.Duplicated, item.Xml)

			continue
		}
		duplicateXmls[item.Xml] = true

		xmlId, err := uuid.Parse(item.Xml)
		if err != nil {
			continue
		}

		id, isExist := existingBrands[item.Xml]
		if !isExist {
			createBrands = append(createBrands, entity.Brands{
				Xml:      xmlId,
				Name:     item.Name,
				IsActive: !item.DeletionMark,
				Sort:     500,
			})

			continue
		} else {
			updatedIds = append(updatedIds, id)
			updateBrands = append(updateBrands, entity.Brands{
				Id:       id,
				Xml:      xmlId,
				Name:     item.Name,
				IsActive: !item.DeletionMark,
				Sort:     500,
			})
		}
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		updateBatchReport.SoftDeleted, _ = s.brandsRepository.ChangeStatus(gctx, false, updatedIds)
		return nil
	})

	if len(updateBrands) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Updated, err = s.brandsRepository.UpdateBatch(gctx, updateBrands)
			return err
		})
	}

	if len(createBrands) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Created, err = s.brandsRepository.CreateBatch(gctx, createBrands)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return updateBatchReport, err
	}

	return updateBatchReport, nil
}

func (s *BrandsService) prepareAndCreateBrandsList(ctx context.Context, request []dto_brands.BrandDto, result *[]entity.Brands) (int64, error) {
	brands := make([]entity.Brands, 0, len(request))
	duplicateXmls := make(map[string]bool, len(request))

	for _, item := range request {
		if duplicateXmls[item.Xml] {
			continue
		}
		duplicateXmls[item.Xml] = true

		xmlId, err := uuid.Parse(item.Xml)
		if err != nil {
			continue
		}

		brands = append(brands, entity.Brands{Xml: xmlId, Name: item.Name, IsActive: !item.DeletionMark})
		*result = append(*result, entity.Brands{Xml: xmlId})
	}

	return s.brandsRepository.CreateBatch(ctx, brands)
}

func (s *BrandsService) buildProductIdsMap(productStoragesList []entity.Brands) map[string]uint {
	existingProductsStorage := make(map[string]uint, len(productStoragesList))

	for _, item := range productStoragesList {
		key := item.Xml.String()
		existingProductsStorage[key] = item.Id
	}

	return existingProductsStorage
}
