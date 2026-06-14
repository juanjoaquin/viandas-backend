package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type DeliveryHandler struct {
	serv service.Service
}

func NewDeliveryHandler(serv service.Service) *DeliveryHandler {
	return &DeliveryHandler{serv: serv}
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

	delivery, err := h.serv.CreateDelivery(ctx, params.Name)
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
	deliveries, err := h.serv.GetDeliveries(ctx)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", deliveries)
}

func (h *DeliveryHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

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
	id := c.Param("id")

	var params dtos.UpdateDelivery
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateDelivery(ctx, id, params.Name, params.Active); err != nil {
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
	id := c.Param("id")

	if err := h.serv.DeleteDelivery(ctx, id); err != nil {
		if err == service.ErrDeliveryNotFound {
			return respond(c, http.StatusNotFound, "delivery not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "delivery deleted", nil)
}
