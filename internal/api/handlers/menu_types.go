package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type MenuTypeHandler struct {
	serv service.Service
}

func NewMenuTypeHandler(serv service.Service) *MenuTypeHandler {
	return &MenuTypeHandler{serv: serv}
}

func (h *MenuTypeHandler) Create(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateMenuType
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	mt, err := h.serv.CreateMenuType(ctx, params.Name, params.Price)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "menu type created", mt)
}

func (h *MenuTypeHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	types, err := h.serv.GetMenuTypes(ctx)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", types)
}

func (h *MenuTypeHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	mt, err := h.serv.GetMenuTypeByID(ctx, id)
	if err != nil {
		if err == service.ErrMenuTypeNotFound {
			return respond(c, http.StatusNotFound, "menu type not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", mt)
}

func (h *MenuTypeHandler) Update(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.UpdateMenuType
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateMenuType(ctx, id, params.Name, params.Price, params.Active); err != nil {
		if err == service.ErrMenuTypeNotFound {
			return respond(c, http.StatusNotFound, "menu type not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "menu type updated", nil)
}

func (h *MenuTypeHandler) Delete(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.serv.DeleteMenuType(ctx, id); err != nil {
		if err == service.ErrMenuTypeNotFound {
			return respond(c, http.StatusNotFound, "menu type not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "menu type deleted", nil)
}
