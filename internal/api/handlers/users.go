package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/encryption"
	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	authmw "github.com/juanjoaquin/viandas-backend/internal/api/middleware"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	serv service.Service
}

func NewUserHandler(serv service.Service) *UserHandler {
	return &UserHandler{serv: serv}
}

func respond(c *echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
	})
}

func (h *UserHandler) Register(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.RegisterUser
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.RegisterUser(ctx, params.Name, params.Email, params.Password, params.Role); err != nil {
		if err == service.ErrUserAlreadyExists {
			return respond(c, http.StatusConflict, err.Error(), nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "user created successfully", nil)
}

func (h *UserHandler) Login(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.LoginUser
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	user, err := h.serv.LoginUser(ctx, params.Email, params.Password)
	if err != nil {
		if err == service.ErrInvalidPassword || err == service.ErrUserNotFound {
			return respond(c, http.StatusUnauthorized, "invalid credentials", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	token, err := encryption.SignedLoginToken(user)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, "could not generate token", nil)
	}

	return respond(c, http.StatusOK, "login successful", map[string]string{
		"access_token": token,
	})
}

func (h *UserHandler) Me(c *echo.Context) error {
	ctx := c.Request().Context()

	authHeader := c.Request().Header.Get("Authorization")
	token, err := authmw.ExtractToken(authHeader)
	if err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	claims, err := encryption.ParseLoginJWT(token)
	if err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	user, err := h.serv.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return respond(c, http.StatusNotFound, "user not found", nil)
	}

	return respond(c, http.StatusOK, "ok", user)
}
