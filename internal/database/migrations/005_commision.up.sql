CREATE TABLE
    IF NOT EXISTS commisions (
        commision_id BIGSERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        service TEXT NOT NULL CHECK (service in ('PAYOUT', 'DMT', 'AEPS', 'BBPS')),
        total_commision NUMERIC(20, 2) NOT NULL CHECK (total_commision >= 0),
        admin_commision NUMERIC(20, 2) NOT NULL,
        master_distributor_commision NUMERIC(20, 2) NOT NULL,
        distributor_commision NUMERIC(20, 2) NOT NULL,
        retailer_commision NUMERIC(20, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (user_id, service)
    );

CREATE TABLE
    IF NOT EXISTS tds_commision (
        tds_commision_id BIGSERIAL NOT NULL,
        transaction_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        user_name TEXT NOT NULL,
        commision NUMERIC(20, 2) NOT NULL,
        tds NUMERIC(20, 2) NOT NULL,
        paid_commision NUMERIC(20, 2) NOT NULL,
        pan_number TEXT NOT NULL,
        status TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX IF NOT EXISTS idx_commisions_user_id ON commisions (user_id);

CREATE INDEX IF NOT EXISTS idx_commisions_created_at ON commisions (created_at);