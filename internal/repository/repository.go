package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/juanjoaquin/viandas-backend/internal/entity"
)

//go:generate mockery --name=Repository --output=repository --inpackage=true
type Repository interface {
	// Users
	SaveUser(ctx context.Context, name, email, passwordHash, role string) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// RefreshTokens
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error

	// Customers
	SaveCustomer(ctx context.Context, name, customerType string, phone, address *string) (*entity.Customer, error)
	GetCustomers(ctx context.Context, nameQuery, typeFilter string) ([]entity.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, id, name, customerType string, phone, address *string) error
	DeleteCustomer(ctx context.Context, id string) error

	// Deliveries
	SaveDelivery(ctx context.Context, name string, phone *string) (*entity.Delivery, error)
	GetDeliveries(ctx context.Context, nameQuery string) ([]entity.Delivery, error)
	GetDeliveryByID(ctx context.Context, id string) (*entity.Delivery, error)
	UpdateDelivery(ctx context.Context, id, name string, phone *string, active bool) error
	DeleteDelivery(ctx context.Context, id string) error

	// MenuTypes
	SaveMenuType(ctx context.Context, name string, price *float64) (*entity.MenuType, error)
	GetMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool) ([]entity.MenuType, error)
	GetMenuTypeByID(ctx context.Context, id string) (*entity.MenuType, error)
	UpdateMenuType(ctx context.Context, id, name string, price *float64, active bool) error
	DeleteMenuType(ctx context.Context, id string) error

	// ProductCategories
	SaveProductCategory(ctx context.Context, name string) (*entity.ProductCategory, error)
	GetProductCategories(ctx context.Context, nameQuery string, activeFilter *bool) ([]entity.ProductCategory, error)
	GetProductCategoryByID(ctx context.Context, id string) (*entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, id, name string, active bool) error
	DeleteProductCategory(ctx context.Context, id string) error

	// Dishes
	SaveDish(ctx context.Context, name, description, menuTypeID string) (*entity.Dish, error)
	GetDishes(ctx context.Context, nameQuery string) ([]entity.Dish, error)
	GetDishesByMenuTypeID(ctx context.Context, menuTypeID string) ([]entity.Dish, error)
	GetDishByID(ctx context.Context, id string) (*entity.Dish, error)
	UpdateDish(ctx context.Context, id, name, description, menuTypeID string, active bool) error
	DeleteDish(ctx context.Context, id string) error

	// ExtraProducts
	SaveExtraProduct(ctx context.Context, name, categoryID string, price float64) (*entity.ExtraProduct, error)
	GetExtraProducts(ctx context.Context, nameQuery string) ([]entity.ExtraProduct, error)
	GetExtraProductByID(ctx context.Context, id string) (*entity.ExtraProduct, error)
	UpdateExtraProduct(ctx context.Context, id, name, categoryID string, price float64, active bool) error
	DeleteExtraProduct(ctx context.Context, id string) error

	// WeekMenus
	SaveWeekMenu(ctx context.Context, weekStartDate, weekEndDate time.Time, createdBy string) (*entity.WeekMenu, error)
	GetWeekMenus(ctx context.Context) ([]entity.WeekMenu, error)
	GetWeekMenuByID(ctx context.Context, id string) (*entity.WeekMenu, error)
	GetWeekMenuByDate(ctx context.Context, date time.Time) (*entity.WeekMenu, error)
	DeleteWeekMenu(ctx context.Context, id string) error

	// WeekMenuItems
	SaveWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, menuTypeID, dishID string) (*entity.WeekMenuItem, error)
	GetWeekMenuItems(ctx context.Context, weekMenuID string) ([]entity.WeekMenuItem, error)
	GetWeekMenuItemsByDate(ctx context.Context, date time.Time) ([]entity.WeekMenuItem, error)
	UpdateWeekMenuItem(ctx context.Context, id, dishID string) error
	DeleteWeekMenuItem(ctx context.Context, id string) error

	// DailyProductions
	SaveDailyProduction(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string) (*entity.DailyProduction, error)
	SaveDailyProductionWithLines(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string, lines []entity.ProductionLineInput) (*entity.DailyProduction, []entity.DailyProductionLine, error)
	GetDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string) ([]entity.DailyProduction, error)
	GetDailyProductionByID(ctx context.Context, id string) (*entity.DailyProduction, error)
	GetDailyProductionByDateAndCustomer(ctx context.Context, productionDate time.Time, customerID string) (*entity.DailyProduction, error)
	UpdateDailyProduction(ctx context.Context, id, fulfillmentType, deliveryID, notes string) error
	DeleteDailyProduction(ctx context.Context, id string) error

	// DailyProductionLines
	SaveDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*entity.DailyProductionLine, error)
	GetDailyProductionLines(ctx context.Context, dailyProductionID string) ([]entity.DailyProductionLine, error)
	UpsertDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*entity.DailyProductionLine, error)
	DeleteDailyProductionLine(ctx context.Context, id string) error

	// DailyProductionExtras
	SaveDailyProductionExtra(ctx context.Context, dailyProductionID, extraProductID string, quantity int) (*entity.DailyProductionExtra, error)
	GetDailyProductionExtras(ctx context.Context, dailyProductionID string) ([]entity.DailyProductionExtra, error)
	UpdateDailyProductionExtra(ctx context.Context, dailyProductionID, id, extraProductID string, quantity int) (*entity.DailyProductionExtra, error)
	DeleteDailyProductionExtra(ctx context.Context, id string) error
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}
