package service

import (
	"context"
	"fmt"
	brandsEntity "ms_exchange/internal/app/brands/entity"
	brands "ms_exchange/internal/app/brands/repository"
	"ms_exchange/internal/app/products/dto/dto_products"
	"ms_exchange/internal/app/products/entity"

	products "ms_exchange/internal/app/products/repository"
	"ms_exchange/internal/app/stocks/dto/dto_stocks"
	"ms_exchange/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ProductRepeatRequest struct {
}

type RelatedData struct {
	CategoriesB2c []string
	CategoriesB2b []string
	Brands        []brandsEntity.Brands
}

type CategoryMaps struct {
	categoriesB2cXmlMap map[string]bool
	categoriesB2bXmlMap map[string]bool
	brandsXmlMap        map[string]bool
}

type ProductsService struct {
	productsRepository             *products.ProductsRepository
	productCategoriesB2cRepository *products.ProductCategoriesB2cRepository
	productCategoriesB2bRepository *products.ProductCategoriesB2bRepository
	brandsRepository               *brands.BrandsRepository
}

func NewProductsService(
	productsRepository *products.ProductsRepository,
	productCategoriesB2cRepository *products.ProductCategoriesB2cRepository,
	productCategoriesB2bRepository *products.ProductCategoriesB2bRepository,
	brandsRepository *brands.BrandsRepository,
) *ProductsService {
	return &ProductsService{
		productsRepository:             productsRepository,
		productCategoriesB2cRepository: productCategoriesB2cRepository,
		productCategoriesB2bRepository: productCategoriesB2bRepository,
		brandsRepository:               brandsRepository,
	}
}

func (s *ProductsService) CreateOrUpdateBatch(ctx context.Context, request dto_products.ProductsRequestDto) (*dto_stocks.UpdateBatchResponseDto, error) {
	updateBatchReport := &dto_stocks.UpdateBatchResponseDto{}

	// сопутствующая инфа
	relatedData, err := s.getRelatedData(ctx)
	if err != nil || relatedData == nil {
		return updateBatchReport, err
	}

	if len(relatedData.CategoriesB2c) == 0 && len(relatedData.CategoriesB2b) == 0 || len(relatedData.Brands) == 0 {
		return updateBatchReport, nil
	}

	сategoryMaps := CategoryMaps{}

	if len(relatedData.CategoriesB2c) > 0 {
		сategoryMaps.categoriesB2cXmlMap = utils.BuildLookupMap(relatedData.CategoriesB2c)
	}
	if len(relatedData.CategoriesB2b) > 0 {
		сategoryMaps.categoriesB2bXmlMap = utils.BuildLookupMap(relatedData.CategoriesB2b)
	}
	сategoryMaps.brandsXmlMap = buildLookupBrandsMap(relatedData.Brands)

	productsXmlMap := make(map[string]uint, 100)
	countProducts := s.productsRepository.CountProducts(ctx)

	if countProducts > 0 {
		xmlIds := buildXmIdsLookupMap(request.Data)
		productsXmls, _ := s.productsRepository.GetListByXmlIds(ctx, []string{"id", "xml"}, xmlIds)
		productsXmlMap = BuildLookupProductsMap(productsXmls)
	}

	productUnits := buildLookupUnitsMap(request.General.Units)

	// подготовка данных
	prepared := s.prepareAndCreateProducts(
		request.Data,
		productsXmlMap,
		productUnits,
		сategoryMaps,
	)

	// создание
	if countProducts == 0 {
		var err error
		updateBatchReport.Created, err = s.productsRepository.CreateBatch(ctx, prepared.CreateList)

		// *** отправка ***
		// url := "http://api-ka.toledo24.ru/ka_ivankov/hs/toledo_ecxhange/data/group"
		// data := map[string]interface{}{
		// 	"general": map[string]string{
		// 		"resource": "B2B",
		// 	},
		// 	"data": []map[string]string{},
		// }

		//  *** НЕДОСТАЮЩИЕ КАТЕГОРИИ ***
		if len(prepared.ForRepeatRequestBrandMap) > 0 {
			for xml, productItem := range prepared.ForRepeatRequestBrandMap {
				fmt.Print("\nБРЕНДЫ: ", xml, " ", productItem.Article)
			}
		}

		if len(prepared.ForRepeatRequestCatB2bMap) > 0 {
			// for xml := range prepared.ForRepeatRequestCatB2bMap {
			// item := map[string]string{"xml": xml.String()}
			// data["data"] = append(data["data"].([]map[string]string), item)
			// fmt.Print("\nForRepeatRequestCatB2bMap: ", xml, " ", productItem.Article)
			// }
		}

		if len(prepared.ForRepeatRequestCatB2cMap) > 0 {
			// for xml, _ := range prepared.ForRepeatRequestCatB2cMap {
			// item := map[string]string{"xml": xml.String()}
			// data["data"] = append(data["data"].([]map[string]string), item)
			// fmt.Print("\nForRepeatRequestCatB2cMap: ", xml, " ", productItem.Article)
			// }
		}

		// jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка при маршализации данных:", err)
			return updateBatchReport, err
		}

		// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		// if err != nil {
		// 	fmt.Println("Ошибка при создании запроса:", err)
		// 	return updateBatchReport, err
		// }

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println("Ошибка при отправке запроса:", err)
		// 	return updateBatchReport, err
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode != http.StatusOK {
		// 	fmt.Println("Ошибка: получен статус", resp.Status)
		// 	return updateBatchReport, fmt.Errorf("получен статус: %s", resp.Status)
		// }

		// answer, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println("Ошибка при чтении ответа:", err)
		// 	return updateBatchReport, err
		// }
		// fmt.Print("ответ:", string(answer))
		// *** конец отправка ***

		return updateBatchReport, err
	}

	// обновление
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		updateBatchReport.SoftDeleted, _ = s.productsRepository.ChangeStatus(gctx, false, prepared.UpdatedXmls)
		return nil
	})

	if len(prepared.UpdateList) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Updated, err = s.productsRepository.UpdateBatch(gctx, prepared.UpdateList)
			return err
		})
	}

	if len(prepared.CreateList) > 0 {
		g.Go(func() error {
			var err error
			updateBatchReport.Created, err = s.productsRepository.CreateBatch(gctx, prepared.CreateList)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return updateBatchReport, err
	}

	return updateBatchReport, nil
}

func (s *ProductsService) prepareAndCreateProducts(
	data []dto_products.ProductDto,
	productsXmlMap map[string]uint,
	productUnits map[string]string,
	categoryMaps CategoryMaps,
) *dto_products.PreparedData[entity.Products] {
	duplicateXmls := make(map[string]bool, len(data))
	prepared := dto_products.PreparedData[entity.Products]{
		CreateList:                make([]entity.Products, 0, 1000),
		UpdateList:                make([]entity.Products, 0, 1000),
		UpdatedXmls:               make([]uuid.UUID, 0, 100),
		ForRepeatRequestCatB2cMap: make(map[uuid.UUID]*entity.Products, 10),
		ForRepeatRequestCatB2bMap: make(map[uuid.UUID]*entity.Products, 10),
		ForRepeatRequestBrandMap:  make(map[uuid.UUID]*entity.Products, 10),
	}

	for _, item := range data {
		unit, isExist := productUnits[item.UnitXml]
		if !isExist {
			unit = ""
		}

		isNeedRepeatRequest := !categoryMaps.categoriesB2bXmlMap[item.Categories.CategoryXmlB2b.String()] ||
			!categoryMaps.categoriesB2cXmlMap[item.Categories.CategoryXmlB2c.String()] ||
			!categoryMaps.brandsXmlMap[item.BrandXml]

		if isNeedRepeatRequest {
			// сохранение товаров с отсутсвующими категориями для дозапроса
			fillRequests(&item, unit, &categoryMaps, &prepared)
			continue
		}

		duplicateXmls[item.ProductXml.String()] = true

		var categoryXmlB2c *uuid.UUID
		var categoryXmlB2b *uuid.UUID
		if categoryMaps.categoriesB2cXmlMap[item.Categories.CategoryXmlB2c.String()] {
			categoryXmlB2c = &item.Categories.CategoryXmlB2c
		}
		if categoryMaps.categoriesB2bXmlMap[item.Categories.CategoryXmlB2b.String()] {
			categoryXmlB2b = &item.Categories.CategoryXmlB2b
		}

		productId, isExist := productsXmlMap[item.ProductXml.String()]

		if !isExist {
			prepared.CreateList = append(prepared.CreateList, entity.Products{
				Xml:            item.ProductXml,
				CategoryXmlB2c: categoryXmlB2c,
				CategoryXmlB2b: categoryXmlB2b,
				IsActive:       !item.DeletionMark,
				Name:           item.Name,
				Article:        item.Article,
				CodeToledo:     item.CodeToledo,
				Description:    item.Description,
				BrandXml:       item.BrandXml,
				Unit:           unit,
				Step:           item.Step,
			})

			continue
		}

		prepared.UpdatedXmls = append(prepared.UpdatedXmls, item.ProductXml)
		prepared.UpdateList = append(prepared.UpdateList, entity.Products{
			Id:             productId,
			Xml:            item.ProductXml,
			CategoryXmlB2c: categoryXmlB2c,
			CategoryXmlB2b: categoryXmlB2b,
			IsActive:       !item.DeletionMark,
			Name:           item.Name,
			Article:        item.Article,
			CodeToledo:     item.CodeToledo,
			Description:    item.Description,
			BrandXml:       item.BrandXml,
			Unit:           unit,
			Step:           item.Step,
		})
	}

	return &prepared
}

func (s *ProductsService) getRelatedData(ctx context.Context) (*RelatedData, error) {
	relatedData := &RelatedData{}
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		relatedData.CategoriesB2c, err = s.productCategoriesB2cRepository.GetList(gctx, []string{"xml"})
		return err
	})

	g.Go(func() error {
		var err error
		relatedData.CategoriesB2b, err = s.productCategoriesB2bRepository.GetList(gctx, []string{"xml"})
		return err
	})

	g.Go(func() error {
		var err error
		relatedData.Brands, err = s.brandsRepository.GetList(gctx, []string{"xml"})
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return relatedData, nil
}

func buildLookupBrandsMap(brands []brandsEntity.Brands) map[string]bool {
	result := make(map[string]bool, len(brands))

	for _, item := range brands {
		result[item.Xml.String()] = true
	}

	return result
}

func buildLookupUnitsMap(units []dto_products.UnitsDto) map[string]string {
	result := make(map[string]string, len(units))

	for _, item := range units {
		result[item.Xml] = item.Name
	}

	return result
}

func buildXmIdsLookupMap(data []dto_products.ProductDto) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(data))
	duplicateXmls := make(map[string]bool, len(data))

	for _, item := range data {
		if duplicateXmls[item.ProductXml.String()] {
			continue
		}
		duplicateXmls[item.ProductXml.String()] = true

		result = append(result, item.ProductXml)
	}

	return result
}

func BuildLookupProductsMap(slice []entity.Products) map[string]uint {
	mappedSlice := make(map[string]uint, len(slice))

	for _, key := range slice {
		mappedSlice[key.Xml.String()] = key.Id
	}

	return mappedSlice
}

// нужно сохранить товары без категорий, брендов... для дозапроса из 1С
func fillRequests(
	currentItem *dto_products.ProductDto,
	currentUnit string,
	categoryMaps *CategoryMaps,
	prepared *dto_products.PreparedData[entity.Products],
) {
	product := &entity.Products{
		Xml:            currentItem.ProductXml,
		CategoryXmlB2c: nil,
		CategoryXmlB2b: nil,
		IsActive:       currentItem.DeletionMark,
		Name:           currentItem.Name,
		Article:        currentItem.Article,
		CodeToledo:     currentItem.CodeToledo,
		Description:    currentItem.Description,
		BrandXml:       currentItem.BrandXml,
		Unit:           currentUnit,
		Step:           currentItem.Step,
	}

	if !categoryMaps.categoriesB2bXmlMap[currentItem.Categories.CategoryXmlB2b.String()] {
		prepared.ForRepeatRequestCatB2bMap[currentItem.Categories.CategoryXmlB2b] = product
	}
	if !categoryMaps.categoriesB2cXmlMap[currentItem.Categories.CategoryXmlB2c.String()] {
		prepared.ForRepeatRequestCatB2cMap[currentItem.Categories.CategoryXmlB2c] = product
	}
	if !categoryMaps.brandsXmlMap[currentItem.BrandXml] {
		xml, err := uuid.Parse(currentItem.BrandXml)
		if err == nil {
			prepared.ForRepeatRequestBrandMap[xml] = product
		}
	}
}
