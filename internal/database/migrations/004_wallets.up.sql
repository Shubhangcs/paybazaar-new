CREATE TABLE
    IF NOT EXISTS wallet_transactions (
        wallet_transaction_id BIGSERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        reference_id TEXT NOT NULL,
        credit_amount NUMERIC(20, 2),
        debit_amount NUMERIC(20, 2),
        before_balance NUMERIC(20, 2) NOT NULL,
        after_balance NUMERIC(20, 2) NOT NULL,
        transaction_reason TEXT NOT NULL CHECK (
            transaction_reason IN (
                'FUND_TRANSFER',
                'FUND_REQUEST',
                'MOBILE_RECHARGE',
                'POSTPAID_MOBILE_RECHARGE',
                'MOBILE_RECHARGE_REFUND',
                'DTH_RECHARGE_REFUND',
                'PAYOUT_REFUND',
                'DTH_RECHARGE',
                'TOPUP',
                'REVERT',
                'PAYOUT'
            )
        ),
        remarks TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS revert_transactions (
        revert_transaction_id BIGSERIAL PRIMARY KEY,
        revert_by_id TEXT NOT NULL,
        revert_on_id TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        revert_status TEXT NOT NULL CHECK (revert_status IN ('PENDING', 'SUCCESS', 'FAILED')),
        remarks TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS fund_transfers (
        fund_transfer_id BIGSERIAL PRIMARY KEY,
        fund_transferer_id TEXT NOT NULL,
        fund_receiver_id TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        fund_transfer_status TEXT NOT NULL CHECK (
            fund_transfer_status IN ('PENDING', 'SUCCESS', 'FAILED')
        ),
        remarks TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS fund_requests (
        fund_request_id BIGSERIAL PRIMARY KEY,
        requester_id TEXT NOT NULL,
        request_to_id TEXT NOT NULL,
        amount NUMERIC(20, 2) NOT NULL,
        bank_name TEXT,
        request_date DATE NOT NULL,
        utr_number TEXT,
        request_type TEXT NOT NULL CHECK (request_type IN ('NORMAL', 'ADVANCE')),
        request_status TEXT NOT NULL CHECK (
            request_status IN ('PENDING', 'ACCEPTED', 'REJECTED')
        ),
        remarks TEXT NOT NULL,
        reject_remarks TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS transaction_limit (
        limit_id BIGSERIAL PRIMARY KEY,
        retailer_id TEXT REFERENCES retailers (retailer_id) ON DELETE CASCADE,
        limit_amount NUMERIC(20, 2) NOT NULL,
        service TEXT NOT NULL CHECK (service IN ('PAYOUT', 'DMT', 'AEPS')),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT unique_retailer_service UNIQUE (retailer_id, service)
    );