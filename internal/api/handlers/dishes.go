package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type DishHandler struct {
	serv service.Service
}

func NewDishHandler(serv service.Service) *DishHandler {
	return &DishHandler{serv: serv}
}

func (h *DishHandler) Create(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateDish
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	dish, err := h.serv.CreateDish(ctx, params.Name, params.Description, params.MenuType)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "dish created", dish)
}

func (h *DishHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()

	menuType := c.QueryParam("menu_type")
	if menuType != "" {
		dishes, err := h.serv.GetDishesByMenuType(ctx, menuType)
		if err != nil {
			log.Println(err)
			return respond(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return respond(c, http.StatusOK, "ok", dishes)
	}

	dishes, err := h.serv.GetDishes(ctx)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", dishes)
}

func (h *DishHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("dishId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "dishId is required", nil)
	}

	dish, err := h.serv.GetDishByID(ctx, id)
	if err != nil {
		if err == service.ErrDishNotFound {
			return respond(c, http.StatusNotFound, "dish not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", dish)
}

func (h *DishHandler) Update(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.UpdateDish
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateDish(ctx, id, params.Name, params.Description, params.Active); err != nil {
		if err == service.ErrDishNotFound {
			return respond(c, http.StatusNotFound, "dish not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "dish updated", nil)
}

func (h *DishHandler) Delete(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.serv.DeleteDish(ctx, id); err != nil {
		if err == service.ErrDishNotFound {
			return respond(c, http.StatusNotFound, "dish not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "dish deleted", nil)
}
