package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateCommisionQuery(
	ctx context.Context,
	req models.CreateCommisionModel,
) (int64, error) {

	// 1️⃣ Calculate split sum
	splitTotal :=
		req.MasterDistributorCommision +
			req.DistributorCommision +
			req.RetailerCommision

	// 2️⃣ Validate split does not exceed total
	if splitTotal > req.TotalCommision {
		return 0, fmt.Errorf(
			"invalid commission split: split total (%.2f) exceeds total commission (%.2f)",
			splitTotal,
			req.TotalCommision,
		)
	}

	// 3️⃣ Calculate admin commission (remaining)
	adminCommision := req.TotalCommision - splitTotal

	query := `
		INSERT INTO commisions (
			user_id,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		) VALUES (
			@user_id,
			@total,
			@admin,
			@md,
			@dist,
			@ret
		)
		RETURNING commision_id;
	`

	var id int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": req.UserID,
		"total":   req.TotalCommision,
		"admin":   adminCommision,
		"md":      req.MasterDistributorCommision,
		"dist":    req.DistributorCommision,
		"ret":     req.RetailerCommision,
	}).Scan(&id)

	return id, err
}

func (db *Database) GetCommisionByIDQuery(
	ctx context.Context,
	commisionID int64,
) (*models.GetCommisionModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		WHERE commision_id = @id;
	`

	var c models.GetCommisionModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": commisionID,
	}).Scan(
		&c.CommisionID,
		&c.UserID,
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *Database) GetCommisionByUserIDQuery(
	ctx context.Context,
	userID string,
) (*models.GetCommisionModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
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

	var c models.GetCommisionModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(
		&c.CommisionID,
		&c.UserID,
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *Database) GetAllCommisionsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetCommisionModel, error) {

	query := `
		SELECT
			commision_id,
			user_id,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetCommisionModel

	for rows.Next() {
		var c models.GetCommisionModel
		if err := rows.Scan(
			&c.CommisionID,
			&c.UserID,
			&c.TotalCommision,
			&c.AdminCommision,
			&c.MasterDistributorCommision,
			&c.DistributorCommision,
			&c.RetailerCommision,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, rows.Err()
}

func (db *Database) UpdateCommisionQuery(
	ctx context.Context,
	commisionID int64,
	req models.UpdateCommisionModel,
) error {

	// 1️⃣ Fetch existing values
	var existing struct {
		Total float64
		MD    float64
		Dist  float64
		Ret   float64
	}

	getQuery := `
		SELECT
			total_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		FROM commisions
		WHERE commision_id = @id;
	`

	err := db.pool.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"id": commisionID,
	}).Scan(
		&existing.Total,
		&existing.MD,
		&existing.Dist,
		&existing.Ret,
	)

	if err != nil {
		return err
	}

	// 2️⃣ Apply updates (COALESCE logic in Go)
	total := existing.Total
	if req.TotalCommision != nil {
		total = *req.TotalCommision
	}

	md := existing.MD
	if req.MasterDistributorCommision != nil {
		md = *req.MasterDistributorCommision
	}

	dist := existing.Dist
	if req.DistributorCommision != nil {
		dist = *req.DistributorCommision
	}

	ret := existing.Ret
	if req.RetailerCommision != nil {
		ret = *req.RetailerCommision
	}

	// 3️⃣ Validate split
	splitTotal := md + dist + ret
	if splitTotal > total {
		return fmt.Errorf(
			"invalid commission split: split total (%.2f) exceeds total commission (%.2f)",
			splitTotal,
			total,
		)
	}

	// 4️⃣ Calculate admin commission
	admin := total - splitTotal

	// 5️⃣ Update record
	updateQuery := `
		UPDATE commisions
		SET
			total_commision = @total,
			admin_commision = @admin,
			master_distributor_commision = @md,
			distributor_commision = @dist,
			retailer_commision = @ret,
			updated_at = NOW()
		WHERE commision_id = @id;
	`

	tag, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"id":    commisionID,
		"total": total,
		"admin": admin,
		"md":    md,
		"dist":  dist,
		"ret":   ret,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
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
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
