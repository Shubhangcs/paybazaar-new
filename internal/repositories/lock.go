package repositories

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/app"
)

type AdminLockRepository struct{}

func NewAdminLockRepository() *AdminLockRepository {
	return &AdminLockRepository{}
}

func (r *AdminLockRepository) LockAPI(c echo.Context) error {
	app.LockAPI()
	return c.JSON(http.StatusOK, map[string]string{
		"message": "All APIs locked successfully",
	})
}

func (r *AdminLockRepository) UnlockAPI(c echo.Context) error {
	app.UnlockAPI()
	return c.JSON(http.StatusOK, map[string]string{
		"message": "All APIs unlocked successfully",
	})
}
