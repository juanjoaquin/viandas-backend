package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDailyProductionNotFound = errors.New("daily production not found")
var ErrDailyProductionExtraNotFound = errors.New("daily production extra not found")
var ErrDailyProductionAlreadyExists = errors.New("este cliente ya tiene una producción para esta fecha")
var ErrInvalidFulfillment = errors.New("PICKUP cannot have a delivery_id")

const (
	FulfillmentPending  = "PENDING"
	FulfillmentDelivery = "DELIVERY"
	FulfillmentPickup   = "PICKUP"
)

func normalizeFulfillment(ft string) string {
	switch ft {
	case FulfillmentDelivery, FulfillmentPickup, FulfillmentPending:
		return ft
	default:
		return FulfillmentPending
	}
}

func (s *serv) CreateDailyProduction(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string, lines []entity.ProductionLineInput) (*models.DailyProduction, error) {
	fulfillmentType = normalizeFulfillment(fulfillmentType)

	if fulfillmentType == FulfillmentPickup && deliveryID != "" {
		return nil, ErrInvalidFulfillment
	}
	if fulfillmentType != FulfillmentDelivery {
		deliveryID = ""
	}

	existing, err := s.repo.GetDailyProductionByDateAndCustomer(ctx, productionDate, customerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDailyProductionAlreadyExists
	}

	var dp *entity.DailyProduction
	var savedLines []entity.DailyProductionLine

	if len(lines) == 0 {
		var err error
		dp, err = s.repo.SaveDailyProduction(ctx, productionDate, customerID, fulfillmentType, deliveryID, notes, createdBy)
		if err != nil {
			if isDailyProductionDuplicate(err) {
				return nil, ErrDailyProductionAlreadyExists
			}
			return nil, err
		}
	} else {
		var err error
		dp, savedLines, err = s.repo.SaveDailyProductionWithLines(ctx, productionDate, customerID, fulfillmentType, deliveryID, notes, createdBy, lines)
		if err != nil {
			if isDailyProductionDuplicate(err) {
				return nil, ErrDailyProductionAlreadyExists
			}
			return nil, err
		}
	}

	lineModels := make([]models.DailyProductionLine, 0, len(savedLines))
	for _, line := range savedLines {
		lm := models.DailyProductionLine{
			ID:       line.ID,
			Quantity: line.Quantity,
		}
		if mt, err := s.repo.GetMenuTypeByID(ctx, line.MenuTypeID); err == nil {
			lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
		}
		lineModels = append(lineModels, lm)
	}

	result := &models.DailyProduction{
		ID:              dp.ID,
		ProductionDate:  dp.ProductionDate.Format("2006-01-02"),
		FulfillmentType: dp.FulfillmentType,
		Notes:           dp.Notes,
		CreatedBy:       dp.CreatedBy,
		CreatedAt:       dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if len(lineModels) > 0 {
		result.Lines = lineModels
	}
	return result, nil
}

func (s *serv) CountDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string) (int, error) {
	return s.repo.CountDailyProductions(ctx, date, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder)
}

func (s *serv) GetDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string, offset, limit int) ([]models.DailyProduction, error) {
	entities, err := s.repo.GetDailyProductions(ctx, date, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder, offset, limit)
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
		if dp.DeliveryID != nil && *dp.DeliveryID != "" {
			delivery, _ := s.repo.GetDeliveryByID(ctx, *dp.DeliveryID)
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
				lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
			}
			lineModels[j] = lm
		}

		extras, _ := s.repo.GetDailyProductionExtras(ctx, dp.ID)
		extraModels := make([]models.DailyProductionExtra, len(extras))
		for j, ex := range extras {
			ep, _ := s.repo.GetExtraProductByID(ctx, ex.ExtraProductID)
			epModel := s.mapExtraProduct(ctx, ep)
			var totalAmount *float64
			if epModel != nil {
				totalAmount = extraLineTotalAmount(epModel.Price, ex.Quantity)
			}
			extraModels[j] = models.DailyProductionExtra{
				ID:           ex.ID,
				ExtraProduct: epModel,
				Quantity:     ex.Quantity,
				TotalAmount:  totalAmount,
			}
		}

		productions[i] = models.DailyProduction{
			ID:              dp.ID,
			ProductionDate:  dp.ProductionDate.Format("2006-01-02"),
			FulfillmentType: dp.FulfillmentType,
			Customer:        customerModel,
			Delivery:        deliveryModel,
			Lines:           lineModels,
			Notes:           dp.Notes,
			Extras:          extraModels,
			CreatedBy:       dp.CreatedBy,
			CreatedAt:       dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
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
			lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
		}
		lineModels[j] = lm
	}

	return &models.DailyProduction{
		ID:              dp.ID,
		ProductionDate:  dp.ProductionDate.Format("2006-01-02"),
		FulfillmentType: dp.FulfillmentType,
		Lines:           lineModels,
		Notes:           dp.Notes,
		CreatedBy:       dp.CreatedBy,
		CreatedAt:       dp.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateDailyProduction(ctx context.Context, id string, fulfillmentType, deliveryID, notes *string) error {
	current, err := s.repo.GetDailyProductionByID(ctx, id)
	if err != nil {
		return ErrDailyProductionNotFound
	}

	// Partir del estado actual y aplicar solo los campos enviados
	ft := current.FulfillmentType
	did := ""
	if current.DeliveryID != nil {
		did = *current.DeliveryID
	}
	n := ""
	if current.Notes != nil {
		n = *current.Notes
	}

	if fulfillmentType != nil {
		ft = normalizeFulfillment(*fulfillmentType)
	}
	if deliveryID != nil {
		did = *deliveryID
	}
	if notes != nil {
		n = *notes
	}

	// Validar estado final
	if ft == FulfillmentPickup && did != "" {
		return ErrInvalidFulfillment
	}
	if ft != FulfillmentDelivery {
		did = ""
	}

	return s.repo.UpdateDailyProduction(ctx, id, ft, did, n)
}

func (s *serv) DeleteDailyProduction(ctx context.Context, id string) error {
	if _, err := s.repo.GetDailyProductionByID(ctx, id); err != nil {
		return ErrDailyProductionNotFound
	}

	return s.repo.DeleteDailyProduction(ctx, id)
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
		lm.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
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
	epModel := s.mapExtraProduct(ctx, ep)
	var totalAmount *float64
	if epModel != nil {
		totalAmount = extraLineTotalAmount(epModel.Price, extra.Quantity)
	}

	return &models.DailyProductionExtra{
		ID:           extra.ID,
		ExtraProduct: epModel,
		Quantity:     extra.Quantity,
		TotalAmount:  totalAmount,
	}, nil
}

func (s *serv) UpdateDailyProductionExtra(ctx context.Context, dailyProductionID, id, extraProductID string, quantity int) (*models.DailyProductionExtra, error) {
	extra, err := s.repo.UpdateDailyProductionExtra(ctx, dailyProductionID, id, extraProductID, quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDailyProductionExtraNotFound
		}
		return nil, err
	}

	ep, _ := s.repo.GetExtraProductByID(ctx, extra.ExtraProductID)
	epModel := s.mapExtraProduct(ctx, ep)
	var totalAmount *float64
	if epModel != nil {
		totalAmount = extraLineTotalAmount(epModel.Price, extra.Quantity)
	}

	return &models.DailyProductionExtra{
		ID:           extra.ID,
		ExtraProduct: epModel,
		Quantity:     extra.Quantity,
		TotalAmount:  totalAmount,
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
			line.MenuType = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
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
	productions, err := s.repo.GetDailyProductions(ctx, date, "", "", "", "", "", "", 0, unlimitedListLimit)
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
					mtModel = &models.MenuType{ID: mt.ID, Name: mt.Name, Price: mt.Price}
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

	var grandTotal float64
	hasAnyPrice := false

	for _, t := range totalsMap {
		if t.MenuType != nil && t.MenuType.Price != nil {
			amount := float64(t.TotalQty) * *t.MenuType.Price
			t.TotalAmount = &amount
			grandTotal += amount
			hasAnyPrice = true
		}
		totals.Totals = append(totals.Totals, *t)
	}

	if hasAnyPrice {
		totals.GrandTotal = &grandTotal
	}

	return totals, nil
}

func (s *serv) GetExtrasTotals(ctx context.Context, date time.Time) (*models.ExtraTotals, error) {
	productions, err := s.repo.GetDailyProductions(ctx, date, "", "", "", "", "", "", 0, unlimitedListLimit)
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
				epModel := s.mapExtraProduct(ctx, ep)
				totalsMap[ex.ExtraProductID] = &models.ExtraTotal{ExtraProduct: epModel, TotalQty: 0}
			}
			totalsMap[ex.ExtraProductID].TotalQty += ex.Quantity
		}
	}

	totals := &models.ExtraTotals{
		Date:   date.Format("2006-01-02"),
		Totals: make([]models.ExtraTotal, 0, len(totalsMap)),
	}

	var grandTotal float64
	hasAnyPrice := false

	for _, t := range totalsMap {
		if t.ExtraProduct != nil && t.ExtraProduct.Price > 0 {
			amount := float64(t.TotalQty) * t.ExtraProduct.Price
			t.TotalAmount = &amount
			grandTotal += amount
			hasAnyPrice = true
		}
		totals.Totals = append(totals.Totals, *t)
	}

	if hasAnyPrice {
		totals.GrandTotal = &grandTotal
	}

	return totals, nil
}
