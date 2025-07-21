package service

import (
	"context"
	"ms_exchange/internal/app/products/dto/dto_categories"
	"ms_exchange/internal/app/products/entity"
	"ms_exchange/internal/app/products/repository"

	"github.com/google/uuid"
)

type ProductCategoriesService struct {
	productCategoriesRepositoryB2c *repository.ProductCategoriesB2cRepository
	productCategoriesRepositoryB2b *repository.ProductCategoriesB2bRepository
}

func NewProductCategoriesService(
	productCategoriesRepositoryB2c *repository.ProductCategoriesB2cRepository,
	productCategoriesRepositoryB2b *repository.ProductCategoriesB2bRepository,
) *ProductCategoriesService {
	return &ProductCategoriesService{
		productCategoriesRepositoryB2c: productCategoriesRepositoryB2c,
		productCategoriesRepositoryB2b: productCategoriesRepositoryB2b,
	}
}

func (s *ProductCategoriesService) CreateOrUpdateBatchB2c(ctx context.Context, request *dto_categories.ProductCategoriesRequestDto) (*dto_categories.UpdateBatchResponseDto, error) {
	if request.General.Clean {
		err := s.productCategoriesRepositoryB2c.ClearCategories()
		if err != nil {
			return nil, err
		}
	}

	return s.prepareAndCreateCategoriesB2c(ctx, request.Data)
}

func (s *ProductCategoriesService) CreateOrUpdateBatchB2b(ctx context.Context, request *dto_categories.ProductCategoriesRequestDto) (*dto_categories.UpdateBatchResponseDto, error) {
	if request.General.Clean {
		err := s.productCategoriesRepositoryB2b.ClearCategories()
		if err != nil {
			return nil, err
		}
	}

	return s.prepareAndCreateCategoriesB2b(ctx, request.Data)
}

func (s *ProductCategoriesService) prepareAndCreateCategoriesB2b(ctx context.Context, categories []dto_categories.CategoryDto) (*dto_categories.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_categories.UpdateBatchResponseDto{}
	duplicateXmls := map[string]bool{}
	productCategory := &entity.ProductCategoriesB2b{}

	createCategory(categories, productCategory, &duplicateXmls, &updateBatchReport.Duplicated)
	rows, err := s.productCategoriesRepositoryB2b.CreateBatch(ctx, productCategory)
	if err != nil {
		return updateBatchReport, err
	}

	updateBatchReport.Created = rows
	return updateBatchReport, nil
}

func (s *ProductCategoriesService) prepareAndCreateCategoriesB2c(ctx context.Context, categories []dto_categories.CategoryDto) (*dto_categories.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_categories.UpdateBatchResponseDto{}
	duplicateXmls := map[string]bool{}
	productCategory := &entity.ProductCategoriesB2c{}

	createCategory(categories, productCategory, &duplicateXmls, &updateBatchReport.Duplicated)
	rows, err := s.productCategoriesRepositoryB2c.CreateBatch(ctx, productCategory)
	if err != nil {
		return updateBatchReport, err
	}

	updateBatchReport.Created = rows
	return updateBatchReport, nil
}

func createCategory(
	categories []dto_categories.CategoryDto,
	productCategories any,
	duplicateXmls *map[string]bool,
	duplicated *[]string,
) {
	for _, item := range categories {
		if (*duplicateXmls)[item.Xml] {
			*duplicated = append(*duplicated, item.Xml)
			continue
		}
		(*duplicateXmls)[item.Xml] = true

		switch categoryItem := productCategories.(type) {
		case *entity.ProductCategoriesB2b:
			if item.Parent == "" {
				xmlId, err := uuid.Parse(item.Xml)
				if err != nil {
					continue
				}

				categoryItem.Xml = xmlId
				categoryItem.Name = item.Name
				categoryItem.IsActive = !item.DeletionMark
			} else {
				xmlId, err := uuid.Parse(item.Xml)
				if err != nil {
					continue
				}
				categoryItem.Children = append(categoryItem.Children, entity.ProductCategoriesB2b{
					Xml:      xmlId,
					Name:     item.Name,
					IsActive: !item.DeletionMark,
				})
			}
		case *entity.ProductCategoriesB2c:
			if item.Parent == "" {
				xmlId, err := uuid.Parse(item.Xml)
				if err != nil {
					continue
				}

				categoryItem.Xml = xmlId
				categoryItem.Name = item.Name
				categoryItem.IsActive = !item.DeletionMark
			} else {
				xmlId, err := uuid.Parse(item.Xml)
				if err != nil {
					continue
				}
				categoryItem.Children = append(categoryItem.Children, entity.ProductCategoriesB2c{
					Xml:      xmlId,
					Name:     item.Name,
					IsActive: !item.DeletionMark,
				})
			}
		}
	}
}
