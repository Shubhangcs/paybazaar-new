package repositories

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func bindAndValidate(c echo.Context, req any) error {
	if c.Bind(req) != nil {
		return fmt.Errorf("invalid request body")
	}
	if err := c.Validate(req); err != nil {
		log.Println(err)
		return fmt.Errorf("invalid request format")
	}
	return nil
}

func parsePagination(c echo.Context) (limit int, offset int) {
	// Defaults
	const (
		defaultLimit = 100
		maxLimit     = 100
	)

	limit = defaultLimit
	page := 1

	// Parse limit
	if l := c.QueryParam("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			if v > maxLimit {
				limit = maxLimit
			} else {
				limit = v
			}
		}
	}

	// Parse page
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	offset = (page - 1) * limit
	return
}

func parseInt64Param(c echo.Context, key string) (int64, error) {
	idStr := c.Param(key)
	if idStr == "" {
		return 0, fmt.Errorf("param is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("param is required")
	}

	return id, nil
}
