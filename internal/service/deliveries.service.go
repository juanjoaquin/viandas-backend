package service

import (
	"context"
	"errors"

	"github.com/juanjoaquin/viandas-backend/internal/models"
)

var ErrDeliveryNotFound = errors.New("delivery not found")

func (s *serv) CreateDelivery(ctx context.Context, name string, phone *string) (*models.Delivery, error) {
	d, err := s.repo.SaveDelivery(ctx, name, phone)
	if err != nil {
		return nil, err
	}
	return &models.Delivery{
		ID:        d.ID,
		Name:      d.Name,
		Phone:     d.Phone,
		Active:    d.Active,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) CountDeliveries(ctx context.Context, nameQuery string) (int, error) {
	return s.repo.CountDeliveries(ctx, nameQuery)
}

func (s *serv) GetDeliveries(ctx context.Context, nameQuery string, offset, limit int) ([]models.Delivery, error) {
	entities, err := s.repo.GetDeliveries(ctx, nameQuery, offset, limit)
	if err != nil {
		return nil, err
	}

	deliveries := make([]models.Delivery, len(entities))
	for i, d := range entities {
		deliveries[i] = models.Delivery{
			ID:        d.ID,
			Name:      d.Name,
			Phone:     d.Phone,
			Active:    d.Active,
			CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return deliveries, nil
}

func (s *serv) GetDeliveryByID(ctx context.Context, id string) (*models.Delivery, error) {
	d, err := s.repo.GetDeliveryByID(ctx, id)
	if err != nil {
		return nil, ErrDeliveryNotFound
	}
	return &models.Delivery{
		ID:        d.ID,
		Name:      d.Name,
		Phone:     d.Phone,
		Active:    d.Active,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *serv) UpdateDelivery(ctx context.Context, id, name string, phone *string, active bool) error {
	_, err := s.repo.GetDeliveryByID(ctx, id)
	if err != nil {
		return ErrDeliveryNotFound
	}
	return s.repo.UpdateDelivery(ctx, id, name, phone, active)
}

func (s *serv) DeleteDelivery(ctx context.Context, id string) error {
	_, err := s.repo.GetDeliveryByID(ctx, id)
	if err != nil {
		return ErrDeliveryNotFound
	}
	return s.repo.DeleteDelivery(ctx, id)
}
