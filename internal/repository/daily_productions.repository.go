package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveDailyProduction(ctx context.Context, productionDate time.Time, customerID, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes, createdBy string) (*entity.DailyProduction, error) {
	var dp entity.DailyProduction

	var deliveryArg interface{}
	if deliveryID != "" {
		deliveryArg = deliveryID
	}

	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO daily_productions (production_date, customer_id, delivery_id, traditional_qty, healthy_qty, vegetarian_qty, notes, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
		productionDate, customerID, deliveryArg, traditionalQty, healthyQty, vegetarianQty, notes, createdBy,
	).StructScan(&dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *repo) GetDailyProductions(ctx context.Context, date time.Time) ([]entity.DailyProduction, error) {
	var productions []entity.DailyProduction
	err := r.db.SelectContext(ctx, &productions,
		`SELECT * FROM daily_productions WHERE production_date = $1 ORDER BY created_at`,
		date,
	)
	return productions, err
}

func (r *repo) GetDailyProductionByID(ctx context.Context, id string) (*entity.DailyProduction, error) {
	var dp entity.DailyProduction
	err := r.db.GetContext(ctx, &dp, `SELECT * FROM daily_productions WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *repo) UpdateDailyProduction(ctx context.Context, id, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes string) error {
	var deliveryArg interface{}
	if deliveryID != "" {
		deliveryArg = deliveryID
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE daily_productions SET delivery_id=$1, traditional_qty=$2, healthy_qty=$3, vegetarian_qty=$4, notes=$5, updated_at=NOW() WHERE id=$6`,
		deliveryArg, traditionalQty, healthyQty, vegetarianQty, notes, id,
	)
	return err
}

func (r *repo) DeleteDailyProduction(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM daily_productions WHERE id = $1`, id)
	return err
}

func (r *repo) SaveDailyProductionExtra(ctx context.Context, dailyProductionID, extraProductID string, quantity int) (*entity.DailyProductionExtra, error) {
	var extra entity.DailyProductionExtra
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO daily_production_extras (daily_production_id, extra_product_id, quantity) VALUES ($1, $2, $3) RETURNING *`,
		dailyProductionID, extraProductID, quantity,
	).StructScan(&extra)
	if err != nil {
		return nil, err
	}
	return &extra, nil
}

func (r *repo) GetDailyProductionExtras(ctx context.Context, dailyProductionID string) ([]entity.DailyProductionExtra, error) {
	var extras []entity.DailyProductionExtra
	err := r.db.SelectContext(ctx, &extras,
		`SELECT * FROM daily_production_extras WHERE daily_production_id = $1`,
		dailyProductionID,
	)
	return extras, err
}

func (r *repo) DeleteDailyProductionExtra(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM daily_production_extras WHERE id = $1`, id)
	return err
}
