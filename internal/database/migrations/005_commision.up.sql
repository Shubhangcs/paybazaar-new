CREATE TABLE
    IF NOT EXISTS commisions (
        commision_id UUID NOT NULL DEFAULT gen_random_uuid (),
        user_id TEXT NOT NULL,
        admin_commision TEXT NOT NULL,
        master_distributor_commision TEXT NOT NULL,
        distributor_commision TEXT NOT NULL,
        retailer_commision TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );