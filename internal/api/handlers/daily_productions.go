package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/entity"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type DailyProductionHandler struct {
	serv service.Service
}

func NewDailyProductionHandler(serv service.Service) *DailyProductionHandler {
	return &DailyProductionHandler{serv: serv}
}

func (h *DailyProductionHandler) Create(c *echo.Context) error {
	claims, err := requireStaff(c)
	if err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateDailyProduction
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	productionDate, err := time.Parse("2006-01-02", params.ProductionDate)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	lineInputs := make([]entity.ProductionLineInput, len(params.Lines))
	for i, l := range params.Lines {
		lineInputs[i] = entity.ProductionLineInput{MenuTypeID: l.MenuTypeID, Quantity: l.Quantity}
	}

	dp, err := h.serv.CreateDailyProduction(ctx, productionDate, params.CustomerID, params.FulfillmentType, params.DeliveryID, params.Notes, claims.UserID, lineInputs)
	if err != nil {
		if err == service.ErrInvalidFulfillment {
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "daily production created", dp)
}

func (h *DailyProductionHandler) GetByDate(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	dateStr := c.QueryParam("date")
	nameQuery := strings.TrimSpace(c.QueryParam("q"))
	fulfillmentType := strings.TrimSpace(c.QueryParam("fulfillment_type"))
	deliveryID := strings.TrimSpace(c.QueryParam("delivery_id"))
	menuTypeID := strings.TrimSpace(c.QueryParam("menu_type_id"))
	sortBy := strings.TrimSpace(c.QueryParam("sort"))
	sortOrder := strings.TrimSpace(c.QueryParam("order"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	if fulfillmentType != "" &&
		fulfillmentType != service.FulfillmentPending &&
		fulfillmentType != service.FulfillmentDelivery &&
		fulfillmentType != service.FulfillmentPickup {
		return respond(c, http.StatusBadRequest, "invalid fulfillment_type", nil)
	}

	if sortBy != "" && sortBy != "quantity" {
		return respond(c, http.StatusBadRequest, "invalid sort", nil)
	}
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return respond(c, http.StatusBadRequest, "invalid order", nil)
	}

	productions, err := h.serv.GetDailyProductions(ctx, date, nameQuery, fulfillmentType, deliveryID, menuTypeID, sortBy, sortOrder)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", productions)
}

func (h *DailyProductionHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("dailyProductionId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "dailyProductionId is required", nil)
	}

	dp, err := h.serv.GetDailyProductionByID(ctx, id)
	if err != nil {
		if err == service.ErrDailyProductionNotFound {
			return respond(c, http.StatusNotFound, "daily production not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", dp)
}

func (h *DailyProductionHandler) Update(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()

	var params dtos.UpdateDailyProduction
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateDailyProduction(ctx, params.ID, params.FulfillmentType, params.DeliveryID, params.Notes); err != nil {
		if err == service.ErrDailyProductionNotFound {
			return respond(c, http.StatusNotFound, "daily production not found", nil)
		}
		if err == service.ErrInvalidFulfillment {
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "daily production updated", nil)
}

func (h *DailyProductionHandler) Delete(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	var params dtos.DeleteDailyProduction
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteDailyProduction(ctx, params.ID); err != nil {
		if err == service.ErrDailyProductionNotFound {
			return respond(c, http.StatusNotFound, "daily production not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "daily production deleted", nil)
}

func (h *DailyProductionHandler) UpsertLine(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.UpsertDailyProductionLine
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	line, err := h.serv.UpsertDailyProductionLine(ctx, id, params.MenuTypeID, params.Quantity)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "line updated", line)
}

func (h *DailyProductionHandler) DeleteLine(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	lineID := c.Param("lineId")

	if err := h.serv.DeleteDailyProductionLine(ctx, lineID); err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "line deleted", nil)
}

func (h *DailyProductionHandler) AddExtra(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.AddDailyProductionExtra
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	extra, err := h.serv.AddDailyProductionExtra(ctx, id, params.ExtraProductID, params.Quantity)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "extra added", extra)
}

func (h *DailyProductionHandler) DeleteExtra(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	extraID := c.Param("extraId")

	if err := h.serv.DeleteDailyProductionExtra(ctx, extraID); err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "extra deleted", nil)
}

func (h *DailyProductionHandler) GetKitchenTotals(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	dateStr := c.QueryParam("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	totals, err := h.serv.GetKitchenTotals(ctx, date)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", totals)
}

func (h *DailyProductionHandler) GetExtrasTotals(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	dateStr := c.QueryParam("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return respond(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD", nil)
	}

	totals, err := h.serv.GetExtrasTotals(ctx, date)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", totals)
}
