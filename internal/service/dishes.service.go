package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDishNotFound = errors.New("dish not found")

func (s *serv) dishToModel(ctx context.Context, id, name, description, menuTypeID string, active bool, createdAt string) *models.Dish {
	var mt *models.MenuType
	if menuTypeID != "" {
		if entity, err := s.repo.GetMenuTypeByID(ctx, menuTypeID); err == nil {
			mt = &models.MenuType{
				ID:     entity.ID,
				Name:   entity.Name,
				Price:  entity.Price,
				Active: entity.Active,
			}
		}
	}
	return &models.Dish{
		ID:          id,
		Name:        name,
		Description: description,
		MenuType:    mt,
		Active:      active,
		CreatedAt:   createdAt,
	}
}

func (s *serv) CreateDish(ctx context.Context, name, description, menuTypeID string) (*models.Dish, error) {
	d, err := s.repo.SaveDish(ctx, name, description, menuTypeID)
	if err != nil {
		return nil, err
	}
	return s.dishToModel(ctx, d.ID, d.Name, d.Description, d.MenuTypeID, d.Active, d.CreatedAt.Format("2006-01-02T15:04:05Z")), nil
}

func (s *serv) GetDishes(ctx context.Context) ([]models.Dish, error) {
	entities, err := s.repo.GetDishes(ctx)
	if err != nil {
		return nil, err
	}

	dishes := make([]models.Dish, len(entities))
	for i, d := range entities {
		dishes[i] = *s.dishToModel(ctx, d.ID, d.Name, d.Description, d.MenuTypeID, d.Active, d.CreatedAt.Format("2006-01-02T15:04:05Z"))
	}
	return dishes, nil
}

func (s *serv) GetDishesByMenuTypeID(ctx context.Context, menuTypeID string) ([]models.Dish, error) {
	entities, err := s.repo.GetDishesByMenuTypeID(ctx, menuTypeID)
	if err != nil {
		return nil, err
	}

	dishes := make([]models.Dish, len(entities))
	for i, d := range entities {
		dishes[i] = *s.dishToModel(ctx, d.ID, d.Name, d.Description, d.MenuTypeID, d.Active, d.CreatedAt.Format("2006-01-02T15:04:05Z"))
	}
	return dishes, nil
}

func (s *serv) GetDishByID(ctx context.Context, id string) (*models.Dish, error) {
	d, err := s.repo.GetDishByID(ctx, id)
	if err != nil {
		return nil, ErrDishNotFound
	}
	return s.dishToModel(ctx, d.ID, d.Name, d.Description, d.MenuTypeID, d.Active, d.CreatedAt.Format("2006-01-02T15:04:05Z")), nil
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
