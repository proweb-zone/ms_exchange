package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/stocks/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PriceTypesRepository struct {
	gormIns *gorm.DB
}

func InitPriceTypesRepository(gormIns *gorm.DB) *PriceTypesRepository {
	return &PriceTypesRepository{gormIns: gormIns}
}

func (r *PriceTypesRepository) CreateBatch(ctx context.Context, priceTypes []entity.PriceTypes) (int64, error) {
	if result := r.gormIns.WithContext(ctx).Model(entity.PriceTypes{}).Create(priceTypes); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении данных")
}

func (r *PriceTypesRepository) GetList(ctx context.Context, fields []string) ([]entity.PriceTypes, error) {
	list := make([]entity.PriceTypes, 0, 5)

	if result := r.gormIns.WithContext(ctx).Model(entity.PriceTypes{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *PriceTypesRepository) ClearPriceTypes() error {
	return r.gormIns.Delete(entity.PriceTypes{}, "1 = 1").Error
}

func (r *PriceTypesRepository) Update(ctx context.Context, xml uuid.UUID, priceType entity.PriceTypes) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.PriceTypes{}).Updates(priceType).Where("xml = ?", xml)

	return result.RowsAffected, nil
}

func (r *PriceTypesRepository) UpdateBatch(ctx context.Context, storages []entity.PriceTypes) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "xml"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "is_active"}),
	}).CreateInBatches(storages, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении цен")
	}

	return result.RowsAffected, nil
}
