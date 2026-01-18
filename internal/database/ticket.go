package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateTicketQuery(
	ctx context.Context,
	req models.CreateTicketModel,
) (int64, error) {

	query := `
		INSERT INTO ticket (
			admin_id,
			user_id,
			ticket_title,
			ticket_description
		) VALUES (
			@admin_id,
			@user_id,
			@title,
			@description
		)
		RETURNING ticket_id;
	`

	var ticketID int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_id":    req.AdminID,
		"user_id":     req.UserID,
		"title":       req.TicketTitle,
		"description": req.TicketDescription,
	}).Scan(&ticketID)

	return ticketID, err
}

func (db *Database) GetTicketByIDQuery(
	ctx context.Context,
	ticketID int64,
) (*models.TicketResponseModel, error) {

	query := `
		SELECT
			ticket_id,
			admin_id,
			user_id,
			ticket_title,
			ticket_description,
			is_ticket_cleared,
			created_at
		FROM ticket
		WHERE ticket_id = @ticket_id;
	`

	var t models.TicketResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"ticket_id": ticketID,
	}).Scan(
		&t.TicketID,
		&t.AdminID,
		&t.UserID,
		&t.TicketTitle,
		&t.TicketDescription,
		&t.IsTicketCleared,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *Database) GetTicketsByAdminIDQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.TicketResponseModel, error) {

	query := `
		SELECT
			ticket_id,
			admin_id,
			user_id,
			ticket_title,
			ticket_description,
			is_ticket_cleared,
			created_at
		FROM ticket
		WHERE admin_id = @admin_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.TicketResponseModel
	for rows.Next() {
		var t models.TicketResponseModel
		if err := rows.Scan(
			&t.TicketID,
			&t.AdminID,
			&t.UserID,
			&t.TicketTitle,
			&t.TicketDescription,
			&t.IsTicketCleared,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (db *Database) GetTicketsByUserIDQuery(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]models.TicketResponseModel, error) {

	query := `
		SELECT
			ticket_id,
			admin_id,
			user_id,
			ticket_title,
			ticket_description,
			is_ticket_cleared,
			created_at
		FROM ticket
		WHERE user_id = @user_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.TicketResponseModel
	for rows.Next() {
		var t models.TicketResponseModel
		if err := rows.Scan(
			&t.TicketID,
			&t.AdminID,
			&t.UserID,
			&t.TicketTitle,
			&t.TicketDescription,
			&t.IsTicketCleared,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (db *Database) GetAllTicketsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.TicketResponseModel, error) {

	query := `
		SELECT
			ticket_id,
			admin_id,
			user_id,
			ticket_title,
			ticket_description,
			is_ticket_cleared,
			created_at
		FROM ticket
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

	var tickets []models.TicketResponseModel
	for rows.Next() {
		var t models.TicketResponseModel
		if err := rows.Scan(
			&t.TicketID,
			&t.AdminID,
			&t.UserID,
			&t.TicketTitle,
			&t.TicketDescription,
			&t.IsTicketCleared,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (db *Database) UpdateTicketQuery(
	ctx context.Context,
	ticketID int64,
	req models.UpdateTicketModel,
) error {

	query := `
		UPDATE ticket
		SET
			ticket_title = COALESCE(@title, ticket_title),
			ticket_description = COALESCE(@description, ticket_description)
		WHERE ticket_id = @ticket_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"ticket_id":  ticketID,
		"title":      req.TicketTitle,
		"description": req.TicketDescription,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) UpdateTicketStatusQuery(
	ctx context.Context,
	ticketID int64,
	isCleared bool,
) error {

	query := `
		UPDATE ticket
		SET
			is_ticket_cleared = @is_cleared
		WHERE ticket_id = @ticket_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"ticket_id":  ticketID,
		"is_cleared": isCleared,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteTicketQuery(
	ctx context.Context,
	ticketID int64,
) error {

	tag, err := db.pool.Exec(ctx, `
		DELETE FROM ticket
		WHERE ticket_id = @ticket_id;
	`, pgx.NamedArgs{
		"ticket_id": ticketID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
