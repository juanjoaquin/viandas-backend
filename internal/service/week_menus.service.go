package service

import (
	"context"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrWeekMenuNotFound = errors.New("week menu not found")

func (s *serv) CreateWeekMenu(ctx context.Context, weekStartDate, weekEndDate time.Time, createdBy string) (*models.WeekMenu, error) {
	m, err := s.repo.SaveWeekMenu(ctx, weekStartDate, weekEndDate, createdBy)
	if err != nil {
		return nil, err
	}
	return &models.WeekMenu{
		ID:            m.ID,
		WeekStartDate: m.WeekStartDate.Format("2006-01-02"),
		WeekEndDate:   m.WeekEndDate.Format("2006-01-02"),
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
			WeekEndDate:   m.WeekEndDate.Format("2006-01-02"),
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
		mi := models.WeekMenuItem{
			ID:         item.ID,
			WeekMenuID: item.WeekMenuID,
			MenuDate:   item.MenuDate.Format("2006-01-02"),
		}

		if mtEntity, err := s.repo.GetMenuTypeByID(ctx, item.MenuTypeID); err == nil {
			mi.MenuType = &models.MenuType{
				ID:     mtEntity.ID,
				Name:   mtEntity.Name,
				Price:  mtEntity.Price,
				Active: mtEntity.Active,
			}
		}

		if dish, err := s.repo.GetDishByID(ctx, item.DishID); err == nil {
			mi.Dish = &models.Dish{ID: dish.ID, Name: dish.Name}
		}

		menuItems[i] = mi
	}

	return &models.WeekMenu{
		ID:            m.ID,
		WeekStartDate: m.WeekStartDate.Format("2006-01-02"),
		WeekEndDate:   m.WeekEndDate.Format("2006-01-02"),
		CreatedBy:     m.CreatedBy,
		Items:         menuItems,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) ResolveWeekMenu(ctx context.Context, requestedID string, date time.Time) (*models.WeekMenu, error) {
	if requestedID != "" {
		if menu, err := s.GetWeekMenuByID(ctx, requestedID); err == nil {
			return menu, nil
		}
	}

	menus, err := s.repo.GetWeekMenus(ctx)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return nil, ErrWeekMenuNotFound
	}

	for _, menu := range menus {
		if !menu.WeekStartDate.After(date) && !menu.WeekEndDate.Before(date) {
			return s.GetWeekMenuByID(ctx, menu.ID)
		}
	}

	var upcomingID string
	var upcomingStart time.Time
	for _, menu := range menus {
		if menu.WeekStartDate.After(date) && (upcomingID == "" || menu.WeekStartDate.Before(upcomingStart)) {
			upcomingID = menu.ID
			upcomingStart = menu.WeekStartDate
		}
	}
	if upcomingID != "" {
		return s.GetWeekMenuByID(ctx, upcomingID)
	}

	return s.GetWeekMenuByID(ctx, menus[0].ID)
}

func (s *serv) DeleteWeekMenu(ctx context.Context, id string) error {
	if _, err := s.repo.GetWeekMenuByID(ctx, id); err != nil {
		return ErrWeekMenuNotFound
	}

	return s.repo.DeleteWeekMenu(ctx, id)
}

func (s *serv) AddWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, menuTypeID, dishID string) (*models.WeekMenuItem, error) {
	item, err := s.repo.SaveWeekMenuItem(ctx, weekMenuID, menuDate, menuTypeID, dishID)
	if err != nil {
		return nil, err
	}

	mi := &models.WeekMenuItem{
		ID:         item.ID,
		WeekMenuID: item.WeekMenuID,
		MenuDate:   item.MenuDate.Format("2006-01-02"),
	}

	if mtEntity, err := s.repo.GetMenuTypeByID(ctx, item.MenuTypeID); err == nil {
		mi.MenuType = &models.MenuType{ID: mtEntity.ID, Name: mtEntity.Name}
	}

	if dish, err := s.repo.GetDishByID(ctx, item.DishID); err == nil {
		mi.Dish = &models.Dish{ID: dish.ID, Name: dish.Name}
	}

	return mi, nil
}

func (s *serv) UpdateWeekMenuItem(ctx context.Context, id, dishID string) error {
	return s.repo.UpdateWeekMenuItem(ctx, id, dishID)
}

func (s *serv) DeleteWeekMenuItem(ctx context.Context, id string) error {
	return s.repo.DeleteWeekMenuItem(ctx, id)
}
