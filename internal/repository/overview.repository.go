package repository

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) GetOverviewMenuTotals(ctx context.Context, from, to time.Time) ([]entity.OverviewMenuTotal, error) {
	var totals []entity.OverviewMenuTotal
	err := r.db.SelectContext(ctx, &totals,
		`SELECT
		     mt.id AS menu_type_id,
		     mt.name AS menu_type_name,
		     mt.price AS menu_type_price,
		     COALESCE(SUM(dpl.quantity), 0)::int AS total_qty,
		     CASE
		         WHEN mt.price IS NULL THEN NULL
		         ELSE (COALESCE(SUM(dpl.quantity), 0) * mt.price)::float
		     END AS total_amount
		 FROM daily_productions dp
		 JOIN daily_production_lines dpl ON dpl.daily_production_id = dp.id
		 JOIN menu_types mt ON mt.id = dpl.menu_type_id
		 WHERE dp.production_date BETWEEN $1 AND $2
		 GROUP BY mt.id, mt.name, mt.price
		 ORDER BY mt.name`,
		from, to,
	)
	return totals, err
}

func (r *repo) GetOverviewProductTotals(ctx context.Context, from, to time.Time) ([]entity.OverviewProductTotal, error) {
	var totals []entity.OverviewProductTotal
	err := r.db.SelectContext(ctx, &totals,
		`SELECT
		     ep.id AS extra_product_id,
		     ep.name AS extra_product_name,
		     ep.price AS extra_product_price,
		     pc.id AS category_id,
		     pc.name AS category_name,
		     COALESCE(SUM(dpe.quantity), 0)::int AS total_qty,
		     (COALESCE(SUM(dpe.quantity), 0) * ep.price)::float AS total_amount
		 FROM daily_productions dp
		 JOIN daily_production_extras dpe ON dpe.daily_production_id = dp.id
		 JOIN extra_products ep ON ep.id = dpe.extra_product_id
		 JOIN product_categories pc ON pc.id = ep.category_id
		 WHERE dp.production_date BETWEEN $1 AND $2
		 GROUP BY ep.id, ep.name, ep.price, pc.id, pc.name
		 ORDER BY ep.name`,
		from, to,
	)
	return totals, err
}
