CREATE TABLE
    IF NOT EXISTS payout_service (
        payout_transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        partner_request_id UUID NOT NULL,
        operator_transaction_id TEXT NOT NULL,
        retailer_id TEXT NOT NULL,
        order_id TEXT NOT NULL,
        mobile_number TEXT NOT NULL,
        beneficiary_bank_name TEXT NOT NULL,
        beneficiary_name TEXT NOT NULL,
        beneficiary_account_number TEXT NOT NULL,
        beneficiary_ifsc_code TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        transfer_type TEXT NOT NULL CHECK (transfer_type IN ('5', '6')),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (retailer_id) REFERENCES retailers (retailer_id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS retailer_beneficiaries (
        retailer_id TEXT NOT NULL,
        mobile_number TEXT NOT NULL,
        beneficiary_id BIGSERIAL PRIMARY KEY,
        beneficiary_bank_name TEXT NOT NULL,
        beneficiary_name TEXT NOT NULL,
        beneficiary_account_number TEXT UNIQUE NOT NULL,
        beneficiary_ifsc_code TEXT NOT NULL,
        beneficiary_phone TEXT UNIQUE NOT NULL,
        is_beneficiary_verified BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (retailer_id) REFERENCES retailers (retailer_id) ON DELETE CASCADE
    );