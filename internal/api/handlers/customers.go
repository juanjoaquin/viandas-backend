package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type CustomerHandler struct {
	serv service.Service
}

func NewCustomerHandler(serv service.Service) *CustomerHandler {
	return &CustomerHandler{serv: serv}
}

func (h *CustomerHandler) Create(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.CreateCustomer
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	customer, err := h.serv.CreateCustomer(ctx, params.Name, params.Type, optionalString(params.Phone), optionalString(params.Address))
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "customer created", customer)
}

func (h *CustomerHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()

	typeFilter := c.QueryParam("type")
	if typeFilter != "" && typeFilter != "COMPANY" && typeFilter != "PERSON" {
		return respond(c, http.StatusBadRequest, "type must be COMPANY or PERSON", nil)
	}

	customers, err := h.serv.GetCustomers(ctx, c.QueryParam("q"), typeFilter)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", customers)
}

func (h *CustomerHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("customerId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "customerId is required", nil)
	}

	customer, err := h.serv.GetCustomerByID(ctx, id)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return respond(c, http.StatusNotFound, "customer not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", customer)
}

func (h *CustomerHandler) Update(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	id := c.Param("id")

	var params dtos.UpdateCustomer
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.UpdateCustomer(ctx, id, params.Name, params.Type, optionalString(params.Phone), optionalString(params.Address)); err != nil {
		if err == service.ErrCustomerNotFound {
			return respond(c, http.StatusNotFound, "customer not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "customer updated", nil)
}

func (h *CustomerHandler) Delete(c *echo.Context) error {
	if _, err := requireAdmin(c); err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()

	var params dtos.DeleteCustomer
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.DeleteCustomer(ctx, params.ID); err != nil {
		if err == service.ErrCustomerNotFound {
			return respond(c, http.StatusNotFound, "customer not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "customer deleted", nil)
}
