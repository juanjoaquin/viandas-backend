package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type WeekMenuHandler struct {
	serv service.Service
}

func NewWeekMenuHandler(serv service.Service) *WeekMenuHandler {
	return &WeekMenuHandler{serv: serv}
}

func (h *WeekMenuHandler) Create(c *echo.Context) error {
	claims, err := requireAdmin(c)
	if err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateWeekMenu
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	startDate, err := time.Parse("2006-01-02", params.WeekStartDate)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	menu, err := h.serv.CreateWeekMenu(ctx, startDate, claims.UserID)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "week menu created", menu)
}

func (h *WeekMenuHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	menus, err := h.serv.GetWeekMenus(ctx)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", menus)
}

func (h *WeekMenuHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	menu, err := h.serv.GetWeekMenuByID(ctx, id)
	if err != nil {
		if err == service.ErrWeekMenuNotFound {
			return respond(c, http.StatusNotFound, "week menu not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", menu)
}

func (h *WeekMenuHandler) AddItem(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	weekMenuID := c.Param("id")

	var params dtos.AddWeekMenuItem
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	menuDate, err := time.Parse("2006-01-02", params.MenuDate)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	item, err := h.serv.AddWeekMenuItem(ctx, weekMenuID, menuDate, params.TraditionalDishID, params.HealthyDishID, params.VegetarianDishID)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "week menu item added", item)
}

func (h *WeekMenuHandler) UpdateItem(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	itemID := c.Param("itemId")

	var params dtos.UpdateWeekMenuItem
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateWeekMenuItem(ctx, itemID, params.TraditionalDishID, params.HealthyDishID, params.VegetarianDishID); err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "week menu item updated", nil)
}

func (h *WeekMenuHandler) DeleteItem(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	itemID := c.Param("itemId")

	if err := h.serv.DeleteWeekMenuItem(ctx, itemID); err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "week menu item deleted", nil)
}

func (h *WeekMenuHandler) GetMenuByDate(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	dateStr := c.QueryParam("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	menu, err := h.serv.GetMenuByDate(ctx, date)
	if err != nil {
		return respond(c, http.StatusNotFound, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", menu)
}
