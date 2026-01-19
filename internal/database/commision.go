package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateCommisionQuery(
	ctx context.Context,
	req models.CreateCommisionRequestModel,
) error {

	query := `
		INSERT INTO commisions (
			user_id,
			service,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		) VALUES (
			@user_id,
			@service,
			@total_commision,
			@admin_commision,
			@md_commision,
			@distributor_commision,
			@retailer_commision
		);
	`

	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id":               req.UserID,
		"service":               req.Service,
		"total_commision":       req.TotalCommision,
		"admin_commision":       req.AdminCommision,
		"md_commision":          req.MasterDistributorCommision,
		"distributor_commision": req.DistributorCommision,
		"retailer_commision":    req.RetailerCommision,
	}); err != nil {
		return fmt.Errorf("failed to create commision")
	}

	return nil
}

func (db *Database) GetCommisionDetailsByCommisionIDQuery(
	ctx context.Context,
	commisionID int64,
) (*models.GetCommisionResponseModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
			service,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		WHERE commision_id = @commision_id;
	`

	var c models.GetCommisionResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"commision_id": commisionID,
	}).Scan(
		&c.CommisionID,
		&c.UserID,
		&c.Service,
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch commision details")
	}

	return &c, nil
}

func (db *Database) GetCommisionsByUserIDQuery(
	ctx context.Context,
	userID string,
) ([]models.GetCommisionResponseModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
			service,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		WHERE user_id = @user_id;
	`

	res, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commisions")
	}
	defer res.Close()

	var list []models.GetCommisionResponseModel

	for res.Next() {
		var commision models.GetCommisionResponseModel
		if err := res.Scan(
			&commision.CommisionID,
			&commision.UserID,
			&commision.Service,
			&commision.TotalCommision,
			&commision.AdminCommision,
			&commision.MasterDistributorCommision,
			&commision.DistributorCommision,
			&commision.RetailerCommision,
			&commision.CreatedAt,
			&commision.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch commisions")
		}
		list = append(list, commision)
	}

	return list, res.Err()
}

func (db *Database) GetCommisionByUserIDAndServiceQuery(
	ctx context.Context,
	userID string,
	service string,
) (*models.GetCommisionResponseModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
			service,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		WHERE user_id = @user_id
		AND service = @service;
	`

	res := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"service": service,
	})

	var commision models.GetCommisionResponseModel

	if err := res.Scan(
		&commision.CommisionID,
		&commision.UserID,
		&commision.Service,
		&commision.TotalCommision,
		&commision.AdminCommision,
		&commision.MasterDistributorCommision,
		&commision.DistributorCommision,
		&commision.RetailerCommision,
		&commision.CreatedAt,
		&commision.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch commision")
	}

	return &commision, nil
}

func (db *Database) UpdateCommisionQuery(
	ctx context.Context,
	req models.UpdateCommisionRequestModel,
) error {
	query := `
		UPDATE commisions
		SET admin_commision = @admin_commision,
		master_distributor_commision = @md_commision,
		distributor_commision = @distributor_commision,
		retailer_commision = @retailer_commision,
		updated_at = NOW()
		WHERE commision_id = @commision_id;
	`
	res, err := db.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"admin_commision":              req.AdminCommision,
			"master_distributor_commision": req.MasterDistributorCommision,
			"distributor_commision":        req.DistributorCommision,
			"retailer_commision":           req.RetailerCommision,
			"commision_id":                 req.CommisionID,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to update commision")
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("invalid commision id or commision is not found")
	}

	return nil
}

func (db *Database) DeleteCommisionQuery(
	ctx context.Context,
	commisionID int64,
) error {

	tag, err := db.pool.Exec(ctx, `
		DELETE FROM commisions
		WHERE commision_id = @id;
	`, pgx.NamedArgs{"id": commisionID})

	if err != nil {
		return fmt.Errorf("failed to delete commision")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid commision or commision not found")
	}

	return nil
}
