package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/lib/pq"
)

// Se fixeo el menu types que daba error SQL

func (r *repo) SaveMenuType(ctx context.Context, name string, price *float64) (*entity.MenuType, error) {
	var mt entity.MenuType
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO menu_types (name, price) VALUES ($1, $2) RETURNING *`,
		name, price,
	).StructScan(&mt)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func buildMenuTypeWhere(nameQuery string, activeFilter *bool) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if nameQuery != "" {
		args = append(args, "%"+nameQuery+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if activeFilter != nil {
		args = append(args, *activeFilter)
		conditions = append(conditions, fmt.Sprintf("active = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *repo) CountMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool) (int, error) {
	where, args := buildMenuTypeWhere(nameQuery, activeFilter)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM menu_types`+where, args...)
	return count, err
}

func (r *repo) GetMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]entity.MenuType, error) {
	where, args := buildMenuTypeWhere(nameQuery, activeFilter)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM menu_types%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var types []entity.MenuType
	err := r.db.SelectContext(ctx, &types, query, args...)
	return types, err
}

func (r *repo) GetMenuTypeByID(ctx context.Context, id string) (*entity.MenuType, error) {
	var mt entity.MenuType
	err := r.db.GetContext(ctx, &mt, `SELECT * FROM menu_types WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repo) UpdateMenuType(ctx context.Context, id, name string, price *float64, active bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu_types SET name=$1, price=$2, active=$3, updated_at=NOW() WHERE id=$4`,
		name, price, active, id,
	)
	return err
}

func agentDebugLog(hypothesisId, location, message string, data map[string]interface{}) {
	_ = os.MkdirAll("/home/juan/github/viandas-backend/.cursor", 0755)
	f, err := os.OpenFile("/home/juan/github/viandas-backend/.cursor/debug-d8bacb.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	b, _ := json.Marshal(map[string]interface{}{
		"sessionId": "d8bacb", "hypothesisId": hypothesisId, "runId": "post-fix",
		"location": location, "message": message, "data": data,
		"timestamp": time.Now().UnixMilli(),
	})
	_, _ = f.Write(append(b, '\n'))
}

func pqMeta(err error) (string, string, string) {
	if err == nil {
		return "", "", ""
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return err.Error(), string(pqErr.Code), pqErr.Constraint
	}
	return err.Error(), "", ""
}

func (r *repo) DeleteMenuType(ctx context.Context, id string) error {
	// #region agent log
	var dishesN, wmiN, dplN int
	_ = r.db.GetContext(ctx, &dishesN, `SELECT COUNT(*) FROM dishes WHERE menu_type_id = $1`, id)
	_ = r.db.GetContext(ctx, &wmiN, `SELECT COUNT(*) FROM week_menu_items WHERE menu_type_id = $1 OR dish_id IN (SELECT id FROM dishes WHERE menu_type_id = $1)`, id)
	_ = r.db.GetContext(ctx, &dplN, `SELECT COUNT(*) FROM daily_production_lines WHERE menu_type_id = $1`, id)
	agentDebugLog("A", "menu_types.repository.go:DeleteMenuType", "related rows before DELETE menu_types", map[string]interface{}{
		"id": id, "dishes": dishesN, "week_menu_items": wmiN, "daily_production_lines": dplN,
	})
	// #endregion

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
		DELETE FROM week_menu_items
		 WHERE menu_type_id = $1
		    OR dish_id IN (SELECT id FROM dishes WHERE menu_type_id = $1)`, id); err != nil {
		errStr, pgCode, constraint := pqMeta(err)
		agentDebugLog("B", "menu_types.repository.go:DeleteMenuType:week_menu_items", "failed deleting week_menu_items", map[string]interface{}{
			"id": id, "err": errStr, "pgCode": pgCode, "constraint": constraint,
		})
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM daily_production_lines WHERE menu_type_id = $1`, id); err != nil {
		errStr, pgCode, constraint := pqMeta(err)
		agentDebugLog("C", "menu_types.repository.go:DeleteMenuType:daily_production_lines", "failed deleting daily_production_lines", map[string]interface{}{
			"id": id, "err": errStr, "pgCode": pgCode, "constraint": constraint,
		})
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM dishes WHERE menu_type_id = $1`, id); err != nil {
		errStr, pgCode, constraint := pqMeta(err)
		agentDebugLog("A", "menu_types.repository.go:DeleteMenuType:dishes", "failed deleting dishes", map[string]interface{}{
			"id": id, "err": errStr, "pgCode": pgCode, "constraint": constraint,
		})
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM menu_types WHERE id = $1`, id); err != nil {
		errStr, pgCode, constraint := pqMeta(err)
		agentDebugLog("A", "menu_types.repository.go:DeleteMenuType:after", "DELETE menu_types result", map[string]interface{}{
			"id": id, "err": errStr, "pgCode": pgCode, "constraint": constraint,
		})
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	// #region agent log
	agentDebugLog("A", "menu_types.repository.go:DeleteMenuType:after", "DELETE menu_types result", map[string]interface{}{
		"id": id, "err": "", "pgCode": "", "constraint": "", "committed": true,
		"relatedBefore": map[string]int{"dishes": dishesN, "week_menu_items": wmiN, "daily_production_lines": dplN},
	})
	// #endregion
	return nil
}
