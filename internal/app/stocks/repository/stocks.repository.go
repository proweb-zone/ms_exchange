package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/stocks/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StocksRepository struct {
	gormIns *gorm.DB
}

func InitStocksRepository(gormIns *gorm.DB) *StocksRepository {
	return &StocksRepository{gormIns: gormIns}
}

func (r *StocksRepository) CreateBatch(ctx context.Context, storages []entity.Storages) (int64, error) {
	if result := r.gormIns.WithContext(ctx).Model(entity.Storages{}).Create(storages); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении данных")
}

func (r *StocksRepository) GetList(ctx context.Context, fields []string) ([]entity.Storages, error) {
	list := make([]entity.Storages, 0, 5)

	if result := r.gormIns.WithContext(ctx).Model(entity.Storages{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *StocksRepository) ClearStorages() error {
	return r.gormIns.Delete(entity.Storages{}, "1 = 1").Error
}

func (r *StocksRepository) Update(ctx context.Context, xml uuid.UUID, productStorages entity.ProductStorages) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.ProductStorages{}).Updates(productStorages).Where("xml = ?", xml)

	return result.RowsAffected, nil
}

func (r *StocksRepository) ChangeStatus(ctx context.Context, DeletionMark bool, xmls []uuid.UUID) (int64, error) {
	if len(xmls) == 0 {
		return 0, nil
	}

	if result := r.gormIns.WithContext(ctx).
		Model(entity.Storages{}).
		Where("xml NOT IN ?", xmls).
		Update("deletion_mark", DeletionMark); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("статус не обновлен")
}

func (r *StocksRepository) UpdateBatch(ctx context.Context, storages []entity.Storages) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "xml"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "is_active"}),
	}).CreateInBatches(storages, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении складских остатков")
	}

	return result.RowsAffected, nil
}
