package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type ExtraProductHandler struct {
	serv                  service.Service
	paginatorLimitDefault string
}

func NewExtraProductHandler(serv service.Service, paginatorLimitDefault string) *ExtraProductHandler {
	return &ExtraProductHandler{serv: serv, paginatorLimitDefault: paginatorLimitDefault}
}

func (h *ExtraProductHandler) Create(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateExtraProduct
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	product, err := h.serv.CreateExtraProduct(ctx, params.Name, params.CategoryID, params.Price)
	if err != nil {
		if err == service.ErrProductCategoryNotFound {
			return respond(c, http.StatusNotFound, "product category not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "extra product created", product)
}

func (h *ExtraProductHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	nameQuery := c.QueryParam("q")

	return paginatedListResponse(c, h.paginatorLimitDefault,
		func() (int, error) {
			return h.serv.CountExtraProducts(ctx, nameQuery)
		},
		func(offset, limit int) (interface{}, error) {
			return h.serv.GetExtraProducts(ctx, nameQuery, offset, limit)
		},
	)
}

func (h *ExtraProductHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("extraProductId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "extraProductId is required", nil)
	}

	product, err := h.serv.GetExtraProductByID(ctx, id)
	if err != nil {
		if err == service.ErrExtraProductNotFound {
			return respond(c, http.StatusNotFound, "extra product not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", product)
}

func (h *ExtraProductHandler) Update(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.UpdateExtraProduct
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateExtraProduct(ctx, params.ID, params.Name, params.CategoryID, params.Price, params.Active); err != nil {
		if err == service.ErrExtraProductNotFound {
			return respond(c, http.StatusNotFound, "extra product not found", nil)
		}
		if err == service.ErrProductCategoryNotFound {
			return respond(c, http.StatusNotFound, "product category not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "extra product updated", nil)
}

func (h *ExtraProductHandler) Delete(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.DeleteExtraProduct
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteExtraProduct(ctx, params.ID); err != nil {
		if err == service.ErrExtraProductNotFound {
			return respond(c, http.StatusNotFound, "extra product not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "extra product deleted", nil)
}
