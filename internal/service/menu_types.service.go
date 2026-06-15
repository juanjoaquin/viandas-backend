package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrMenuTypeNotFound = errors.New("menu type not found")

func (s *serv) CreateMenuType(ctx context.Context, name string, price *float64) (*models.MenuType, error) {
	mt, err := s.repo.SaveMenuType(ctx, name, price)
	if err != nil {
		return nil, err
	}
	return &models.MenuType{
		ID:        mt.ID,
		Name:      mt.Name,
		Price:     mt.Price,
		Active:    mt.Active,
		CreatedAt: mt.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetMenuTypes(ctx context.Context) ([]models.MenuType, error) {
	entities, err := s.repo.GetMenuTypes(ctx)
	if err != nil {
		return nil, err
	}

	types := make([]models.MenuType, len(entities))
	for i, mt := range entities {
		types[i] = models.MenuType{
			ID:        mt.ID,
			Name:      mt.Name,
			Price:     mt.Price,
			Active:    mt.Active,
			CreatedAt: mt.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return types, nil
}

func (s *serv) GetMenuTypeByID(ctx context.Context, id string) (*models.MenuType, error) {
	mt, err := s.repo.GetMenuTypeByID(ctx, id)
	if err != nil {
		return nil, ErrMenuTypeNotFound
	}
	return &models.MenuType{
		ID:        mt.ID,
		Name:      mt.Name,
		Price:     mt.Price,
		Active:    mt.Active,
		CreatedAt: mt.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateMenuType(ctx context.Context, id, name string, price *float64, active bool) error {
	_, err := s.repo.GetMenuTypeByID(ctx, id)
	if err != nil {
		return ErrMenuTypeNotFound
	}
	return s.repo.UpdateMenuType(ctx, id, name, price, active)
}

func (s *serv) DeleteMenuType(ctx context.Context, id string) error {
	_, err := s.repo.GetMenuTypeByID(ctx, id)
	if err != nil {
		return ErrMenuTypeNotFound
	}
	return s.repo.DeleteMenuType(ctx, id)
}
