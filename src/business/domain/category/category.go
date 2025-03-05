package category

import (
	"context"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(param entity.CategoryParam) ([]entity.Category, error)
	Get(ctx context.Context, param entity.CategoryParam) (entity.Category, error)
}

type category struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	c := &category{
		db: db,
	}

	return c
}

func (c *category) GetList(param entity.CategoryParam) ([]entity.Category, error) {
	res := []entity.Category{}

	if err := c.db.Where(param).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (c *category) Get(ctx context.Context, param entity.CategoryParam) (entity.Category, error) {
	res := entity.Category{}

	if err := c.db.WithContext(ctx).Where(param).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
