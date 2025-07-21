package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/products/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductPropertyValuesRepository struct {
	gormIns *gorm.DB
}

func InitProductPropertyValuesRepository(gormIns *gorm.DB) *ProductPropertyValuesRepository {
	return &ProductPropertyValuesRepository{gormIns: gormIns}
}

func (r *ProductPropertyValuesRepository) CreateBatch(ctx context.Context, propertyValues []entity.ProductPropertyValues) (int64, error) {
	if result := r.gormIns.WithContext(ctx).Model(entity.ProductPropertyValues{}).CreateInBatches(propertyValues, 1000); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении данных")
}

func (r *ProductPropertyValuesRepository) GetListByXmls(ctx context.Context, fields []string, xmls []uuid.UUID) ([]entity.ProductPropertyValues, error) {
	list := make([]entity.ProductPropertyValues, 0, 500)

	result := r.gormIns.WithContext(ctx).Model(entity.ProductPropertyValues{}).Select(fields).Where("property_xml IN ?", xmls).Find(&list)
	if result.Error == nil {
		return list, nil
	}

	return list, result.Error
}

func (r *ProductPropertyValuesRepository) UpdateBatch(ctx context.Context, products []entity.ProductPropertyValues) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).CreateInBatches(products, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении складских остатков")
	}

	return result.RowsAffected, nil
}

func (r *ProductPropertyValuesRepository) Clear() error {
	return r.gormIns.Delete(entity.ProductPropertyValues{}, "1 = 1").Error
}
