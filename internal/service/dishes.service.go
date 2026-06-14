package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDishNotFound = errors.New("dish not found")

func (s *serv) CreateDish(ctx context.Context, name, description, menuType string) (*models.Dish, error) {
	d, err := s.repo.SaveDish(ctx, name, description, menuType)
	if err != nil {
		return nil, err
	}
	return &models.Dish{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		MenuType:    d.MenuType,
		Active:      d.Active,
		CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetDishes(ctx context.Context) ([]models.Dish, error) {
	entities, err := s.repo.GetDishes(ctx)
	if err != nil {
		return nil, err
	}

	dishes := make([]models.Dish, len(entities))
	for i, d := range entities {
		dishes[i] = models.Dish{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			MenuType:    d.MenuType,
			Active:      d.Active,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return dishes, nil
}

func (s *serv) GetDishesByMenuType(ctx context.Context, menuType string) ([]models.Dish, error) {
	entities, err := s.repo.GetDishesByMenuType(ctx, menuType)
	if err != nil {
		return nil, err
	}

	dishes := make([]models.Dish, len(entities))
	for i, d := range entities {
		dishes[i] = models.Dish{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			MenuType:    d.MenuType,
			Active:      d.Active,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return dishes, nil
}

func (s *serv) GetDishByID(ctx context.Context, id string) (*models.Dish, error) {
	d, err := s.repo.GetDishByID(ctx, id)
	if err != nil {
		return nil, ErrDishNotFound
	}
	return &models.Dish{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		MenuType:    d.MenuType,
		Active:      d.Active,
		CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateDish(ctx context.Context, id, name, description string, active bool) error {
	_, err := s.repo.GetDishByID(ctx, id)
	if err != nil {
		return ErrDishNotFound
	}
	return s.repo.UpdateDish(ctx, id, name, description, active)
}

func (s *serv) DeleteDish(ctx context.Context, id string) error {
	_, err := s.repo.GetDishByID(ctx, id)
	if err != nil {
		return ErrDishNotFound
	}
	return s.repo.DeleteDish(ctx, id)
}
