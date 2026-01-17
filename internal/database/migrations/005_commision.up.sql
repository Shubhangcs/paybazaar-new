CREATE TABLE
    IF NOT EXISTS commisions (
        commision_id BIGSERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        total_commision NUMERIC(20, 2) NOT NULL,
        admin_commision NUMERIC(20, 2) NOT NULL,
        master_distributor_commision NUMERIC(20, 2) NOT NULL,
        distributor_commision NUMERIC(20, 2) NOT NULL,
        retailer_commision NUMERIC(20, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX IF NOT EXISTS idx_commisions_user_id ON commisions (user_id);

CREATE INDEX IF NOT EXISTS idx_commisions_created_at ON commisions (created_at);