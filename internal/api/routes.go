package api

import (
	"github.com/juanjoaquin/viandas-backend/internal/api/handlers"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

func (a *API) RegisterRoutes(e *echo.Echo, serv service.Service, paginatorLimitDefault string) {
	userH := handlers.NewUserHandler(serv)
	customerH := handlers.NewCustomerHandler(serv, paginatorLimitDefault)
	deliveryH := handlers.NewDeliveryHandler(serv, paginatorLimitDefault)
	menuTypeH := handlers.NewMenuTypeHandler(serv, paginatorLimitDefault)
	dishH := handlers.NewDishHandler(serv, paginatorLimitDefault)
	productCategoryH := handlers.NewProductCategoryHandler(serv, paginatorLimitDefault)
	extraH := handlers.NewExtraProductHandler(serv, paginatorLimitDefault)
	weekMenuH := handlers.NewWeekMenuHandler(serv)
	dailyH := handlers.NewDailyProductionHandler(serv, paginatorLimitDefault)
	overviewH := handlers.NewOverviewHandler(serv)

	// Auth
	auth := e.Group("/auth")
	auth.POST("/register", userH.Register)
	auth.POST("/register-with-invite", userH.RegisterWithInvite)
	auth.POST("/login", userH.Login)
	auth.GET("/me", userH.Me)
	auth.POST("/refresh", userH.Refresh)
	auth.POST("/logout", userH.Logout)

	// Users
	users := e.Group("/users")
	users.POST("/invites", userH.Invite)

	// Customers
	customers := e.Group("/customers")
	customers.POST("", customerH.Create)
	customers.GET("", customerH.GetAll) // ?q=<name>
	customers.GET("/one", customerH.GetByID)
	customers.PUT("/:id", customerH.Update)
	customers.DELETE("", customerH.Delete)

	// Deliveries
	deliveries := e.Group("/deliveries")
	deliveries.POST("", deliveryH.Create)
	deliveries.GET("", deliveryH.GetAll) // ?q=<name>
	deliveries.GET("/one", deliveryH.GetByID)
	deliveries.PUT("", deliveryH.Update)
	deliveries.DELETE("", deliveryH.Delete)

	// Menu Types
	menuTypes := e.Group("/menu-types")
	menuTypes.POST("", menuTypeH.Create)
	menuTypes.GET("", menuTypeH.GetAll)
	menuTypes.GET("/one", menuTypeH.GetByID)
	menuTypes.PUT("", menuTypeH.Update)
	menuTypes.DELETE("", menuTypeH.Delete)

	// Dishes
	dishes := e.Group("/dishes")
	dishes.POST("", dishH.Create)
	dishes.GET("", dishH.GetAll) // ?q=<name>&menu_type_id=<uuid>
	dishes.GET("/one", dishH.GetByID)
	dishes.PUT("", dishH.Update)
	dishes.DELETE("", dishH.Delete)

	// Product Categories
	productCategories := e.Group("/product-categories")
	productCategories.POST("", productCategoryH.Create)
	productCategories.GET("", productCategoryH.GetAll)
	productCategories.GET("/one", productCategoryH.GetByID)
	productCategories.PUT("", productCategoryH.Update)
	productCategories.DELETE("", productCategoryH.Delete)

	// Extra Products
	extras := e.Group("/extra-products")
	extras.POST("", extraH.Create)
	extras.GET("", extraH.GetAll)
	extras.GET("/one", extraH.GetByID)
	extras.PUT("", extraH.Update)
	extras.DELETE("", extraH.Delete)

	// Week Menus
	weekMenus := e.Group("/week-menus")
	weekMenus.POST("", weekMenuH.Create)
	weekMenus.GET("", weekMenuH.GetAll)
	weekMenus.GET("/menu", weekMenuH.GetMenuByDate) // ?date=2026-06-16
	weekMenus.GET("/one", weekMenuH.GetByID)
	weekMenus.GET("/resolved", weekMenuH.Resolve)
	weekMenus.DELETE("/:id", weekMenuH.Delete)
	weekMenus.POST("/:id/items", weekMenuH.AddItem)
	weekMenus.PUT("/:id/items/:itemId", weekMenuH.UpdateItem)
	weekMenus.DELETE("/:id/items/:itemId", weekMenuH.DeleteItem)

	// Daily Productions
	daily := e.Group("/daily-productions")
	daily.POST("", dailyH.Create)
	daily.GET("", dailyH.GetByDate)                       // ?date=2026-06-16
	daily.GET("/totals/kitchen", dailyH.GetKitchenTotals) // ?date=2026-06-16
	daily.GET("/totals/extras", dailyH.GetExtrasTotals)   // ?date=2026-06-16
	daily.GET("/one", dailyH.GetByID)
	daily.PUT("", dailyH.Update)
	daily.DELETE("", dailyH.Delete)
	daily.PUT("/:id/lines", dailyH.UpsertLine)
	daily.DELETE("/:id/lines/:lineId", dailyH.DeleteLine)
	daily.POST("/:id/extras", dailyH.AddExtra)
	daily.PUT("/:id/extras/:extraId", dailyH.UpdateExtra)
	daily.DELETE("/:id/extras/:extraId", dailyH.DeleteExtra)

	// Overview
	overview := e.Group("/overview")
	overview.GET("/production", overviewH.GetProductionOverview) // ?from=2026-06-01&to=2026-06-30
}
