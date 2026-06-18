package repository

import (
	"context"

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

func (r *repo) GetCustomers(ctx context.Context, nameQuery, typeFilter string) ([]entity.Customer, error) {
	var customers []entity.Customer
	var err error

	switch {
	case nameQuery == "" && typeFilter == "":
		err = r.db.SelectContext(ctx, &customers, `SELECT * FROM customers ORDER BY name`)
	case nameQuery != "" && typeFilter == "":
		err = r.db.SelectContext(ctx, &customers,
			`SELECT * FROM customers WHERE name ILIKE $1 ORDER BY name`,
			"%"+nameQuery+"%",
		)
	case nameQuery == "" && typeFilter != "":
		err = r.db.SelectContext(ctx, &customers,
			`SELECT * FROM customers WHERE type = $1 ORDER BY name`,
			typeFilter,
		)
	default:
		err = r.db.SelectContext(ctx, &customers,
			`SELECT * FROM customers WHERE name ILIKE $1 AND type = $2 ORDER BY name`,
			"%"+nameQuery+"%", typeFilter,
		)
	}
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
