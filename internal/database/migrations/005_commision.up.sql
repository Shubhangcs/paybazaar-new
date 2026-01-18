CREATE TABLE
    IF NOT EXISTS commisions (
        commision_id BIGSERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        total_commision NUMERIC(20, 2) NOT NULL CHECK (total_commision >= 0),
        admin_commision NUMERIC(20, 2) NOT NULL CHECK (admin_commision >= 0),
        master_distributor_commision NUMERIC(20, 2) NOT NULL CHECK (master_distributor_commision >= 0),
        distributor_commision NUMERIC(20, 2) NOT NULL CHECK (distributor_commision >= 0),
        retailer_commision NUMERIC(20, 2) NOT NULL CHECK (retailer_commision >= 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        -- 🔐 Core business rule protection
        CONSTRAINT commision_split_valid CHECK (
            master_distributor_commision + distributor_commision + retailer_commision <= total_commision
        ),
        CONSTRAINT commision_admin_correct CHECK (
            admin_commision = total_commision - (
                master_distributor_commision + distributor_commision + retailer_commision
            )
        )
    );

CREATE INDEX IF NOT EXISTS idx_commisions_user_id ON commisions (user_id);

CREATE INDEX IF NOT EXISTS idx_commisions_created_at ON commisions (created_at);