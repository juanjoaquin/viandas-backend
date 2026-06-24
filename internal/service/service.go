package service

import (
	"context"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/models"
	"github.com/juanjoaquin/viandas-backend/internal/repository"
	"github.com/juanjoaquin/viandas-backend/settings"
)

type Mailer interface {
	SendInvite(ctx context.Context, toEmail, inviteURL string) error
	SendPasswordReset(ctx context.Context, toEmail, resetURL string) error
}

//go:generate mockery --name=Service --output=service --inpackage=true
type Service interface {
	// Auth
	RegisterUser(ctx context.Context, name, email, password, role string) error
	InviteUser(ctx context.Context, email, role, invitedBy string) (*models.UserInvite, error)
	RegisterWithInvite(ctx context.Context, token, name, password string) error
	LoginUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CountUsers(ctx context.Context, nameQuery string, activeFilter *bool) (int, error)
	GetUsers(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]models.User, error)
	UpdateUserActive(ctx context.Context, id string, active bool, requestingUserID string) error
	CreateRefreshToken(ctx context.Context, userID string) (string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	// Customers
	CreateCustomer(ctx context.Context, name, customerType string, phone, address *string) (*models.Customer, error)
	CountCustomers(ctx context.Context, nameQuery, typeFilter string) (int, error)
	GetCustomers(ctx context.Context, nameQuery, typeFilter string, offset, limit int) ([]models.Customer, error)
	GetCustomerByID(ctx context.Context, id string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id, name, customerType string, phone, address *string) error
	DeleteCustomer(ctx context.Context, id string) error

	// Deliveries
	CreateDelivery(ctx context.Context, name string, phone *string) (*models.Delivery, error)
	CountDeliveries(ctx context.Context, nameQuery string) (int, error)
	GetDeliveries(ctx context.Context, nameQuery string, offset, limit int) ([]models.Delivery, error)
	GetDeliveryByID(ctx context.Context, id string) (*models.Delivery, error)
	UpdateDelivery(ctx context.Context, id, name string, phone *string, active bool) error
	DeleteDelivery(ctx context.Context, id string) error

	// MenuTypes
	CreateMenuType(ctx context.Context, name string, price *float64) (*models.MenuType, error)
	CountMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool) (int, error)
	GetMenuTypes(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]models.MenuType, error)
	GetMenuTypeByID(ctx context.Context, id string) (*models.MenuType, error)
	UpdateMenuType(ctx context.Context, id, name string, price *float64, active bool) error
	DeleteMenuType(ctx context.Context, id string) error

	// ProductCategories
	CreateProductCategory(ctx context.Context, name string) (*models.ProductCategory, error)
	CountProductCategories(ctx context.Context, nameQuery string, activeFilter *bool) (int, error)
	GetProductCategories(ctx context.Context, nameQuery string, activeFilter *bool, offset, limit int) ([]models.ProductCategory, error)
	GetProductCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, id, name string, active bool) error
	DeleteProductCategory(ctx context.Context, id string) error

	// Dishes
	CreateDish(ctx context.Context, name, description, menuTypeID string) (*models.Dish, error)
	CountDishes(ctx context.Context, nameQuery, menuTypeID string) (int, error)
	GetDishes(ctx context.Context, nameQuery, menuTypeID string, offset, limit int) ([]models.Dish, error)
	GetDishesByMenuTypeID(ctx context.Context, menuTypeID string, offset, limit int) ([]models.Dish, error)
	GetDishByID(ctx context.Context, id string) (*models.Dish, error)
	UpdateDish(ctx context.Context, id, name, description, menuTypeID string, active bool) error
	DeleteDish(ctx context.Context, id string) error

	// ExtraProducts
	CreateExtraProduct(ctx context.Context, name, categoryID string, price float64) (*models.ExtraProduct, error)
	CountExtraProducts(ctx context.Context, nameQuery string) (int, error)
	GetExtraProducts(ctx context.Context, nameQuery string, offset, limit int) ([]models.ExtraProduct, error)
	GetExtraProductByID(ctx context.Context, id string) (*models.ExtraProduct, error)
	UpdateExtraProduct(ctx context.Context, id, name, categoryID string, price float64, active bool) error
	DeleteExtraProduct(ctx context.Context, id string) error

	// WeekMenus
	CreateWeekMenu(ctx context.Context, weekStartDate, weekEndDate time.Time, createdBy string) (*models.WeekMenu, error)
	GetWeekMenus(ctx context.Context) ([]models.WeekMenu, error)
	GetWeekMenuByID(ctx context.Context, id string) (*models.WeekMenu, error)
	ResolveWeekMenu(ctx context.Context, requestedID string, date time.Time) (*models.WeekMenu, error)
	DeleteWeekMenu(ctx context.Context, id string) error
	AddWeekMenuItem(ctx context.Context, weekMenuID string, menuDate time.Time, menuTypeID, dishID string) (*models.WeekMenuItem, error)
	UpdateWeekMenuItem(ctx context.Context, id, dishID string) error
	DeleteWeekMenuItem(ctx context.Context, id string) error

	// DailyProductions
	CreateDailyProduction(ctx context.Context, productionDate time.Time, customerID, fulfillmentType, deliveryID, notes, createdBy string, lines []entity.ProductionLineInput) (*models.DailyProduction, error)
	CountDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string) (int, error)
	GetDailyProductions(ctx context.Context, date time.Time, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder string, offset, limit int) ([]models.DailyProduction, error)
	GetDailyProductionByID(ctx context.Context, id string) (*models.DailyProduction, error)
	UpdateDailyProduction(ctx context.Context, id string, fulfillmentType, deliveryID, notes *string) error
	DeleteDailyProduction(ctx context.Context, id string) error
	UpsertDailyProductionLine(ctx context.Context, dailyProductionID, menuTypeID string, quantity int) (*models.DailyProductionLine, error)
	DeleteDailyProductionLine(ctx context.Context, id string) error
	AddDailyProductionExtra(ctx context.Context, dailyProductionID, extraProductID string, quantity int) (*models.DailyProductionExtra, error)
	UpdateDailyProductionExtra(ctx context.Context, dailyProductionID, id, extraProductID string, quantity int) (*models.DailyProductionExtra, error)
	DeleteDailyProductionExtra(ctx context.Context, id string) error

	// Consultas especiales
	GetMenuByDate(ctx context.Context, date time.Time) (*models.DayMenu, error)
	GetKitchenTotals(ctx context.Context, date time.Time) (*models.KitchenTotals, error)
	GetExtrasTotals(ctx context.Context, date time.Time) (*models.ExtraTotals, error)
	GetProductionOverview(ctx context.Context, from, to time.Time) (*models.ProductionOverview, error)
}

type serv struct {
	repo     repository.Repository
	mailer   Mailer
	settings *settings.Settings
}

func New(repo repository.Repository, mailer Mailer, settings *settings.Settings) Service {
	return &serv{
		repo:     repo,
		mailer:   mailer,
		settings: settings,
	}
}
