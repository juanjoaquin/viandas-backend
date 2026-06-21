package handlers

import (
	"net/http"
	"strconv"

	"github.com/juanjoaquin/back-g-meta/pkg/meta"
	"github.com/labstack/echo/v5"
)

func parsePagination(c *echo.Context, defaultLimit string) (page, limit int) {
	page, _ = strconv.Atoi(c.QueryParam("page"))
	limit, _ = strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		if parsed, err := strconv.Atoi(defaultLimit); err == nil && parsed > 0 {
			limit = parsed
		} else {
			limit = 10
		}
	}
	return page, limit
}

func buildMeta(page, limit, total int, defaultLimit string) (*meta.Meta, error) {
	return meta.New(page, limit, total, defaultLimit)
}

func respondWithMeta(c *echo.Context, code int, message string, data interface{}, pagination *meta.Meta) error {
	return c.JSON(code, map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
		"meta":    pagination,
	})
}

func paginatedListResponse(
	c *echo.Context,
	defaultLimit string,
	countFn func() (int, error),
	listFn func(offset, limit int) (interface{}, error),
) error {
	page, limit := parsePagination(c, defaultLimit)

	total, err := countFn()
	if err != nil {
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	pagination, err := buildMeta(page, limit, total, defaultLimit)
	if err != nil {
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	data, err := listFn(pagination.Offset(), pagination.Limit())
	if err != nil {
		return respond(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return respondWithMeta(c, http.StatusOK, "ok", data, pagination)
}
