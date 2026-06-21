package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type OverviewHandler struct {
	serv service.Service
}

func NewOverviewHandler(serv service.Service) *OverviewHandler {
	return &OverviewHandler{serv: serv}
}

func (h *OverviewHandler) GetProductionOverview(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	var params dtos.ProductionOverviewFilters
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	from, err := time.Parse("2006-01-02", params.From)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid from date format, use YYYY-MM-DD", nil)
	}

	to, err := time.Parse("2006-01-02", params.To)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid to date format, use YYYY-MM-DD", nil)
	}

	if to.Before(from) {
		return respond(c, http.StatusBadRequest, "to date must be greater than or equal to from date", nil)
	}

	overview, err := h.serv.GetProductionOverview(ctx, from, to)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", overview)
}
