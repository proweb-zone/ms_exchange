package service

import (
	"context"
	"ms_exchange/internal/app/products/dto/dto_products"
	"ms_exchange/internal/app/products/dto/dto_properties"
	"ms_exchange/internal/app/products/entity"
	"ms_exchange/internal/app/products/repository"
	"ms_exchange/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ProductPropertyValuesService struct {
	productPropertyTypesRepository *repository.ProductPropertyTypesRepository
	productPropertyValues          *repository.ProductPropertyValuesRepository
	productsRepository             *repository.ProductsRepository
}

func NewProductPropertyValuesService(
	productPropertyTypesRepository *repository.ProductPropertyTypesRepository,
	productPropertyValues *repository.ProductPropertyValuesRepository,
	productsRepository *repository.ProductsRepository,
) *ProductPropertyValuesService {
	return &ProductPropertyValuesService{
		productPropertyTypesRepository: productPropertyTypesRepository,
		productPropertyValues:          productPropertyValues,
		productsRepository:             productsRepository,
	}
}

func (s *ProductPropertyValuesService) CreateBatchPropValues(ctx context.Context, request *dto_properties.ProductPropertyValuesRequestDto) (*dto_properties.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_properties.UpdateBatchResponseDto{}
	productXmlIds := make([]uuid.UUID, 0, len(request.Data))
	duplicateXmls := make(map[uuid.UUID]bool, len(request.Data))

	for _, prop := range request.Data {
		if duplicateXmls[prop.ProductXml] {
			continue
		}
		duplicateXmls[prop.ProductXml] = true

		productXmlIds = append(productXmlIds, prop.ProductXml)
	}

	relatedData, err := s.getRelatedData(ctx, productXmlIds)
	if err != nil {
		return updateBatchReport, err
	}

	propertyValues := preparePropValues(request.Data, relatedData.ProductsXmlsMap, relatedData.propertyTypesXmlsMap)

	rowsAffected, err := s.productPropertyValues.CreateBatch(ctx, propertyValues)
	if err != nil {
		return updateBatchReport, err
	}

	updateBatchReport.Created = rowsAffected
	return updateBatchReport, nil
}

func (s *ProductPropertyValuesService) UpdateOrCreateBatchPropValues(ctx context.Context, request *dto_properties.ProductPropertyValuesRequestDto) (*dto_properties.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_properties.UpdateBatchResponseDto{}
	// xml товаров для обновления свойств
	productXmls := make([]uuid.UUID, 0, len(request.Data))
	// xml кодов свойств для обновления значений
	propertyValueXmls := make([]uuid.UUID, 0, len(request.Data))

	// сбор входящих данных
	updateBatchReport.Duplicated = buildPropetyValueXmls(request.Data, &productXmls, &propertyValueXmls)

	// получение существующих значений свойств и подготовка данных (для обновления)
	existingPropertiesMap, err := s.getExistingPropValuesMap(ctx, &propertyValueXmls)
	if err != nil {
		return updateBatchReport, err
	}

	// получение существующих товаров и типов свойств
	relatedData, err := s.getRelatedData(ctx, productXmls)
	if err != nil {
		return updateBatchReport, err
	}

	// подготовка сущностей
	propertyValues := preparePropValues(request.Data, relatedData.ProductsXmlsMap, relatedData.propertyTypesXmlsMap)
	prepared := buildPropertyValuesEntity(propertyValues, existingPropertiesMap)

	// обновление/создание
	err = s.updateOrCreatePropertyValues(ctx, prepared, updateBatchReport)
	if err != nil {
		return updateBatchReport, err
	}

	return updateBatchReport, nil
}

func (s *ProductPropertyValuesService) getRelatedData(ctx context.Context, productXmlIds []uuid.UUID) (*RelatedPropsData, error) {
	relatedData := &RelatedPropsData{}
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		productsXmls, err := s.productsRepository.GetListByXmlIds(gctx, []string{"xml"}, productXmlIds)
		if err == nil && len(productsXmls) > 0 {
			relatedData.ProductsXmlsMap = BuildLookupProductsMap(productsXmls)
		}

		return err
	})

	g.Go(func() error {
		propertyTypesXmls, err := s.productPropertyTypesRepository.GetList(gctx, []string{"xml"})
		if err == nil && len(propertyTypesXmls) > 0 {
			relatedData.propertyTypesXmlsMap = utils.BuildLookupMap(propertyTypesXmls)
		}

		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return relatedData, nil
}

func prepareProps(data []dto_properties.PropertyDto) []entity.ProductPropertyTypes {
	propertyTypes := make([]entity.ProductPropertyTypes, 0, len(data))
	duplicateXmls := make(map[uuid.UUID]bool, len(data))

	for _, prop := range data {
		if duplicateXmls[prop.Xml] {
			continue
		}
		duplicateXmls[prop.Xml] = true

		propertyTypes = append(propertyTypes, entity.ProductPropertyTypes{
			Xml:      prop.Xml,
			IsActive: !prop.DeletionMark,
			Name:     prop.Name,
			Type:     prop.Type,
		})
	}

	return propertyTypes
}

func preparePropValues(data []dto_properties.PropertyValuesDto, productsXmlsMap map[string]uint, propertyTypesXmlsMap map[string]bool) []entity.ProductPropertyValues {
	propertyValues := make([]entity.ProductPropertyValues, 0, len(data))
	duplicateXmls := make(map[string]bool, len(data))

	for _, item := range data {
		_, isExist := productsXmlsMap[item.ProductXml.String()]
		if !isExist {
			continue
		}

		for _, prop := range item.Properties {
			isDuplicate := duplicateXmls[item.ProductXml.String()+prop.PropertyXml.String()]
			if isDuplicate {
				continue
			}

			duplicateXmls[item.ProductXml.String()+prop.PropertyXml.String()] = true

			if !propertyTypesXmlsMap[prop.PropertyXml.String()] {
				continue
			}

			propertyValues = append(propertyValues, entity.ProductPropertyValues{ProductXml: item.ProductXml, PropertyXml: prop.PropertyXml, Value: prop.Value})
		}
	}

	return propertyValues
}

// *** вспомогательные методы для обновления ***

func (s *ProductPropertyValuesService) getExistingPropValuesMap(ctx context.Context, propertyValueXmls *[]uuid.UUID) (*map[string]uint, error) {
	existingPropertyValues, err := s.productPropertyValues.GetListByXmls(ctx, []string{"id", "property_xml", "product_xml"}, *propertyValueXmls)
	if err != nil {
		return nil, err
	}

	existingPropertiesMap := make(map[string]uint, len(existingPropertyValues))
	for _, pv := range existingPropertyValues {
		existingPropertiesMap[pv.ProductXml.String()+pv.PropertyXml.String()] = pv.Id
	}

	return &existingPropertiesMap, nil
}

func buildPropetyValueXmls(data []dto_properties.PropertyValuesDto, productXmls *[]uuid.UUID, propertyValueXmls *[]uuid.UUID) []string {
	duplicateXmls := make(map[string]bool, len(data))
	var duplicated []string

	for _, item := range data {
		*productXmls = append(*productXmls, item.ProductXml)

		for _, prop := range item.Properties {
			isDuplicate := duplicateXmls[item.ProductXml.String()+prop.PropertyXml.String()]
			if isDuplicate {
				duplicated = append(duplicated, item.ProductXml.String()+"_"+prop.PropertyXml.String())
				continue
			}

			duplicateXmls[item.ProductXml.String()+prop.PropertyXml.String()] = true

			*propertyValueXmls = append(*propertyValueXmls, prop.PropertyXml)
		}
	}

	return duplicated
}

func buildPropertyValuesEntity(propertyValues []entity.ProductPropertyValues, existingPropertiesMap *map[string]uint) *dto_products.PreparedData[entity.ProductPropertyValues] {
	updLength := float64(len(*existingPropertiesMap)) * 0.66

	prepared := dto_products.PreparedData[entity.ProductPropertyValues]{
		CreateList:  make([]entity.ProductPropertyValues, 0, 200),
		UpdateList:  make([]entity.ProductPropertyValues, 0, int(updLength)),
		UpdatedXmls: make([]uuid.UUID, 0, int(updLength)),
	}

	for _, pv := range propertyValues {
		id, isExist := (*existingPropertiesMap)[pv.ProductXml.String()+pv.PropertyXml.String()]
		if isExist {
			prepared.UpdateList = append(prepared.UpdateList, entity.ProductPropertyValues{Id: id, PropertyXml: pv.PropertyXml, ProductXml: pv.ProductXml, Value: pv.Value})
		} else {
			prepared.CreateList = append(prepared.CreateList, entity.ProductPropertyValues{PropertyXml: pv.PropertyXml, ProductXml: pv.ProductXml, Value: pv.Value})
		}
	}

	return &prepared
}

func (s *ProductPropertyValuesService) updateOrCreatePropertyValues(
	ctx context.Context,
	prepared *dto_products.PreparedData[entity.ProductPropertyValues],
	result *dto_properties.UpdateBatchResponseDto,
) error {
	g, gctx := errgroup.WithContext(ctx)

	if len(prepared.CreateList) > 0 {
		g.Go(func() error {
			var err error

			result.Created, err = s.productPropertyValues.CreateBatch(gctx, prepared.CreateList)
			return err
		})
	}

	if len(prepared.UpdateList) > 0 {
		g.Go(func() error {
			var err error

			result.Updated, err = s.productPropertyValues.UpdateBatch(gctx, prepared.UpdateList)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *ProductPropertyValuesService) Clear() error {
	return s.productPropertyValues.Clear()
}
