package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/brands/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BrandsRepository struct {
	gormIns *gorm.DB
}

func InitBrandsRepository(gormIns *gorm.DB) *BrandsRepository {
	return &BrandsRepository{gormIns: gormIns}
}

func (r *BrandsRepository) GetList(ctx context.Context, fields []string) ([]entity.Brands, error) {
	list := make([]entity.Brands, 0, 200)

	if result := r.gormIns.WithContext(ctx).Model(entity.Brands{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *BrandsRepository) CreateBatch(ctx context.Context, brands []entity.Brands) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(entity.Brands{}).CreateInBatches(&brands, 500)
	if result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении брендов")
}

func (r *BrandsRepository) UpdateBatch(ctx context.Context, brands []entity.Brands) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "xml"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "is_active"}),
	}).CreateInBatches(brands, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении брендов")
	}

	return result.RowsAffected, nil
}

func (r *BrandsRepository) ChangeStatus(ctx context.Context, DeletionMark bool, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("список идентификаторов для удаления пуст")
	}

	if result := r.gormIns.WithContext(ctx).
		Model(entity.Brands{}).
		Where("id NOT IN ?", ids).
		Update("is_active", DeletionMark); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("статус не обновлен")
}
