package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type ProductCategoryHandler struct {
	serv                  service.Service
	paginatorLimitDefault string
}

func NewProductCategoryHandler(serv service.Service, paginatorLimitDefault string) *ProductCategoryHandler {
	return &ProductCategoryHandler{serv: serv, paginatorLimitDefault: paginatorLimitDefault}
}

func (h *ProductCategoryHandler) Create(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateProductCategory
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	category, err := h.serv.CreateProductCategory(ctx, params.Name)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "product category created", category)
}

func (h *ProductCategoryHandler) GetAll(c *echo.Context) error {
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
			return h.serv.CountProductCategories(ctx, nameQuery, activeFilter)
		},
		func(offset, limit int) (interface{}, error) {
			return h.serv.GetProductCategories(ctx, nameQuery, activeFilter, offset, limit)
		},
	)
}

func (h *ProductCategoryHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("productCategoryId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "productCategoryId is required", nil)
	}

	category, err := h.serv.GetProductCategoryByID(ctx, id)
	if err != nil {
		if err == service.ErrProductCategoryNotFound {
			return respond(c, http.StatusNotFound, "product category not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", category)
}

func (h *ProductCategoryHandler) Update(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.UpdateProductCategory
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateProductCategory(ctx, params.ID, params.Name, params.Active); err != nil {
		if err == service.ErrProductCategoryNotFound {
			return respond(c, http.StatusNotFound, "product category not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "product category updated", nil)
}

func (h *ProductCategoryHandler) Delete(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.DeleteProductCategory
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteProductCategory(ctx, params.ID); err != nil {
		if err == service.ErrProductCategoryNotFound {
			return respond(c, http.StatusNotFound, "product category not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "product category deleted", nil)
}
