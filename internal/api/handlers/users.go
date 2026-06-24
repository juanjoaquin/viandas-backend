package handlers

import (
	"log"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/encryption"
	"github.com/juanjoaquin/viandas-backend/internal/api/dtos"
	authmw "github.com/juanjoaquin/viandas-backend/internal/api/middleware"
	"github.com/juanjoaquin/viandas-backend/internal/roles"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	serv                  service.Service
	paginatorLimitDefault string
}

func NewUserHandler(serv service.Service, paginatorLimitDefault string) *UserHandler {
	return &UserHandler{serv: serv, paginatorLimitDefault: paginatorLimitDefault}
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

	role := params.Role
	if role == "" {
		role = roles.Employee
	}
	if role != roles.Admin && role != roles.Employee {
		return respond(c, http.StatusBadRequest, "role must be ADMIN or EMPLOYEE", nil)
	}

	if err := h.serv.RegisterUser(ctx, params.Name, params.Email, params.Password, role); err != nil {
		if err == service.ErrUserAlreadyExists {
			return respond(c, http.StatusConflict, err.Error(), nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusCreated, "user created successfully", nil)
}

func (h *UserHandler) Invite(c *echo.Context) error {
	claims, err := requireAdmin(c)
	if err != nil {
		return respond(c, http.StatusForbidden, "forbidden", nil)
	}

	ctx := c.Request().Context()
	var params dtos.InviteUser
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.Role == "" {
		params.Role = roles.Employee
	}

	invite, err := h.serv.InviteUser(ctx, params.Email, params.Role, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			return respond(c, http.StatusConflict, err.Error(), nil)
		case service.ErrInviteRoleNotSupported:
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Println(err)
			return respond(c, http.StatusInternalServerError, err.Error(), nil)
		}
	}

	return respond(c, http.StatusCreated, "invite created successfully", invite)
}

func (h *UserHandler) RegisterWithInvite(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.RegisterWithInvite
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.RegisterWithInvite(ctx, params.Token, params.Name, params.Password); err != nil {
		switch err {
		case service.ErrInvalidInvite, service.ErrInviteExpired, service.ErrInviteAlreadyAccepted:
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		case service.ErrUserAlreadyExists:
			return respond(c, http.StatusConflict, err.Error(), nil)
		default:
			log.Println(err)
			return respond(c, http.StatusInternalServerError, err.Error(), nil)
		}
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
		if err == service.ErrUserInactive {
			return respond(c, http.StatusForbidden, "account disabled", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	accessToken, err := encryption.SignedLoginToken(user)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, "could not generate token", nil)
	}

	refreshToken, err := h.serv.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, "could not generate refresh token", nil)
	}

	return respond(c, http.StatusOK, "login successful", map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *UserHandler) Refresh(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.RefreshTokenRequest
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	accessToken, newRefreshToken, err := h.serv.RefreshAccessToken(ctx, params.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidRefreshToken {
			return respond(c, http.StatusUnauthorized, "invalid or expired refresh token", nil)
		}
		if err == service.ErrUserInactive {
			return respond(c, http.StatusForbidden, "account disabled", nil)
		}
		log.Println(err)
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "token refreshed", map[string]string{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

func (h *UserHandler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.RefreshTokenRequest
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.RevokeRefreshToken(ctx, params.RefreshToken); err != nil {
		log.Println(err)
		return respond(c, http.StatusInternalServerError, "could not revoke token", nil)
	}

	return respond(c, http.StatusOK, "logged out successfully", nil)
}

func (h *UserHandler) ForgotPassword(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.ForgotPassword
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	// Siempre 200 para no revelar si el email existe
	if err := h.serv.RequestPasswordReset(ctx, params.Email); err != nil {
		log.Println(err)
	}

	return respond(c, http.StatusOK, "if that email exists, a reset link has been sent", nil)
}

func (h *UserHandler) ResetPassword(c *echo.Context) error {
	ctx := c.Request().Context()
	var params dtos.ResetPassword
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}

	if err := h.serv.ResetPassword(ctx, params.Token, params.Password); err != nil {
		switch err {
		case service.ErrPasswordResetTokenUsed:
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		case service.ErrPasswordResetTokenInvalid:
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Println(err)
			return respond(c, http.StatusInternalServerError, err.Error(), nil)
		}
	}

	return respond(c, http.StatusOK, "password reset successfully", nil)
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

func (h *UserHandler) GetAll(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	nameQuery := c.QueryParam("q")

	activeParam := c.QueryParam("active")
	var activeFilter *bool
	if activeParam != "" {
		if activeParam != "true" && activeParam != "false" {
			return respond(c, http.StatusBadRequest, "active must be true or false", nil)
		}
		active := activeParam == "true"
		activeFilter = &active
	}

	return paginatedListResponse(c, h.paginatorLimitDefault,
		func() (int, error) {
			return h.serv.CountUsers(ctx, nameQuery, activeFilter)
		},
		func(offset, limit int) (interface{}, error) {
			return h.serv.GetUsers(ctx, nameQuery, activeFilter, offset, limit)
		},
	)
}

func (h *UserHandler) GetByID(c *echo.Context) error {
	if _, err := requireStaff(c); err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()
	id := c.QueryParam("userId")
	if id == "" {
		return respond(c, http.StatusBadRequest, "userId is required", nil)
	}

	user, err := h.serv.GetUserByID(ctx, id)
	if err != nil {
		if err == service.ErrUserNotFound {
			return respond(c, http.StatusNotFound, "user not found", nil)
		}
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respond(c, http.StatusOK, "ok", user)
}

func (h *UserHandler) Update(c *echo.Context) error {
	claims, err := requireStaff(c)
	if err != nil {
		return respond(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	ctx := c.Request().Context()

	var params dtos.UpdateUser
	if err := c.Bind(&params); err != nil {
		return respond(c, http.StatusBadRequest, err.Error(), nil)
	}
	if params.ID == "" {
		return respond(c, http.StatusBadRequest, "id is required", nil)
	}

	if err := h.serv.UpdateUserActive(ctx, params.ID, params.Active, claims.UserID); err != nil {
		switch err {
		case service.ErrUserNotFound:
			return respond(c, http.StatusNotFound, "user not found", nil)
		case service.ErrUserCannotDeactivateSelf:
			return respond(c, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Println(err)
			return respond(c, http.StatusInternalServerError, err.Error(), nil)
		}
	}

	return respond(c, http.StatusOK, "user updated", nil)
}
