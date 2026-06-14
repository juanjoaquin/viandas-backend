package api

import (
	"github.com/juanjoaquin/viandas-backend/internal/api/handlers"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

func (a *API) RegisterRoutes(e *echo.Echo, serv service.Service) {
	userH := handlers.NewUserHandler(serv)
	customerH := handlers.NewCustomerHandler(serv)
	deliveryH := handlers.NewDeliveryHandler(serv)
	menuTypeH := handlers.NewMenuTypeHandler(serv)
	dishH := handlers.NewDishHandler(serv)
	extraH := handlers.NewExtraProductHandler(serv)
	weekMenuH := handlers.NewWeekMenuHandler(serv)
	dailyH := handlers.NewDailyProductionHandler(serv)

	// Auth
	auth := e.Group("/auth")
	auth.POST("/register", userH.Register)
	auth.POST("/login", userH.Login)
	auth.GET("/me", userH.Me)

	// Customers
	customers := e.Group("/customers")
	customers.POST("", customerH.Create)
	customers.GET("", customerH.GetAll)
	customers.GET("/one", customerH.GetByID)
	customers.PUT("/:id", customerH.Update)
	customers.DELETE("/:id", customerH.Delete)

	// Deliveries
	deliveries := e.Group("/deliveries")
	deliveries.POST("", deliveryH.Create)
	deliveries.GET("", deliveryH.GetAll)
	deliveries.GET("/one", deliveryH.GetByID)
	deliveries.PUT("/:id", deliveryH.Update)
	deliveries.DELETE("/:id", deliveryH.Delete)

	// Menu Types
	menuTypes := e.Group("/menu-types")
	menuTypes.POST("", menuTypeH.Create)
	menuTypes.GET("", menuTypeH.GetAll)
	menuTypes.GET("/:id", menuTypeH.GetByID)
	menuTypes.PUT("/:id", menuTypeH.Update)
	menuTypes.DELETE("/:id", menuTypeH.Delete)

	// Dishes
	dishes := e.Group("/dishes")
	dishes.POST("", dishH.Create)
	dishes.GET("", dishH.GetAll) // ?menu_type_id=<uuid>
	dishes.GET("/one", dishH.GetByID)
	dishes.PUT("/:id", dishH.Update)
	dishes.DELETE("/:id", dishH.Delete)

	// Extra Products
	extras := e.Group("/extra-products")
	extras.POST("", extraH.Create)
	extras.GET("", extraH.GetAll)
	extras.GET("/one", extraH.GetByID)
	extras.PUT("/:id", extraH.Update)
	extras.DELETE("/:id", extraH.Delete)

	// Week Menus
	weekMenus := e.Group("/week-menus")
	weekMenus.POST("", weekMenuH.Create)
	weekMenus.GET("", weekMenuH.GetAll)
	weekMenus.GET("/menu", weekMenuH.GetMenuByDate) // ?date=2026-06-16
	weekMenus.GET("/one", weekMenuH.GetByID)
	weekMenus.POST("/:id/items", weekMenuH.AddItem)
	weekMenus.PUT("/:id/items/:itemId", weekMenuH.UpdateItem)
	weekMenus.DELETE("/:id/items/:itemId", weekMenuH.DeleteItem)

	// Daily Productions
	daily := e.Group("/daily-productions")
	daily.POST("", dailyH.Create)
	daily.GET("", dailyH.GetByDate)                         // ?date=2026-06-16
	daily.GET("/totals/kitchen", dailyH.GetKitchenTotals)   // ?date=2026-06-16
	daily.GET("/totals/extras", dailyH.GetExtrasTotals)     // ?date=2026-06-16
	daily.GET("/one", dailyH.GetByID)
	daily.PUT("/:id", dailyH.Update)
	daily.PUT("/:id/lines", dailyH.UpsertLine)
	daily.DELETE("/:id/lines/:lineId", dailyH.DeleteLine)
	daily.POST("/:id/extras", dailyH.AddExtra)
	daily.DELETE("/:id/extras/:extraId", dailyH.DeleteExtra)
}
