CREATE TABLE
    IF NOT EXISTS ticket (
        admin_id TEXT NOT NULL,
        ticket_id BIGSERIAL PRIMARY,
        ticket_title TEXT NOT NULL,
        ticket_description TEXT NOT NULL,
        is_ticket_cleared BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (admin_id) REFERENCES admins (admin_id) ON DELETE CASCADE
    );