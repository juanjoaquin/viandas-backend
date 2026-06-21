package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type DeliveryHandler struct {
	serv                  service.Service
	paginatorLimitDefault string
}

func NewDeliveryHandler(serv service.Service, paginatorLimitDefault string) *DeliveryHandler {
	return &DeliveryHandler{serv: serv, paginatorLimitDefault: paginatorLimitDefault}
}

func (h *DeliveryHandler) Create(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateDelivery
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	delivery, err := h.serv.CreateDelivery(ctx, params.Name, optionalString(params.Phone))
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "delivery created", delivery)
}

func (h *DeliveryHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	nameQuery := c.QueryParam("q")

	return paginatedListResponse(c, h.paginatorLimitDefault,
		func() (int, error) {
			return h.serv.CountDeliveries(ctx, nameQuery)
		},
		func(offset, limit int) (interface{}, error) {
			return h.serv.GetDeliveries(ctx, nameQuery, offset, limit)
		},
	)
}

func (h *DeliveryHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("deliveryId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "deliveryId is required", nil)
	}

	delivery, err := h.serv.GetDeliveryByID(ctx, id)
	if err != nil {
		if err == service.ErrDeliveryNotFound {
			return respond(c, http.StatusNotFound, "delivery not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", delivery)
}

func (h *DeliveryHandler) Update(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.UpdateDelivery
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateDelivery(ctx, params.ID, params.Name, optionalString(params.Phone), params.Active); err != nil {
		if err == service.ErrDeliveryNotFound {
			return respond(c, http.StatusNotFound, "delivery not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "delivery updated", nil)
}

func (h *DeliveryHandler) Delete(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.DeleteDelivery
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteDelivery(ctx, params.ID); err != nil {
		if err == service.ErrDeliveryNotFound {
			return respond(c, http.StatusNotFound, "delivery not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "delivery deleted", nil)
}
