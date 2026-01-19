CREATE TABLE
    IF NOT EXISTS admins (
        admin_id TEXT PRIMARY KEY DEFAULT 'A' || LPAD(nextval('admin_id_sequence')::TEXT, 6, '0'),
        admin_name TEXT NOT NULL,
        admin_email TEXT UNIQUE NOT NULL,
        admin_phone TEXT UNIQUE NOT NULL,
        admin_password TEXT NOT NULL,
        admin_wallet_balance NUMERIC(20, 2) NOT NULL DEFAULT 0.0 CHECK (admin_wallet_balance >= 0.0),
        is_admin_blocked BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS master_distributors (
        master_distributor_id TEXT PRIMARY KEY DEFAULT 'M' || LPAD(nextval('md_id_sequence')::TEXT, 6, '0'),
        admin_id TEXT NOT NULL,
        master_distributor_name TEXT NOT NULL,
        master_distributor_phone TEXT UNIQUE NOT NULL,
        master_distributor_email TEXT UNIQUE NOT NULL,
        master_distributor_password TEXT NOT NULL,
        master_distributor_aadhar_number TEXT UNIQUE NOT NULL CHECK (length (master_distributor_aadhar_number) = 12),
        master_distributor_pan_number TEXT UNIQUE NOT NULL CHECK (length (master_distributor_pan_number) = 10),
        master_distributor_date_of_birth DATE NOT NULL,
        master_distributor_gender TEXT NOT NULL CHECK (
            master_distributor_gender IN ('MALE', 'FEMALE', 'OTHER')
        ),
        master_distributor_city TEXT NOT NULL,
        master_distributor_state TEXT NOT NULL,
        master_distributor_address TEXT NOT NULL,
        master_distributor_pincode TEXT NOT NULL,
        master_distributor_business_name TEXT NOT NULL,
        master_distributor_business_type TEXT NOT NULL,
        master_distributor_mpin INTEGER NOT NULL DEFAULT 1234,
        master_distributor_kyc_status BOOLEAN NOT NULL DEFAULT FALSE,
        master_distributor_documents_url TEXT,
        master_distributor_gst_number TEXT,
        master_distributor_wallet_balance NUMERIC(20, 2) NOT NULL DEFAULT 0.0 CHECK (master_distributor_wallet_balance >= 0.0),
        is_master_distributor_blocked BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (admin_id) REFERENCES admins (admin_id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS distributors (
        distributor_id TEXT PRIMARY KEY DEFAULT 'D' || LPAD(nextval('distributor_id_sequence')::TEXT, 6, '0'),
        master_distributor_id TEXT NOT NULL,
        distributor_name TEXT NOT NULL,
        distributor_phone TEXT UNIQUE NOT NULL,
        distributor_email TEXT UNIQUE NOT NULL,
        distributor_password TEXT NOT NULL,
        distributor_aadhar_number TEXT UNIQUE NOT NULL CHECK (length (distributor_aadhar_number) = 12),
        distributor_pan_number TEXT UNIQUE NOT NULL CHECK (length (distributor_pan_number) = 10),
        distributor_date_of_birth DATE NOT NULL,
        distributor_gender TEXT NOT NULL CHECK (distributor_gender IN ('MALE', 'FEMALE', 'OTHER')),
        distributor_city TEXT NOT NULL,
        distributor_state TEXT NOT NULL,
        distributor_address TEXT NOT NULL,
        distributor_pincode TEXT NOT NULL,
        distributor_business_name TEXT NOT NULL,
        distributor_business_type TEXT NOT NULL,
        distributor_gst_number TEXT,
        distributor_mpin INTEGER NOT NULL DEFAULT 1234,
        distributor_kyc_status BOOLEAN NOT NULL DEFAULT FALSE,
        distributor_documents_url TEXT,
        distributor_wallet_balance NUMERIC(20, 2) NOT NULL DEFAULT 0.0 CHECK (distributor_wallet_balance >= 0.0),
        is_distributor_blocked BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (master_distributor_id) REFERENCES master_distributors (master_distributor_id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS retailers (
        retailer_id TEXT PRIMARY KEY DEFAULT 'R' || LPAD(nextval('retailer_id_sequence')::TEXT, 6, '0'),
        distributor_id TEXT NOT NULL,
        retailer_name TEXT NOT NULL,
        retailer_phone TEXT UNIQUE NOT NULL,
        retailer_email TEXT UNIQUE NOT NULL,
        retailer_password TEXT NOT NULL,
        retailer_aadhar_number TEXT UNIQUE NOT NULL CHECK (length (retailer_aadhar_number) = 12),
        retailer_pan_number TEXT UNIQUE NOT NULL CHECK (length (retailer_pan_number) = 10),
        retailer_date_of_birth DATE NOT NULL,
        retailer_gender TEXT NOT NULL CHECK (retailer_gender IN ('MALE', 'FEMALE', 'OTHER')),
        retailer_city TEXT NOT NULL,
        retailer_state TEXT NOT NULL,
        retailer_address TEXT NOT NULL,
        retailer_pincode TEXT NOT NULL,
        retailer_business_name TEXT NOT NULL,
        retailer_business_type TEXT NOT NULL,
        retailer_gst_number TEXT,
        retailer_mpin INTEGER NOT NULL DEFAULT 1234,
        retailer_kyc_status BOOLEAN NOT NULL DEFAULT FALSE,
        retailer_documents_url TEXT,
        retailer_wallet_balance NUMERIC(20, 2) NOT NULL DEFAULT 0.0 CHECK (retailer_wallet_balance >= 0.0),
        is_retailer_blocked BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (distributor_id) REFERENCES distributors (distributor_id) ON DELETE CASCADE
    );