package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/juanjoaquin/viandas-backend/internal/pkg/logger"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/juanjoaquin/viandas-backend/settings"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type DataResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

type API struct {
	serv          service.Service
	settings      *settings.Settings
	dataValidator *validator.Validate
}

func New(serv service.Service, settings *settings.Settings) *API {
	return &API{
		serv:          serv,
		settings:      settings,
		dataValidator: validator.New(),
	}
}

func (a *API) Start(e *echo.Echo, address string) error {
	e.Use(middleware.RequestID())
	e.Use(logger.RequestLogger(a.settings.LogFormat))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: false,
	}))

	a.RegisterRoutes(e, a.serv, a.settings.PaginatorLimitDefault)

	return e.Start(address)
}

func (a *API) ok(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, DataResponse{
		Message: message,
		Code:    http.StatusOK,
		Data:    data,
	})
}

func (a *API) created(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, DataResponse{
		Message: message,
		Code:    http.StatusCreated,
		Data:    data,
	})
}

func (a *API) badRequest(c *echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, DataResponse{
		Message: message,
		Code:    http.StatusBadRequest,
	})
}

func (a *API) notFound(c *echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, DataResponse{
		Message: message,
		Code:    http.StatusNotFound,
	})
}

func (a *API) unauthorized(c *echo.Context) error {
	return c.JSON(http.StatusUnauthorized, DataResponse{
		Message: "unauthorized",
		Code:    http.StatusUnauthorized,
	})
}

func (a *API) forbidden(c *echo.Context) error {
	return c.JSON(http.StatusForbidden, DataResponse{
		Message: "forbidden: insufficient permissions",
		Code:    http.StatusForbidden,
	})
}

func (a *API) internalError(c *echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, DataResponse{
		Message: err.Error(),
		Code:    http.StatusInternalServerError,
	})
}
