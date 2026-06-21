package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveDailyProduction(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string) (*entity.DailyProduction, error) {
	var dp entity.DailyProduction

	var deliveryArg interface{}
	if deliveryID != "" {
		deliveryArg = deliveryID
	}

	var notesArg interface{}
	if notes != "" {
		notesArg = notes
	}

	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO daily_productions (production_date, customer_id, fulfillment_type, delivery_id, notes, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		productionDate, customerID, fulfillmentType, deliveryArg, notesArg, createdBy,
	).StructScan(&dp)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *repo) SaveDailyProductionWithLines(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string, lines []entity.ProductionLineInput) (*entity.DailyProduction, []entity.DailyProductionLine, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	var deliveryArg interface{}
	if deliveryID != "" {
		deliveryArg = deliveryID
	}

	var notesArg interface{}
	if notes != "" {
		notesArg = notes
	}

	var dp entity.DailyProduction
	err = tx.QueryRowxContext(ctx,
		`INSERT INTO daily_productions (production_date, customer_id, fulfillment_type, delivery_id, notes, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		productionDate, customerID, fulfillmentType, deliveryArg, notesArg, createdBy,
	).StructScan(&dp)
	if err != nil {
		return nil, nil, err
	}

	savedLines := make([]entity.DailyProductionLine, 0, len(lines))
	for _, l := range lines {
		var line entity.DailyProductionLine
		err = tx.QueryRowxContext(ctx,
			`INSERT INTO daily_production_lines (daily_production_id, menu_type_id, quantity)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (daily_production_id, menu_type_id)
			 DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = NOW()
			 RETURNING *`,
			dp.ID, l.MenuTypeID, l.Quantity,
		).StructScan(&line)
		if err != nil {
			return nil, nil, err
		}
		savedLines = append(savedLines, line)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	return &dp, savedLines, nil
}

type dailyProductionListQuery struct {
	joins   string
	where   string
	orderBy string
	args    []interface{}
}

func buildDailyProductionListQuery(date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string) dailyProductionListQuery {
	args := []interface{}{date}
	joins := []string{}
	where := []string{"dp.production_date = $1"}

	if nameQuery != "" {
		joins = append(joins, "JOIN customers c ON c.id = dp.customer_id")
		args = append(args, "%"+nameQuery+"%")
		where = append(where, fmt.Sprintf("c.name ILIKE $%d", len(args)))
	}

	if fulfillmentType != "" {
		args = append(args, fulfillmentType)
		where = append(where, fmt.Sprintf("dp.fulfillment_type = $%d", len(args)))
	}

	if deliveryID != "" {
		args = append(args, deliveryID)
		where = append(where, fmt.Sprintf("dp.delivery_id = $%d", len(args)))
	}

	if menuTypeID != "" {
		args = append(args, menuTypeID)
		where = append(where, fmt.Sprintf(`EXISTS (
			SELECT 1
			FROM daily_production_lines dpl
			WHERE dpl.daily_production_id = dp.id
			  AND dpl.menu_type_id = $%d
		)`, len(args)))
	}

	orderBy := "dp.created_at"
	if sortBy == "quantity" {
		direction := "DESC"
		if sortOrder == "asc" {
			direction = "ASC"
		}
		orderBy = fmt.Sprintf(`(
			SELECT COALESCE(SUM(dpl.quantity), 0)
			FROM daily_production_lines dpl
			WHERE dpl.daily_production_id = dp.id
		) %s, dp.created_at`, direction)
	}

	return dailyProductionListQuery{
		joins:   strings.Join(joins, " "),
		where:   strings.Join(where, " AND "),
		orderBy: orderBy,
		args:    args,
	}
}

func (r *repo) CountDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string) (int, error) {
	q := buildDailyProductionListQuery(date, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder)
	query := fmt.Sprintf(
		`SELECT COUNT(*)
		 FROM daily_productions dp
		 %s
		 WHERE %s`,
		q.joins,
		q.where,
	)

	var count int
	err := r.db.GetContext(ctx, &count, query, q.args...)
	return count, err
}

func (r *repo) GetDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string, offset, limit int) ([]entity.DailyProduction, error) {
	q := buildDailyProductionListQuery(date, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder)
	args := append(q.args, limit, offset)
	query := fmt.Sprintf(
		`SELECT dp.*
		 FROM daily_productions dp
		 %s
		 WHERE %s
		 ORDER BY %s
		 LIMIT $%d OFFSET $%d`,
		q.joins,
		q.where,
		q.orderBy,
		len(q.args)+1,
		len(q.args)+2,
	)

	var productions []entity.DailyProduction
	err := r.db.SelectContext(ctx, &productions, query, args...)
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

func (r *repo) GetDailyProductionByDateAndCustomer(ctx context.Context, productionDate time.Time, customerID string) (*entity.DailyProduction, error) {
	var dp entity.DailyProduction
	err := r.db.GetContext(ctx, &dp,
		`SELECT * FROM daily_productions WHERE production_date = $1 AND customer_id = $2`,
		productionDate, customerID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dp, nil
}

func (r *repo) UpdateDailyProduction(ctx context.Context, id, fulfillmentType, deliveryID, notes string) error {
	var deliveryArg interface{}
	if deliveryID != "" {
		deliveryArg = deliveryID
	}

	var notesArg interface{}
	if notes != "" {
		notesArg = notes
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE daily_productions
		 SET fulfillment_type = $1,
		     delivery_id      = $2,
		     notes            = $3,
		     updated_at       = NOW()
		 WHERE id = $4`,
		fulfillmentType, deliveryArg, notesArg, id,
	)
	return err
}

func (r *repo) DeleteDailyProduction(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM daily_productions WHERE id = $1`, id)
	return err
}

func (r *repo) SaveDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*entity.DailyProductionLine, error) {
	var line entity.DailyProductionLine
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO daily_production_lines (daily_production_id, menu_type_id, quantity)
		 VALUES ($1, $2, $3) RETURNING *`,
		dailyProductionID, menuTypeID, quantity,
	).StructScan(&line)
	if err != nil {
		return nil, err
	}
	return &line, nil
}

func (r *repo) GetDailyProductionLines(ctx context.Context, dailyProductionID string) ([]entity.DailyProductionLine, error) {
	var lines []entity.DailyProductionLine
	err := r.db.SelectContext(ctx, &lines,
		`SELECT dpl.* FROM daily_production_lines dpl
		 JOIN menu_types mt ON mt.id = dpl.menu_type_id
		 WHERE dpl.daily_production_id = $1
		 ORDER BY mt.name`,
		dailyProductionID,
	)
	return lines, err
}

func (r *repo) UpsertDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*entity.DailyProductionLine, error) {
	var line entity.DailyProductionLine
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO daily_production_lines (daily_production_id, menu_type_id, quantity)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (daily_production_id, menu_type_id)
		 DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = NOW()
		 RETURNING *`,
		dailyProductionID, menuTypeID, quantity,
	).StructScan(&line)
	if err != nil {
		return nil, err
	}
	return &line, nil
}

func (r *repo) DeleteDailyProductionLine(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM daily_production_lines WHERE id = $1`, id)
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

func (r *repo) UpdateDailyProductionExtra(ctx context.Context, dailyProductionID, id, extraProductID string, quantity int) (*entity.DailyProductionExtra, error) {
	var extra entity.DailyProductionExtra
	err := r.db.QueryRowxContext(ctx,
		`UPDATE daily_production_extras
		 SET extra_product_id = $1,
		     quantity = $2,
		     updated_at = NOW()
		 WHERE id = $3 AND daily_production_id = $4
		 RETURNING *`,
		extraProductID, quantity, id, dailyProductionID,
	).StructScan(&extra)
	if err != nil {
		return nil, err
	}
	return &extra, nil
}

func (r *repo) DeleteDailyProductionExtra(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM daily_production_extras WHERE id = $1`, id)
	return err
}
