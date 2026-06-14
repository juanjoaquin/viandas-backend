package service

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/models"
	"github.com/juanjoaquin/viandas-backend/internal/repository"
)

//go:generate mockery --name=Service --output=service --inpackage=true
type Service interface {
	// Auth
	RegisterUser(ctx context.Context, name, email, password, role string) error
	LoginUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// Customers
	CreateCustomer(ctx context.Context, name, customerType string, phone, address *string) (*models.Customer, error)
	GetCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id, name, customerType string, phone, address *string) error
	DeleteCustomer(ctx context.Context, id string) error

	// Deliveries
	CreateDelivery(ctx context.Context, name string, phone *string) (*models.Delivery, error)
	GetDeliveries(ctx context.Context) ([]models.Delivery, error)
	GetDeliveryByID(ctx context.Context, id string) (*models.Delivery, error)
	UpdateDelivery(ctx context.Context, id, name string, phone *string, active bool) error
	DeleteDelivery(ctx context.Context, id string) error

	// Dishes
	CreateDish(ctx context.Context, name, description, menuType string) (*models.Dish, error)
	GetDishes(ctx context.Context) ([]models.Dish, error)
	GetDishesByMenuType(ctx context.Context, menuType string) ([]models.Dish, error)
	GetDishByID(ctx context.Context, id string) (*models.Dish, error)
	UpdateDish(ctx context.Context, id, name, description string, active bool) error
	DeleteDish(ctx context.Context, id string) error

	// ExtraProducts
	CreateExtraProduct(ctx context.Context, name, category string) (*models.ExtraProduct, error)
	GetExtraProducts(ctx context.Context) ([]models.ExtraProduct, error)
	GetExtraProductByID(ctx context.Context, id string) (*models.ExtraProduct, error)
	UpdateExtraProduct(ctx context.Context, id, name string, active bool) error
	DeleteExtraProduct(ctx context.Context, id string) error

	// WeekMenus
	CreateWeekMenu(ctx context.Context, weekStartDate time.Time, createdBy string) (*models.WeekMenu, error)
	GetWeekMenus(ctx context.Context) ([]models.WeekMenu, error)
	GetWeekMenuByID(ctx context.Context, id string) (*models.WeekMenu, error)
	AddWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, traditionalDishID, healthyDishID, vegetarianDishID string) (*models.WeekMenuItem, error)
	UpdateWeekMenuItem(ctx context.Context, id, traditionalDishID, healthyDishID, vegetarianDishID string) error
	DeleteWeekMenuItem(ctx context.Context, id string) error

	// DailyProductions
	CreateDailyProduction(ctx context.Context, productionDate time.Time, customerID, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes, createdBy string) (*models.DailyProduction, error)
	GetDailyProductions(ctx context.Context, date time.Time) ([]models.DailyProduction, error)
	GetDailyProductionByID(ctx context.Context, id string) (*models.DailyProduction, error)
	UpdateDailyProduction(ctx context.Context, id, deliveryID string, traditionalQty, healthyQty, vegetarianQty int, notes string) error
	AddDailyProductionExtra(ctx context.Context, dailyProductionID, extraProductID string, quantity int) (*models.DailyProductionExtra, error)
	DeleteDailyProductionExtra(ctx context.Context, id string) error

	// Consultas especiales
	GetMenuByDate(ctx context.Context, date time.Time) (*models.DayMenu, error)
	GetKitchenTotals(ctx context.Context, date time.Time) (*models.KitchenTotals, error)
	GetExtrasTotals(ctx context.Context, date time.Time) (*models.ExtraTotals, error)
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &serv{repo: repo}
}
