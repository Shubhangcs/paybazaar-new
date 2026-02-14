CREATE TABLE
    IF NOT EXISTS payout_transactions (
        payout_transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        partner_request_id UUID NOT NULL,
        operator_transaction_id TEXT NOT NULL,
        retailer_id TEXT NOT NULL,
        order_id TEXT NOT NULL,
        mobile_number TEXT NOT NULL,
        bank_name TEXT NOT NULL,
        beneficiary_name TEXT NOT NULL,
        account_number TEXT NOT NULL,
        ifsc_code TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        transfer_type TEXT NOT NULL CHECK (transfer_type IN ('IMPS', 'NEFT')),
        admin_commision NUMERIC(20, 2) NOT NULL,
        master_distributor_commision NUMERIC(20, 2) NOT NULL,
        distributor_commision NUMERIC(20, 2) NOT NULL,
        retailer_commision NUMERIC(20, 2) NOT NULL,
        payout_transaction_status TEXT NOT NULL CHECK (
            payout_transaction_status IN ('SUCCESS', 'PENDING', 'FAILED', 'REFUND')
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (retailer_id) REFERENCES retailers (retailer_id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS beneficiaries (
        beneficiary_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        mobile_number TEXT NOT NULL,
        bank_name TEXT NOT NULL,
        ifsc_code TEXT NOT NULL,
        account_number TEXT NOT NULL,
        beneficiary_name TEXT NOT NULL,
        beneficiary_phone TEXT NOT NULL,
        beneficiary_verified BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS mobile_recharge (
        mobile_recharge_transaction_id BIGSERIAL PRIMARY KEY,
        retailer_id TEXT NOT NULL REFERENCES retailers (retailer_id) ON DELETE CASCADE,
        partner_request_id TEXT NOT NULL,
        mobile_number TEXT NOT NULL,
        operator_name TEXT NOT NULL,
        circle_name TEXT NOT NULL,
        operator_code INTEGER NOT NULL,
        circle_code INTEGER NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        commision NUMERIC(20, 2) NOT NULL,
        recharge_type INTEGER NOT NULL,
        status TEXT NOT NULL CHECK (
            status IN ('SUCCESS', 'FAILED', 'PENDING', 'REFUND')
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS mobile_recharge_circles (
        circle_code INTEGER NOT NULL,
        circle_name TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS mobile_recharge_operators (
        operator_code INTEGER NOT NULL,
        operator_name TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS dth_recharge (
        dth_transaction_id BIGSERIAL PRIMARY KEY,
        retailer_id TEXT NOT NULL,
        partner_request_id TEXT NOT NULL,
        customer_id TEXT NOT NULL,
        operator_name TEXT NOT NULL,
        operator_code INTEGER NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        commision NUMERIC(20, 2) NOT NULL,
        status TEXT NOT NULL CHECK (
            status IN ('SUCCESS', 'PENDING', 'FAILED', 'REFUND')
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS dth_recharge_operators (
        operator_code INTEGER NOT NULL,
        operator_name TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS mobile_recharge_postpaid (
        postpaid_recharge_transaction_id BIGSERIAL PRIMARY KEY,
        retailer_id TEXT NOT NULL REFERENCES retailers (retailer_id) ON DELETE CASCADE,
        partner_request_id TEXT NOT NULL,
        operator_transaction_id TEXT,
        order_id TEXT,
        mobile_number TEXT NOT NULL,
        operator_code INTEGER NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        circle_code INTEGER NOT NULL,
        circle_name TEXT NOT NULL,
        operator_name TEXT NOT NULL,
        recharge_type TEXT NOT NULL,
        recharge_status TEXT NOT NULL CHECK (
            recharge_status IN ('PENDING', 'SUCCESS', 'FAILED' , 'REFUND')
        ),
        commision NUMERIC(20, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS electricity_bill_payments (
        electricity_bill_transaction_id BIGSERIAL PRIMARY KEY,
        retailer_id TEXT NOT NULL REFERENCES retailers (retailer_id) ON DELETE CASCADE,
        order_id TEXT,
        operator_transaction_id TEXT,
        partner_request_id TEXT NOT NULL,
        customer_id TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        operator_code INTEGER NOT NULL,
        operator_name TEXT NOT NULL,
        customer_email TEXT NOT NULL,
        commision NUMERIC(20, 2) NOT NULL,
        transaction_status TEXT NOT NULL CHECK (
            transaction_status IN ('PENDING', 'SUCCESS', 'FAILED' , 'REFUND')
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS electricity_operators (
        operator_name TEXT NOT NULL,
        operator_code INTEGER NOT NULL
    );