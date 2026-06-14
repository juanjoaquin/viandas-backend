package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type ExtraProductHandler struct {
	serv service.Service
}

func NewExtraProductHandler(serv service.Service) *ExtraProductHandler {
	return &ExtraProductHandler{serv: serv}
}

func (h *ExtraProductHandler) Create(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateExtraProduct
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	product, err := h.serv.CreateExtraProduct(ctx, params.Name, params.Category)
	if err != nil {
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
	products, err := h.serv.GetExtraProducts(ctx)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", products)
}

func (h *ExtraProductHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

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
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.UpdateExtraProduct
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateExtraProduct(ctx, id, params.Name, params.Active); err != nil {
		if err == service.ErrExtraProductNotFound {
			return respond(c, http.StatusNotFound, "extra product not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "extra product updated", nil)
}

func (h *ExtraProductHandler) Delete(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.serv.DeleteExtraProduct(ctx, id); err != nil {
		if err == service.ErrExtraProductNotFound {
			return respond(c, http.StatusNotFound, "extra product not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "extra product deleted", nil)
}
