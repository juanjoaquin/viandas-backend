package service

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

func (s *serv) GetProductionOverview(ctx context.Context, from, to time.Time) (*models.ProductionOverview, error) {
	menuTotals, err := s.repo.GetOverviewMenuTotals(ctx, from, to)
	if err != nil {
		return nil, err
	}

	productTotals, err := s.repo.GetOverviewProductTotals(ctx, from, to)
	if err != nil {
		return nil, err
	}

	overview := &models.ProductionOverview{
		Period: models.OverviewPeriod{
			From: from.Format("2006-01-02"),
			To:   to.Format("2006-01-02"),
		},
		MenusByType: make([]models.ProductionOverviewMenu, 0, len(menuTotals)),
		Products:    make([]models.ProductionOverviewProduct, 0, len(productTotals)),
	}

	for _, total := range menuTotals {
		overview.Summary.TotalMenusQty += total.TotalQty
		if total.TotalAmount != nil {
			overview.Summary.TotalMenusAmount += *total.TotalAmount
		}

		overview.MenusByType = append(overview.MenusByType, models.ProductionOverviewMenu{
			MenuType: &models.MenuType{
				ID:    total.MenuTypeID,
				Name:  total.MenuTypeName,
				Price: total.MenuTypePrice,
			},
			TotalQty:    total.TotalQty,
			TotalAmount: total.TotalAmount,
		})
	}

	for _, total := range productTotals {
		overview.Summary.TotalProductsQty += total.TotalQty
		overview.Summary.TotalProductsAmount += total.TotalAmount

		overview.Products = append(overview.Products, models.ProductionOverviewProduct{
			ExtraProduct: &models.ExtraProduct{
				ID:    total.ExtraProductID,
				Name:  total.ExtraProductName,
				Price: total.ExtraProductPrice,
				Category: &models.ProductCategory{
					ID:   total.CategoryID,
					Name: total.CategoryName,
				},
			},
			TotalQty:    total.TotalQty,
			TotalAmount: total.TotalAmount,
		})
	}

	overview.Summary.GrandTotalAmount = overview.Summary.TotalMenusAmount + overview.Summary.TotalProductsAmount

	return overview, nil
}
