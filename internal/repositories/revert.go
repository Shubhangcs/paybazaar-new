package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type RevertInterface interface {
	CreateRevert(echo.Context) error
}

type revertRepository struct {
	db *database.Database
}

func NewRevertRepository(db *database.Database) *revertRepository {
	return &revertRepository{
		db,
	}
}

func (rr *revertRepository) CreateRevert(c echo.Context) error {
	var req models.CreateRevertRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return rr.db.CreateRevertQuery(ctx, req)
}
