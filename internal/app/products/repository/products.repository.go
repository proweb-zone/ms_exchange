package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/products/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductsRepository struct {
	gormIns *gorm.DB
}

func InitProductsRepository(gormIns *gorm.DB) *ProductsRepository {
	return &ProductsRepository{gormIns: gormIns}
}

func (r *ProductsRepository) GetListByXmlIds(ctx context.Context, fields []string, xmlIds []uuid.UUID) ([]entity.Products, error) {
	if len(xmlIds) == 0 {
		return []entity.Products{}, nil
	}

	list := make([]entity.Products, 0, len(xmlIds))

	result := r.gormIns.WithContext(ctx).Model(entity.Products{}).Select(fields).Where("xml IN ?", xmlIds).Find(&list)
	if result.Error == nil {
		return list, nil
	}

	return list, result.Error
}

func (r *ProductsRepository) GetList(ctx context.Context, fields []string) ([]string, error) {
	list := []string{}

	if result := r.gormIns.WithContext(ctx).Model(entity.Products{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *ProductsRepository) UpdateBatch(ctx context.Context, products []entity.Products) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "xml_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"article", "is_active", "name", "description", "code_toledo", "unit", "step"}),
	}).CreateInBatches(products, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении складских остатков")
	}

	return result.RowsAffected, nil
}

func (r *ProductsRepository) CreateBatch(ctx context.Context, products []entity.Products) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.Products{}).CreateInBatches(&products, 1000)
	if result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении товаров")
}

func (r *ProductsRepository) ChangeStatus(ctx context.Context, DeletionMark bool, updatedXmls []uuid.UUID) (int64, error) {
	if len(updatedXmls) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Model(entity.Products{}).
		Where("xml_id NOT IN ?", updatedXmls).
		Update("is_active", DeletionMark); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("статус не обновлен")
}

func (r *ProductsRepository) CountProducts(ctx context.Context) int64 {
	var count int64

	if err := r.gormIns.WithContext(ctx).Model(entity.Products{}).Select("id").Count(&count).Error; err == nil {
		return count
	}

	return 0
}
