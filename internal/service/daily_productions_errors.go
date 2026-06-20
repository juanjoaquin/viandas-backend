package service

import (
	"errors"

	"github.com/lib/pq"
)

func isDailyProductionDuplicate(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" &&
			pqErr.Constraint == "daily_productions_production_date_customer_id_key"
	}
	return false
}
