package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

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
	if _, err := requireStaff(c); err != nil {
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
	if _, err := requireStaff(c); err != nil {
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
	if _, err := requireStaff(c); err != nil {
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

	// #region agent log
	func() {
		b, _ := json.Marshal(map[string]interface{}{
			"sessionId": "d8bacb", "hypothesisId": "F", "runId": "post-fix",
			"location": "menu_types.go:Delete", "message": "handler DeleteMenuType entry",
			"data": map[string]interface{}{"id": params.ID}, "timestamp": time.Now().UnixMilli(),
		})
		_ = os.MkdirAll("/home/juan/github/viandas-backend/.cursor", 0755)
		if f, err := os.OpenFile("/home/juan/github/viandas-backend/.cursor/debug-d8bacb.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			_, _ = f.Write(append(b, '\n'))
			_ = f.Close()
		}
		client := &http.Client{Timeout: 2 * time.Second}
		req, _ := http.NewRequest(http.MethodPost, "http://host.docker.internal:7677/ingest/4882f41c-573b-4ad0-8bdf-7268d958ad2e", bytes.NewReader(b))
		if req != nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Debug-Session-Id", "d8bacb")
			if resp, err := client.Do(req); err == nil {
				_ = resp.Body.Close()
			}
		}
	}()
	// #endregion

	if err := h.serv.DeleteMenuType(ctx, params.ID); err != nil {
		if err == service.ErrMenuTypeNotFound {
			return respond(c, http.StatusNotFound, "menu type not found", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "menu type deleted", nil)
}
