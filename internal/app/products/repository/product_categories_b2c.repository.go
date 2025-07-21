package repository

import (
	"context"
	"errors"
	"ms_exchange/internal/app/products/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductCategoriesB2cRepository struct {
	gormIns *gorm.DB
}

func InitProductCategoriesB2cRepository(gormIns *gorm.DB) *ProductCategoriesB2cRepository {
	return &ProductCategoriesB2cRepository{gormIns: gormIns}
}

func (r *ProductCategoriesB2cRepository) CreateBatch(ctx context.Context, categories *entity.ProductCategoriesB2c) (int64, error) {
	if result := r.gormIns.WithContext(ctx).Model(entity.ProductCategoriesB2c{}).CreateInBatches(categories, 1000); result.Error == nil {
		return result.RowsAffected, nil
	}

	return 0, errors.New("ошибка при сохранении данных")
}

func (r *ProductCategoriesB2cRepository) GetListAll(ctx context.Context, fields []string) ([]entity.ProductCategoriesB2c, error) {
	list := []entity.ProductCategoriesB2c{}

	if result := r.gormIns.WithContext(ctx).Model(entity.ProductCategoriesB2c{}).Preload("Parent").Preload("Children").Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *ProductCategoriesB2cRepository) GetList(ctx context.Context, fields []string) ([]string, error) {
	list := []string{}

	if result := r.gormIns.WithContext(ctx).Model(entity.ProductCategoriesB2c{}).Select(fields).Find(&list); result.Error == nil {
		return list, nil
	}

	return list, nil
}

func (r *ProductCategoriesB2cRepository) ClearCategories() error {
	err := r.gormIns.Exec("TRUNCATE TABLE product_categories_b2c CASCADE").Error
	if err != nil {
		return err
	}

	r.gormIns.Exec("ALTER SEQUENCE product_categories_b2c_id_seq RESTART WITH 1")
	return err
}

func (r *ProductCategoriesB2cRepository) Create(ctx context.Context, category *entity.ProductCategoriesB2c) (*uuid.UUID, error) {
	var err error = nil
	if result := r.gormIns.WithContext(ctx).Model(entity.ProductCategoriesB2c{}).Create(category); result.Error != nil {
		err = errors.New("запись не была создана")
	}

	return &category.Xml, err
}

func (r *ProductCategoriesB2cRepository) UpdateBatch(ctx context.Context, categories []entity.ProductCategoriesB2c) (int64, error) {
	result := r.gormIns.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "xml_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "is_active", "parent_id"}),
	}).CreateInBatches(categories, 500)

	if result.Error != nil {
		return result.RowsAffected, errors.New("ошибка при обновлении категорий")
	}

	return result.RowsAffected, nil
}
