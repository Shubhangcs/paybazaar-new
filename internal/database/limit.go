package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateLimitQuery(
	ctx context.Context,
	req models.CreateTransactionLimitRequestModel,
) error {
	query := `
		INSERT INTO transaction_limit(
			retailer_id,
			limit_amount,
			service
		) VALUES (
			@retailer_id,
			@limit_amount,
			@service
		);
	`

	res, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id":  req.RetailerID,
		"limit_amount": req.LimitAmount,
		"service":      req.Service,
	})

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (db *Database) UpdateLimitQuery(
	ctx context.Context,
	req models.UpdateTransactionLimitRequestModel,
) error {
	query := `
		UPDATE transaction_limit
		SET limit_amount = @limit_amount,
		service = @service,
		updated_at = NOW()
		WHERE limit_id = @limit_id;
	`

	res, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"limit_amount": req.LimitAmount,
		"service":      req.Service,
		"limit_id":     req.LimitID,
	})

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteLimitQuery(
	ctx context.Context,
	limitId int,
) error {
	query := `
		DELETE FROM transaction_limit
		WHERE limit_id = @limit_id;
	`

	res, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"limit_id": limitId,
	})

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (db *Database) GetAllLimitsQuery(
	ctx context.Context,
) ([]models.GetLimitResponseModel, error) {
	query := `
		SELECT 
			limit_id,
			retailer_id,
			limit_amount,
			service,
			created_at,
			updated_at
		FROM transaction_limit;
	`

	res, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var limits []models.GetLimitResponseModel
	for res.Next() {
		var limit models.GetLimitResponseModel
		if err := res.Scan(
			&limit.LimitID,
			&limit.RetailerID,
			&limit.LimitAmount,
			&limit.Service,
			&limit.CreatedAT,
			&limit.UpdatedAT,
		); err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}

	fmt.Println(limits)

	return limits, res.Err()
}

func (db *Database) GetLimitAmountByRetailerIDAndServiceQuery(
	ctx context.Context,
	retailerId, service string,
) (float64, error) {
	query := `
		SELECT limit_amount
		FROM transaction_limit
		WHERE retailer_id = @retailer_id
		AND service = @service;
	`

	var limit float64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerId,
		"service":     service,
	}).Scan(
		&limit,
	)

	log.Println(err , "lolo")

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}
	return limit, nil
}

func (db *Database) GetLimitByRetailerIDServiceQuery(
	ctx context.Context,
	retailerId, service string,
) (*models.GetLimitResponseModel, error) {
	query := `
		SELECT 
			limit_id,
			retailer_id,
			limit_amount,
			service,
			created_at,
			updated_at
		FROM transaction_limit
		WHERE retailer_id = @retailer_id
		AND service = @service;
	`
	var limit models.GetLimitResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerId,
		"service":     service,
	}).Scan(
		&limit.LimitID,
		&limit.RetailerID,
		&limit.LimitAmount,
		&limit.Service,
		&limit.CreatedAT,
		&limit.UpdatedAT,
	)
	if err != nil {
		return nil, err
	}
	return &limit, nil
}
