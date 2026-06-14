package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrCustomerNotFound = errors.New("customer not found")

func (s *serv) CreateCustomer(ctx context.Context, name, customerType, phone, address string) (*models.Customer, error) {
	c, err := s.repo.SaveCustomer(ctx, name, customerType, phone, address)
	if err != nil {
		return nil, err
	}
	return &models.Customer{
		ID:        c.ID,
		Name:      c.Name,
		Type:      c.Type,
		Phone:     c.Phone,
		Address:   c.Address,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) GetCustomers(ctx context.Context) ([]models.Customer, error) {
	entities, err := s.repo.GetCustomers(ctx)
	if err != nil {
		return nil, err
	}

	customers := make([]models.Customer, len(entities))
	for i, c := range entities {
		customers[i] = models.Customer{
			ID:        c.ID,
			Name:      c.Name,
			Type:      c.Type,
			Phone:     c.Phone,
			Address:   c.Address,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return customers, nil
}

func (s *serv) GetCustomerByID(ctx context.Context, id string) (*models.Customer, error) {
	c, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, ErrCustomerNotFound
	}
	return &models.Customer{
		ID:        c.ID,
		Name:      c.Name,
		Type:      c.Type,
		Phone:     c.Phone,
		Address:   c.Address,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateCustomer(ctx context.Context, id, name, customerType, phone, address string) error {
	_, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return ErrCustomerNotFound
	}
	return s.repo.UpdateCustomer(ctx, id, name, customerType, phone, address)
}

func (s *serv) DeleteCustomer(ctx context.Context, id string) error {
	_, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return ErrCustomerNotFound
	}
	return s.repo.DeleteCustomer(ctx, id)
}
