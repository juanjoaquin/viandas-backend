package service

import (
	"context"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrWeekMenuNotFound = errors.New("week menu not found")

func (s *serv) CreateWeekMenu(ctx context.Context, weekStartDate time.Time, createdBy string) (*models.WeekMenu, error) {
	m, err := s.repo.SaveWeekMenu(ctx, weekStartDate, createdBy)
	if err != nil {
		return nil, err
	}
	return &models.WeekMenu{
		ID:            m.ID,
		WeekStartDate: m.WeekStartDate.Format("2006-01-02"),
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetWeekMenus(ctx context.Context) ([]models.WeekMenu, error) {
	entities, err := s.repo.GetWeekMenus(ctx)
	if err != nil {
		return nil, err
	}

	menus := make([]models.WeekMenu, len(entities))
	for i, m := range entities {
		menus[i] = models.WeekMenu{
			ID:            m.ID,
			WeekStartDate: m.WeekStartDate.Format("2006-01-02"),
			CreatedBy:     m.CreatedBy,
			CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return menus, nil
}

func (s *serv) GetWeekMenuByID(ctx context.Context, id string) (*models.WeekMenu, error) {
	m, err := s.repo.GetWeekMenuByID(ctx, id)
	if err != nil {
		return nil, ErrWeekMenuNotFound
	}

	items, err := s.repo.GetWeekMenuItems(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	menuItems := make([]models.WeekMenuItem, len(items))
	for i, item := range items {
		tDish, _ := s.repo.GetDishByID(ctx, item.TraditionalDishID)
		hDish, _ := s.repo.GetDishByID(ctx, item.HealthyDishID)
		vDish, _ := s.repo.GetDishByID(ctx, item.VegetarianDishID)

		mi := models.WeekMenuItem{
			ID:         item.ID,
			WeekMenuID: item.WeekMenuID,
			MenuDate:   item.MenuDate.Format("2006-01-02"),
		}
		if tDish != nil {
			mi.TraditionalDish = &models.Dish{ID: tDish.ID, Name: tDish.Name, MenuType: tDish.MenuType}
		}
		if hDish != nil {
			mi.HealthyDish = &models.Dish{ID: hDish.ID, Name: hDish.Name, MenuType: hDish.MenuType}
		}
		if vDish != nil {
			mi.VegetarianDish = &models.Dish{ID: vDish.ID, Name: vDish.Name, MenuType: vDish.MenuType}
		}
		menuItems[i] = mi
	}

	return &models.WeekMenu{
		ID:            m.ID,
		WeekStartDate: m.WeekStartDate.Format("2006-01-02"),
		CreatedBy:     m.CreatedBy,
		Items:         menuItems,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) AddWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, traditionalDishID, healthyDishID, vegetarianDishID string) (*models.WeekMenuItem, error) {
	item, err := s.repo.SaveWeekMenuItem(ctx, weekMenuID, menuDate, traditionalDishID, healthyDishID, vegetarianDishID)
	if err != nil {
		return nil, err
	}
	return &models.WeekMenuItem{
		ID:         item.ID,
		WeekMenuID: item.WeekMenuID,
		MenuDate:   item.MenuDate.Format("2006-01-02"),
	}, nil
}

func (s *serv) UpdateWeekMenuItem(ctx context.Context, id, traditionalDishID, healthyDishID, vegetarianDishID string) error {
	return s.repo.UpdateWeekMenuItem(ctx, id, traditionalDishID, healthyDishID, vegetarianDishID)
}

func (s *serv) DeleteWeekMenuItem(ctx context.Context, id string) error {
	return s.repo.DeleteWeekMenuItem(ctx, id)
}
