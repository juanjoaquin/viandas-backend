package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrProductCategoryNotFound = errors.New("product category not found")

func mapProductCategory(category *entity.ProductCategory) *models.ProductCategory {
	if category == nil {
		return nil
	}
	return &models.ProductCategory{
		ID:        category.ID,
		Name:      category.Name,
		Active:    category.Active,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *serv) CreateProductCategory(ctx context.Context, name string) (*models.ProductCategory, error) {
	category, err := s.repo.SaveProductCategory(ctx, name)
	if err != nil {
		return nil, err
	}
	return mapProductCategory(category), nil
}

func (s *serv) GetProductCategories(ctx context.Context, nameQuery string, activeFilter *bool) ([]models.ProductCategory, error) {
	entities, err := s.repo.GetProductCategories(ctx, nameQuery, activeFilter)
	if err != nil {
		return nil, err
	}

	categories := make([]models.ProductCategory, len(entities))
	for i, category := range entities {
		categories[i] = *mapProductCategory(&category)
	}
	return categories, nil
}

func (s *serv) GetProductCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error) {
	category, err := s.repo.GetProductCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrProductCategoryNotFound
	}
	return mapProductCategory(category), nil
}

func (s *serv) UpdateProductCategory(ctx context.Context, id, name string, active bool) error {
	_, err := s.repo.GetProductCategoryByID(ctx, id)
	if err != nil {
		return ErrProductCategoryNotFound
	}
	return s.repo.UpdateProductCategory(ctx, id, name, active)
}

func (s *serv) DeleteProductCategory(ctx context.Context, id string) error {
	_, err := s.repo.GetProductCategoryByID(ctx, id)
	if err != nil {
		return ErrProductCategoryNotFound
	}
	return s.repo.DeleteProductCategory(ctx, id)
}
