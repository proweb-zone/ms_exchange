package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/stocks/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockProductsRepository struct {
	gormIns *gorm.DB
}

func InitStockProductsRepository(gormIns *gorm.DB) *StockProductsRepository {
	return &StockProductsRepository{gormIns: gormIns}
}

func (r *StockProductsRepository) CreateBatch(ctx context.Context, productStorages []entity.ProductStorages) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.ProductStorages{}).CreateInBatches(&productStorages, 1000)
	if result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении складских остатков")
}

func (r *StockProductsRepository) UpdateBatch(ctx context.Context, productStorages []entity.ProductStorages) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "is_active"}),
	}).CreateInBatches(productStorages, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении складских остатков")
	}

	return result.RowsAffected, nil
}

func (r *StockProductsRepository) GetList(ctx context.Context, fields []string) ([]entity.ProductStorages, error) {
	list := []entity.ProductStorages{}

	if result := r.gormIns.WithContext(ctx).Model(entity.ProductStorages{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *StockProductsRepository) ExcDelete(ctx context.Context, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Where("id NOT IN ?", ids).
		Delete(&entity.ProductStorages{}); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("запись не была удалена")
}

func (r *StockProductsRepository) ChangeStatus(ctx context.Context, DeletionMark bool, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Model(entity.ProductStorages{}).
		Where("id NOT IN ?", ids).
		Update("is_active", DeletionMark); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("статус не обновлен")
}
