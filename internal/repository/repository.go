package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository is a generic base for all GORM repositories.
type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := r.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Save(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Delete(entity).Error
}
