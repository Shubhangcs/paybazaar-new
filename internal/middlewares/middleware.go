package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/app"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

func AuthorizationMiddleware(jwtUtils *pkg.JwtUtils) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, models.ResponseModel{
					Message: "missing authorization header",
					Status:  "failed",
				})
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, models.ResponseModel{
					Message: "invalid authorization header format",
					Status:  "failed",
				})
			}
			tokenString := parts[1]
			claims, err := jwtUtils.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.ResponseModel{
					Message: "invalid access token or token is expired",
					Status:  "failed",
				})
			}
			c.Set("user", claims)
			return next(c)
		}
	}
}

// RequireRoles allows only specified roles to access a route
func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Get user claims from context
			user, ok := c.Get("user").(*models.AccessTokenClaims)
			if !ok || user == nil {
				return c.JSON(http.StatusUnauthorized, models.ResponseModel{
					Status:  "failed",
					Message: "unauthorized access",
				})
			}

			userRole := user.UserRole

			// Check if user role is allowed
			if slices.Contains(allowedRoles, userRole) {
				return next(c)
			}

			// Role not permitted
			return c.JSON(http.StatusForbidden, models.ResponseModel{
				Status:  "failed",
				Message: "you do not have permission to access this resource",
			})
		}
	}
}

func APILockMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Allow health + lock/unlock APIs
			path := c.Path()
			if path == "/admin/login" ||
				path == "/admin/portal/lock" ||
				path == "/admin/portal/unlock" {
				return next(c)
			}

			if app.IsLocked() {
				return c.JSON(http.StatusServiceUnavailable, map[string]any{
					"status":  "locked",
					"message": "System is under maintenance. Please try later.",
				})
			}

			return next(c)
		}
	}
}
