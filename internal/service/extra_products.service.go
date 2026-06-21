package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrExtraProductNotFound = errors.New("extra product not found")

func extraLineTotalAmount(price float64, quantity int) *float64 {
	amount := price * float64(quantity)
	return &amount
}

func (s *serv) mapExtraProduct(ctx context.Context, e *entity.ExtraProduct) *models.ExtraProduct {
	if e == nil {
		return nil
	}

	category, _ := s.repo.GetProductCategoryByID(ctx, e.CategoryID)

	return &models.ExtraProduct{
		ID:        e.ID,
		Name:      e.Name,
		Category:  mapProductCategory(category),
		Price:     e.Price,
		Active:    e.Active,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *serv) CreateExtraProduct(ctx context.Context, name, categoryID string, price float64) (*models.ExtraProduct, error) {
	if _, err := s.repo.GetProductCategoryByID(ctx, categoryID); err != nil {
		return nil, ErrProductCategoryNotFound
	}

	e, err := s.repo.SaveExtraProduct(ctx, name, categoryID, price)
	if err != nil {
		return nil, err
	}
	return s.mapExtraProduct(ctx, e), nil
}

func (s *serv) CountExtraProducts(ctx context.Context, nameQuery string) (int, error) {
	return s.repo.CountExtraProducts(ctx, nameQuery)
}

func (s *serv) GetExtraProducts(ctx context.Context, nameQuery string, offset, limit int) ([]models.ExtraProduct, error) {
	entities, err := s.repo.GetExtraProducts(ctx, nameQuery, offset, limit)
	if err != nil {
		return nil, err
	}

	products := make([]models.ExtraProduct, len(entities))
	for i, e := range entities {
		products[i] = *s.mapExtraProduct(ctx, &e)
	}
	return products, nil
}

func (s *serv) GetExtraProductByID(ctx context.Context, id string) (*models.ExtraProduct, error) {
	e, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return nil, ErrExtraProductNotFound
	}
	return s.mapExtraProduct(ctx, e), nil
}

func (s *serv) UpdateExtraProduct(ctx context.Context, id, name, categoryID string, price float64, active bool) error {
	_, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return ErrExtraProductNotFound
	}
	if _, err := s.repo.GetProductCategoryByID(ctx, categoryID); err != nil {
		return ErrProductCategoryNotFound
	}
	return s.repo.UpdateExtraProduct(ctx, id, name, categoryID, price, active)
}

func (s *serv) DeleteExtraProduct(ctx context.Context, id string) error {
	_, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return ErrExtraProductNotFound
	}
	return s.repo.DeleteExtraProduct(ctx, id)
}
