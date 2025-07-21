package service

import (
	"context"
	"ms_exchange/internal/app/products/dto/dto_properties"
	"ms_exchange/internal/app/products/repository"
)

type RelatedPropsData struct {
	ProductsXmlsMap      map[string]uint
	propertyTypesXmlsMap map[string]bool
}

type ProductPropertyTypesService struct {
	productPropertyTypesRepository *repository.ProductPropertyTypesRepository
	productPropertyValues          *repository.ProductPropertyValuesRepository
	productsRepository             *repository.ProductsRepository
}

func NewProductPropertyTypesService(
	productPropertyTypesRepository *repository.ProductPropertyTypesRepository,
	productPropertyValues *repository.ProductPropertyValuesRepository,
	productsRepository *repository.ProductsRepository,
) *ProductPropertyTypesService {
	return &ProductPropertyTypesService{
		productPropertyTypesRepository: productPropertyTypesRepository,
		productPropertyValues:          productPropertyValues,
		productsRepository:             productsRepository,
	}
}

func (s *ProductPropertyTypesService) CreateBatchProps(ctx context.Context, request *dto_properties.ProductPropertiesRequestDto) (*dto_properties.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_properties.UpdateBatchResponseDto{}

	propertyTypes := prepareProps(request.Data)

	err := s.productPropertyTypesRepository.ClearProperties()
	if err != nil {
		return updateBatchReport, err
	}

	rowsAffected, err := s.productPropertyTypesRepository.CreateBatch(ctx, propertyTypes)
	if err != nil {
		return updateBatchReport, err
	}

	updateBatchReport.Created = rowsAffected
	return updateBatchReport, nil
}
