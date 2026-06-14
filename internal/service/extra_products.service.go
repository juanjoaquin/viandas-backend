package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrExtraProductNotFound = errors.New("extra product not found")

func (s *serv) CreateExtraProduct(ctx context.Context, name, category string) (*models.ExtraProduct, error) {
	e, err := s.repo.SaveExtraProduct(ctx, name, category)
	if err != nil {
		return nil, err
	}
	return &models.ExtraProduct{
		ID:        e.ID,
		Name:      e.Name,
		Category:  e.Category,
		Active:    e.Active,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetExtraProducts(ctx context.Context) ([]models.ExtraProduct, error) {
	entities, err := s.repo.GetExtraProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]models.ExtraProduct, len(entities))
	for i, e := range entities {
		products[i] = models.ExtraProduct{
			ID:        e.ID,
			Name:      e.Name,
			Category:  e.Category,
			Active:    e.Active,
			CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return products, nil
}

func (s *serv) GetExtraProductByID(ctx context.Context, id string) (*models.ExtraProduct, error) {
	e, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return nil, ErrExtraProductNotFound
	}
	return &models.ExtraProduct{
		ID:        e.ID,
		Name:      e.Name,
		Category:  e.Category,
		Active:    e.Active,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateExtraProduct(ctx context.Context, id, name string, active bool) error {
	_, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return ErrExtraProductNotFound
	}
	return s.repo.UpdateExtraProduct(ctx, id, name, active)
}

func (s *serv) DeleteExtraProduct(ctx context.Context, id string) error {
	_, err := s.repo.GetExtraProductByID(ctx, id)
	if err != nil {
		return ErrExtraProductNotFound
	}
	return s.repo.DeleteExtraProduct(ctx, id)
}
