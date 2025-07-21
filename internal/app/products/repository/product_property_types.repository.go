package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/products/entity"

	"gorm.io/gorm"
)

type ProductPropertyTypesRepository struct {
	gormIns *gorm.DB
}

func InitProductPropertyTypesRepository(gormIns *gorm.DB) *ProductPropertyTypesRepository {
	return &ProductPropertyTypesRepository{gormIns: gormIns}
}

func (r *ProductPropertyTypesRepository) CreateBatch(ctx context.Context, propertyTypes []entity.ProductPropertyTypes) (int64, error) {
	if result := r.gormIns.WithContext(ctx).Model(entity.ProductPropertyTypes{}).CreateInBatches(propertyTypes, 1000); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении данных")
}

func (r *ProductPropertyTypesRepository) GetList(ctx context.Context, fields []string) ([]string, error) {
	list := []string{}

	if result := r.gormIns.WithContext(ctx).Model(entity.ProductPropertyTypes{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *ProductPropertyTypesRepository) ClearProperties() error {
	return r.gormIns.Delete(entity.ProductPropertyTypes{}, "1 = 1").Error
}
