package service

import (
	"context"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDailyProductionNotFound = errors.New("daily production not found")

func (s *serv) CreateDailyProduction(ctx context.Context, productionDate time.Time, customerID, deliveryID, notes, createdBy string) (*models.DailyProduction, error) {
	dp, err := s.repo.SaveDailyProduction(ctx, productionDate, customerID, deliveryID, notes, createdBy)
	if err != nil {
		return nil, err
	}
	return &models.DailyProduction{
		ID:             dp.ID,
		ProductionDate: dp.ProductionDate.Format("2006-01-02"),
		Notes:          dp.Notes,
		CreatedBy:      dp.CreatedBy,
		CreatedAt:      dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetDailyProductions(ctx context.Context, date time.Time) ([]models.DailyProduction, error) {
	entities, err := s.repo.GetDailyProductions(ctx, date)
	if err != nil {
		return nil, err
	}

	productions := make([]models.DailyProduction, len(entities))
	for i, dp := range entities {
		customer, _ := s.repo.GetCustomerByID(ctx, dp.CustomerID)
		var customerModel *models.Customer
		if customer != nil {
			customerModel = &models.Customer{ID: customer.ID, Name: customer.Name, Type: customer.Type}
		}

		var deliveryModel *models.Delivery
		if dp.DeliveryID != "" {
			delivery, _ := s.repo.GetDeliveryByID(ctx, dp.DeliveryID)
			if delivery != nil {
				deliveryModel = &models.Delivery{ID: delivery.ID, Name: delivery.Name, Phone: delivery.Phone}
			}
		}

		lines, _ := s.repo.GetDailyProductionLines(ctx, dp.ID)
		lineModels := make([]models.DailyProductionLine, len(lines))
		for j, line := range lines {
			lm := models.DailyProductionLine{
				ID:       line.ID,
				Quantity: line.Quantity,
			}
			if mt, err := s.repo.GetMenuTypeByID(ctx, line.MenuTypeID); err == nil {
				lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, SortOrder: mt.SortOrder}
			}
			lineModels[j] = lm
		}

		extras, _ := s.repo.GetDailyProductionExtras(ctx, dp.ID)
		extraModels := make([]models.DailyProductionExtra, len(extras))
		for j, ex := range extras {
			ep, _ := s.repo.GetExtraProductByID(ctx, ex.ExtraProductID)
			var epModel *models.ExtraProduct
			if ep != nil {
				epModel = &models.ExtraProduct{ID: ep.ID, Name: ep.Name, Category: ep.Category}
			}
			extraModels[j] = models.DailyProductionExtra{
				ID:           ex.ID,
				ExtraProduct: epModel,
				Quantity:     ex.Quantity,
			}
		}

		productions[i] = models.DailyProduction{
			ID:             dp.ID,
			ProductionDate: dp.ProductionDate.Format("2006-01-02"),
			Customer:       customerModel,
			Delivery:       deliveryModel,
			Lines:          lineModels,
			Notes:          dp.Notes,
			Extras:         extraModels,
			CreatedBy:      dp.CreatedBy,
			CreatedAt:      dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return productions, nil
}

func (s *serv) GetDailyProductionByID(ctx context.Context, id string) (*models.DailyProduction, error) {
	dp, err := s.repo.GetDailyProductionByID(ctx, id)
	if err != nil {
		return nil, ErrDailyProductionNotFound
	}

	lines, _ := s.repo.GetDailyProductionLines(ctx, dp.ID)
	lineModels := make([]models.DailyProductionLine, len(lines))
	for j, line := range lines {
		lm := models.DailyProductionLine{
			ID:       line.ID,
			Quantity: line.Quantity,
		}
		if mt, err := s.repo.GetMenuTypeByID(ctx, line.MenuTypeID); err == nil {
			lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, SortOrder: mt.SortOrder}
		}
		lineModels[j] = lm
	}

	return &models.DailyProduction{
		ID:             dp.ID,
		ProductionDate: dp.ProductionDate.Format("2006-01-02"),
		Lines:          lineModels,
		Notes:          dp.Notes,
		CreatedBy:      dp.CreatedBy,
		CreatedAt:      dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateDailyProduction(ctx context.Context, id, deliveryID, notes string) error {
	_, err := s.repo.GetDailyProductionByID(ctx, id)
	if err != nil {
		return ErrDailyProductionNotFound
	}
	return s.repo.UpdateDailyProduction(ctx, id, deliveryID, notes)
}

func (s *serv) UpsertDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*models.DailyProductionLine, error) {
	line, err := s.repo.UpsertDailyProductionLine(ctx, dailyProductionID, menuTypeID, quantity)
	if err != nil {
		return nil, err
	}

	lm := &models.DailyProductionLine{
		ID:       line.ID,
		Quantity: line.Quantity,
	}
	if mt, err := s.repo.GetMenuTypeByID(ctx, line.MenuTypeID); err == nil {
		lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, SortOrder: mt.SortOrder}
	}
	return lm, nil
}

func (s *serv) DeleteDailyProductionLine(ctx context.Context, id string) error {
	return s.repo.DeleteDailyProductionLine(ctx, id)
}

func (s *serv) AddDailyProductionExtra(ctx context.Context, dailyProductionID, extraProductID string, quantity int) (*models.DailyProductionExtra, error) {
	extra, err := s.repo.SaveDailyProductionExtra(ctx, dailyProductionID, extraProductID, quantity)
	if err != nil {
		return nil, err
	}

	ep, _ := s.repo.GetExtraProductByID(ctx, extra.ExtraProductID)
	var epModel *models.ExtraProduct
	if ep != nil {
		epModel = &models.ExtraProduct{ID: ep.ID, Name: ep.Name, Category: ep.Category}
	}

	return &models.DailyProductionExtra{
		ID:           extra.ID,
		ExtraProduct: epModel,
		Quantity:     extra.Quantity,
	}, nil
}

func (s *serv) DeleteDailyProductionExtra(ctx context.Context, id string) error {
	return s.repo.DeleteDailyProductionExtra(ctx, id)
}

func (s *serv) GetMenuByDate(ctx context.Context, date time.Time) (*models.DayMenu, error) {
	items, err := s.repo.GetWeekMenuItemsByDate(ctx, date)
	if err != nil || len(items) == 0 {
		return nil, errors.New("no menu configured for this date")
	}

	dayMenu := &models.DayMenu{
		Date:  date.Format("2006-01-02"),
		Lines: make([]models.DayMenuLine, 0, len(items)),
	}

	for _, item := range items {
		line := models.DayMenuLine{}

		if mt, err := s.repo.GetMenuTypeByID(ctx, item.MenuTypeID); err == nil {
			line.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, SortOrder: mt.SortOrder}
		}

		if dish, err := s.repo.GetDishByID(ctx, item.DishID); err == nil {
			var dishMenuType *models.MenuType
			if mt, err := s.repo.GetMenuTypeByID(ctx, dish.MenuTypeID); err == nil {
				dishMenuType = &models.MenuType{ID: mt.ID, Name: mt.Name}
			}
			line.Dish = &models.Dish{ID: dish.ID, Name: dish.Name, MenuType: dishMenuType}
		}

		dayMenu.Lines = append(dayMenu.Lines, line)
	}

	return dayMenu, nil
}

func (s *serv) GetKitchenTotals(ctx context.Context, date time.Time) (*models.KitchenTotals, error) {
	productions, err := s.repo.GetDailyProductions(ctx, date)
	if err != nil {
		return nil, err
	}

	totalsMap := make(map[string]*models.MenuTypeTotalQty)
	for _, dp := range productions {
		lines, err := s.repo.GetDailyProductionLines(ctx, dp.ID)
		if err != nil {
			continue
		}
		for _, line := range lines {
			if _, exists := totalsMap[line.MenuTypeID]; !exists {
				mt, _ := s.repo.GetMenuTypeByID(ctx, line.MenuTypeID)
				var mtModel *models.MenuType
				if mt != nil {
					mtModel = &models.MenuType{ID: mt.ID, Name: mt.Name, SortOrder: mt.SortOrder}
				}
				totalsMap[line.MenuTypeID] = &models.MenuTypeTotalQty{MenuType: mtModel, TotalQty: 0}
			}
			totalsMap[line.MenuTypeID].TotalQty += line.Quantity
		}
	}

	totals := &models.KitchenTotals{
		Date:   date.Format("2006-01-02"),
		Totals: make([]models.MenuTypeTotalQty, 0, len(totalsMap)),
	}
	for _, t := range totalsMap {
		totals.Totals = append(totals.Totals, *t)
	}
	return totals, nil
}

func (s *serv) GetExtrasTotals(ctx context.Context, date time.Time) (*models.ExtraTotals, error) {
	productions, err := s.repo.GetDailyProductions(ctx, date)
	if err != nil {
		return nil, err
	}

	totalsMap := make(map[string]*models.ExtraTotal)
	for _, dp := range productions {
		extras, err := s.repo.GetDailyProductionExtras(ctx, dp.ID)
		if err != nil {
			continue
		}
		for _, ex := range extras {
			if _, exists := totalsMap[ex.ExtraProductID]; !exists {
				ep, _ := s.repo.GetExtraProductByID(ctx, ex.ExtraProductID)
				var epModel *models.ExtraProduct
				if ep != nil {
					epModel = &models.ExtraProduct{ID: ep.ID, Name: ep.Name, Category: ep.Category}
				}
				totalsMap[ex.ExtraProductID] = &models.ExtraTotal{ExtraProduct: epModel, TotalQty: 0}
			}
			totalsMap[ex.ExtraProductID].TotalQty += ex.Quantity
		}
	}

	totals := &models.ExtraTotals{Date: date.Format("2006-01-02")}
	for _, t := range totalsMap {
		totals.Totals = append(totals.Totals, *t)
	}
	return totals, nil
}
