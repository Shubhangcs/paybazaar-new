CREATE TABLE
    IF NOT EXISTS ticket (
        admin_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        ticket_id BIGSERIAL PRIMARY KEY,
        ticket_title TEXT NOT NULL,
        ticket_description TEXT NOT NULL,
        is_ticket_cleared BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (admin_id) REFERENCES admins (admin_id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_ticket_admin_id ON ticket (admin_id);

CREATE INDEX IF NOT EXISTS idx_ticket_created_at ON ticket (created_at);

CREATE INDEX IF NOT EXISTS idx_ticket_cleared ON ticket (is_ticket_cleared);