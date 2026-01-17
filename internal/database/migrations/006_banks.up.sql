CREATE TABLE
    IF NOT EXISTS banks (
        bank_id BIGSERIAL PRIMARY KEY,
        bank_name TEXT NOT NULL,
        ifsc_code TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS admin_banks (
        admin_id TEXT NOT NULL,
        admin_bank_id BIGSERIAL PRIMARY KEY,
        admin_bank_name TEXT NOT NULL,
        admin_bank_account_number TEXT NOT NULL,
        admin_bank_ifsc_code TEXT NOT NULL
    );

