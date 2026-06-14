package service

import (
	"context"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDailyProductionNotFound = errors.New("daily production not found")

func (s *serv) CreateDailyProduction(ctx context.Context, productionDate time.Time, customerID, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes, createdBy string) (*models.DailyProduction, error) {
	dp, err := s.repo.SaveDailyProduction(ctx, productionDate, customerID, deliveryID, traditionalQty, healthyQty, vegetarianQty, notes, createdBy)
	if err != nil {
		return nil, err
	}
	return &models.DailyProduction{
		ID:             dp.ID,
		ProductionDate: dp.ProductionDate.Format("2006-01-02"),
		TraditionalQty: dp.TraditionalQty,
		HealthyQty:     dp.HealthyQty,
		VegetarianQty:  dp.VegetarianQty,
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
				deliveryModel = &models.Delivery{ID: delivery.ID, Name: delivery.Name}
			}
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
			TraditionalQty: dp.TraditionalQty,
			HealthyQty:     dp.HealthyQty,
			VegetarianQty:  dp.VegetarianQty,
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
	return &models.DailyProduction{
		ID:             dp.ID,
		ProductionDate: dp.ProductionDate.Format("2006-01-02"),
		TraditionalQty: dp.TraditionalQty,
		HealthyQty:     dp.HealthyQty,
		VegetarianQty:  dp.VegetarianQty,
		Notes:          dp.Notes,
		CreatedBy:      dp.CreatedBy,
		CreatedAt:      dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateDailyProduction(ctx context.Context, id, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes string) error {
	_, err := s.repo.GetDailyProductionByID(ctx, id)
	if err != nil {
		return ErrDailyProductionNotFound
	}
	return s.repo.UpdateDailyProduction(ctx, id, deliveryID, traditionalQty, healthyQty, vegetarianQty, notes)
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

// GetMenuByDate obtiene el plato de cada tipo para una fecha específica.
func (s *serv) GetMenuByDate(ctx context.Context, date time.Time) (*models.DayMenu, error) {
	item, err := s.repo.GetWeekMenuItemByDate(ctx, date)
	if err != nil {
		return nil, errors.New("no menu configured for this date")
	}

	tDish, _ := s.repo.GetDishByID(ctx, item.TraditionalDishID)
	hDish, _ := s.repo.GetDishByID(ctx, item.HealthyDishID)
	vDish, _ := s.repo.GetDishByID(ctx, item.VegetarianDishID)

	dayMenu := &models.DayMenu{
		Date: date.Format("2006-01-02"),
	}
	if tDish != nil {
		dayMenu.TraditionalDish = &models.Dish{ID: tDish.ID, Name: tDish.Name, MenuType: tDish.MenuType}
	}
	if hDish != nil {
		dayMenu.HealthyDish = &models.Dish{ID: hDish.ID, Name: hDish.Name, MenuType: hDish.MenuType}
	}
	if vDish != nil {
		dayMenu.VegetarianDish = &models.Dish{ID: vDish.ID, Name: vDish.Name, MenuType: vDish.MenuType}
	}
	return dayMenu, nil
}

// GetKitchenTotals suma las cantidades por tipo de menú para una fecha.
func (s *serv) GetKitchenTotals(ctx context.Context, date time.Time) (*models.KitchenTotals, error) {
	productions, err := s.repo.GetDailyProductions(ctx, date)
	if err != nil {
		return nil, err
	}

	totals := &models.KitchenTotals{Date: date.Format("2006-01-02")}
	for _, dp := range productions {
		totals.TraditionalQty += dp.TraditionalQty
		totals.HealthyQty += dp.HealthyQty
		totals.VegetarianQty += dp.VegetarianQty
	}
	return totals, nil
}

// GetExtrasTotals agrupa y suma cantidades de extras para una fecha.
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
