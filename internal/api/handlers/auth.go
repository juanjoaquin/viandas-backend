package handlers

import (
	"github.com/juanjoaquin/viandas-backend/encryption"
	authmw "github.com/juanjoaquin/viandas-backend/internal/api/middleware"
	"github.com/juanjoaquin/viandas-backend/internal/roles"
	"github.com/labstack/echo/v5"
)

func getTokenClaims(c *echo.Context, allowed ...string) (*encryption.Claims, error) {
	token, err := authmw.ExtractToken(c.Request().Header.Get("Authorization"))
	if err != nil {
		return nil, echo.ErrUnauthorized
	}
	return authmw.RequireRole(token, allowed...)
}

func requireAdmin(c *echo.Context) (*encryption.Claims, error) {
	return getTokenClaims(c, roles.Admin)
}

func requireStaff(c *echo.Context) (*encryption.Claims, error) {
	return getTokenClaims(c, roles.Admin, roles.Employee)
}
