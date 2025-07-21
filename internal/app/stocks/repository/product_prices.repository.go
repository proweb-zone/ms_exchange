package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/stocks/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductPricesRepository struct {
	gormIns *gorm.DB
}

func InitProductPricesRepository(gormIns *gorm.DB) *ProductPricesRepository {
	return &ProductPricesRepository{gormIns: gormIns}
}

func (r *ProductPricesRepository) CreateBatch(ctx context.Context, productPrices []entity.ProductPrices) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.ProductPrices{}).CreateInBatches(&productPrices, 1000)
	if result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении складских остатков")
}

func (r *ProductPricesRepository) UpdateBatch(ctx context.Context, productPrices []entity.ProductPrices) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"price", "is_active"}),
	}).CreateInBatches(productPrices, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении складских остатков")
	}

	return result.RowsAffected, nil
}

func (r *ProductPricesRepository) GetList(ctx context.Context, fields []string) ([]entity.ProductPrices, error) {
	list := []entity.ProductPrices{}

	if result := r.gormIns.WithContext(ctx).Model(entity.ProductPrices{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *ProductPricesRepository) ExcDelete(ctx context.Context, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Where("id NOT IN ?", ids).
		Delete(entity.ProductPrices{}); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("запись не была удалена")
}

func (r *ProductPricesRepository) ChangeStatus(ctx context.Context, DeletionMark bool, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Model(entity.ProductPrices{}).
		Where("id NOT IN ?", ids).
		Update("is_active", DeletionMark); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("статус не обновлен")
}
