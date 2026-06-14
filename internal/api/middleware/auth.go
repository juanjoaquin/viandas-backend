package middleware

import (
	"errors"
	"strings"

	"github.com/juanjoaquin/viandas-backend/encryption"
)

const bearerPrefix = "Bearer "

func ExtractToken(authHeader string) (string, error) {
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return strings.TrimSpace(authHeader[len(bearerPrefix):]), nil
	}
	return "", errors.New("missing or invalid authorization header")
}

func RequireRole(token string, roles ...string) (*encryption.Claims, error) {
	claims, err := encryption.ParseLoginJWT(token)
	if err != nil {
		return nil, err
	}
	for _, r := range roles {
		if claims.Role == r {
			return claims, nil
		}
	}
	return nil, errors.New("insufficient permissions")
}
