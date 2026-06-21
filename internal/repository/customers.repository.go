package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

func (r *repo) SaveCustomer(ctx context.Context, name, customerType string, phone, address *string) (*entity.Customer, error) {
	var c entity.Customer
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO customers (name, type, phone, address) VALUES ($1, $2, $3, $4) RETURNING *`,
		name, customerType, phone, address,
	).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func buildCustomerWhere(nameQuery, typeFilter string) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if nameQuery != "" {
		args = append(args, "%"+nameQuery+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if typeFilter != "" {
		args = append(args, typeFilter)
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *repo) CountCustomers(ctx context.Context, nameQuery, typeFilter string) (int, error) {
	where, args := buildCustomerWhere(nameQuery, typeFilter)
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM customers`+where, args...)
	return count, err
}

func (r *repo) GetCustomers(ctx context.Context, nameQuery, typeFilter string, offset, limit int) ([]entity.Customer, error) {
	where, args := buildCustomerWhere(nameQuery, typeFilter)
	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT * FROM customers%s ORDER BY name LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	var customers []entity.Customer
	err := r.db.SelectContext(ctx, &customers, query, args...)
	return customers, err
}

func (r *repo) GetCustomerByID(ctx context.Context, id string) (*entity.Customer, error) {
	var c entity.Customer
	err := r.db.GetContext(ctx, &c, `SELECT * FROM customers WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repo) UpdateCustomer(ctx context.Context, id, name, customerType string, phone, address *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE customers SET name=$1, type=$2, phone=$3, address=$4, updated_at=NOW() WHERE id=$5`,
		name, customerType, phone, address, id,
	)
	return err
}

func (r *repo) DeleteCustomer(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1`, id)
	return err
}
