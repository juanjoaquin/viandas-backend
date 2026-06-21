package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type MenuTypeHandler struct {
	serv                  service.Service
	paginatorLimitDefault string
}

func NewMenuTypeHandler(serv service.Service, paginatorLimitDefault string) *MenuTypeHandler {
	return &MenuTypeHandler{serv: serv, paginatorLimitDefault: paginatorLimitDefault}
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

	activeParam := c.QueryParam("active")
	var activeFilter *bool
	if activeParam != "" {
		if activeParam != "true" && activeParam != "false" {
			return respond(c, http.StatusBadRequest, "active must be true or false", nil)
		}
		active := activeParam == "true"
		activeFilter = &active
	}

	nameQuery := c.QueryParam("q")

	return paginatedListResponse(c, h.paginatorLimitDefault,
		func() (int, error) {
			return h.serv.CountMenuTypes(ctx, nameQuery, activeFilter)
		},
		func(offset, limit int) (interface{}, error) {
			return h.serv.GetMenuTypes(ctx, nameQuery, activeFilter, offset, limit)
		},
	)
}

func (h *MenuTypeHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("menuTypeId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "menuTypeId is required", nil)
	}

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

	var params dtos.UpdateMenuType
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateMenuType(ctx, params.ID, params.Name, params.Price, params.Active); err != nil {
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

	var params dtos.DeleteMenuType
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteMenuType(ctx, params.ID); err != nil {
		if err == service.ErrMenuTypeNotFound {
			return respond(c, http.StatusNotFound, "menu type not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "menu type deleted", nil)
}
