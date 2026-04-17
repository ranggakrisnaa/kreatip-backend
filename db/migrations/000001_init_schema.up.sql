-- =============================================================
-- 000001_init_schema.up.sql
-- Initial schema for Kreatip — donation platform for creators
--
-- Conventions:
--   - All timestamps use TIMESTAMPTZ (timezone-aware, WIB = UTC+7)
--   - All money columns use BIGINT in IDR (rupiah), never FLOAT
--   - Status enums use CHECK constraints (easier to extend vs PG ENUM)
--   - Soft deletes via deleted_at TIMESTAMPTZ (nullable)
--   - FKs use ON DELETE CASCADE only for owned data (tokens, profile)
-- =============================================================

-- pgcrypto provides gen_random_uuid() for Postgres < 13
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================
-- TABLE: users
-- =============================================================
CREATE TABLE users (
    id             UUID         NOT NULL DEFAULT gen_random_uuid(),
    email          VARCHAR(255) NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,
    role           VARCHAR(20)  NOT NULL DEFAULT 'user',
    email_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    totp_secret    VARCHAR(64),
    totp_enabled   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_uq UNIQUE (email),
    CONSTRAINT users_role_ck CHECK (role IN ('user', 'creator', 'admin'))
);

COMMENT ON TABLE  users IS 'Platform accounts. One user can be a creator (has creator_profiles row).';
COMMENT ON COLUMN users.role IS 'user | creator | admin';
COMMENT ON COLUMN users.totp_secret IS 'Base32 TOTP secret, stored encrypted at rest (todo).';

-- =============================================================
-- TABLE: creator_profiles
-- =============================================================
CREATE TABLE creator_profiles (
    id                UUID         NOT NULL DEFAULT gen_random_uuid(),
    user_id           UUID         NOT NULL,
    username          VARCHAR(30)  NOT NULL,
    display_name      VARCHAR(100) NOT NULL,
    avatar_url        VARCHAR(500),
    bio               TEXT,
    social_links      JSONB        NOT NULL DEFAULT '{}',
    theme_color       VARCHAR(7)   NOT NULL DEFAULT '#6366f1',
    min_donation      BIGINT       NOT NULL DEFAULT 1000,
    max_donation      BIGINT       NOT NULL DEFAULT 10000000,
    thank_you_message TEXT         NOT NULL DEFAULT 'Terima kasih atas dukungannya!',
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT creator_profiles_pkey        PRIMARY KEY (id),
    CONSTRAINT creator_profiles_user_id_fk  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT creator_profiles_user_id_uq  UNIQUE (user_id),
    CONSTRAINT creator_profiles_username_uq UNIQUE (username),
    CONSTRAINT creator_profiles_min_donation_ck   CHECK (min_donation >= 0),
    CONSTRAINT creator_profiles_max_donation_ck   CHECK (max_donation > 0),
    CONSTRAINT creator_profiles_donation_range_ck CHECK (max_donation >= min_donation),
    CONSTRAINT creator_profiles_theme_color_ck    CHECK (theme_color ~ '^#[0-9A-Fa-f]{6}$')
);

COMMENT ON TABLE  creator_profiles IS 'Public-facing profile for each creator. One-to-one with users.';
COMMENT ON COLUMN creator_profiles.social_links IS 'JSON: {"youtube":"","twitch":"","instagram":"","tiktok":""}';
COMMENT ON COLUMN creator_profiles.min_donation  IS 'Minimum donation amount in IDR set by the creator.';
COMMENT ON COLUMN creator_profiles.max_donation  IS 'Maximum donation amount in IDR (anti-abuse).';

-- =============================================================
-- TABLE: wallets
-- =============================================================
CREATE TABLE wallets (
    id                 UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id            UUID        NOT NULL,
    currency           VARCHAR(3)  NOT NULL DEFAULT 'IDR',
    balance_available  BIGINT      NOT NULL DEFAULT 0,
    balance_pending    BIGINT      NOT NULL DEFAULT 0,
    last_reconciled_at TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT wallets_pkey        PRIMARY KEY (id),
    CONSTRAINT wallets_user_id_fk  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT wallets_user_id_uq  UNIQUE (user_id),
    CONSTRAINT wallets_balance_available_ck CHECK (balance_available >= 0),
    CONSTRAINT wallets_balance_pending_ck   CHECK (balance_pending   >= 0)
);

COMMENT ON TABLE  wallets IS 'Cached balance per user. Source of truth is ledger_entries.';
COMMENT ON COLUMN wallets.balance_available IS 'Withdrawable balance (IDR). Cache — reconciled nightly against ledger.';
COMMENT ON COLUMN wallets.balance_pending   IS 'Amount held for in-flight withdrawals.';

-- =============================================================
-- TABLE: refresh_tokens
-- =============================================================
CREATE TABLE refresh_tokens (
    id         UUID         NOT NULL DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL,
    token_hash VARCHAR(64)  NOT NULL,
    user_agent VARCHAR(500) NOT NULL DEFAULT '',
    ip         VARCHAR(45)  NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ  NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT refresh_tokens_pkey          PRIMARY KEY (id),
    CONSTRAINT refresh_tokens_user_id_fk    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT refresh_tokens_token_hash_uq UNIQUE (token_hash)
);

COMMENT ON TABLE  refresh_tokens IS 'Rotated refresh tokens. token_hash = SHA-256 of the raw opaque token.';

-- =============================================================
-- TABLE: alert_tokens
-- =============================================================
CREATE TABLE alert_tokens (
    id         UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    token      VARCHAR(64) NOT NULL,
    settings   JSONB       NOT NULL DEFAULT '{}',
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT alert_tokens_pkey       PRIMARY KEY (id),
    CONSTRAINT alert_tokens_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT alert_tokens_token_uq   UNIQUE (token)
);

COMMENT ON TABLE  alert_tokens IS 'Opaque tokens for OBS browser-source overlay auth. Rotatable.';
COMMENT ON COLUMN alert_tokens.settings IS 'JSON: {"duration":5,"sound_url":"","template":"","animation":"slide"}';

-- =============================================================
-- TABLE: bank_accounts
-- =============================================================
CREATE TABLE bank_accounts (
    id                  UUID         NOT NULL DEFAULT gen_random_uuid(),
    user_id             UUID         NOT NULL,
    type                VARCHAR(10)  NOT NULL,
    bank_code           VARCHAR(20)  NOT NULL,
    account_number      VARCHAR(30)  NOT NULL,
    account_holder_name VARCHAR(100) NOT NULL,
    is_verified         BOOLEAN      NOT NULL DEFAULT FALSE,
    is_default          BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    CONSTRAINT bank_accounts_pkey       PRIMARY KEY (id),
    CONSTRAINT bank_accounts_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT bank_accounts_type_ck    CHECK (type IN ('bank', 'ewallet'))
);

COMMENT ON TABLE  bank_accounts IS 'Destination accounts for creator withdrawals. Soft-deleted.';
COMMENT ON COLUMN bank_accounts.bank_code IS 'e.g. BCA | BNI | BRI | MANDIRI | GOPAY | OVO | DANA | SHOPEEPAY';

-- =============================================================
-- TABLE: donations
-- =============================================================
CREATE TABLE donations (
    id           UUID         NOT NULL DEFAULT gen_random_uuid(),
    creator_id   UUID         NOT NULL,
    donor_name   VARCHAR(100),
    donor_email  VARCHAR(255),
    amount       BIGINT       NOT NULL,
    platform_fee BIGINT       NOT NULL DEFAULT 0,
    payment_fee  BIGINT       NOT NULL DEFAULT 0,
    net_amount   BIGINT       NOT NULL,
    message      TEXT,
    status       VARCHAR(20)  NOT NULL DEFAULT 'pending',
    currency     VARCHAR(3)   NOT NULL DEFAULT 'IDR',
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT donations_pkey          PRIMARY KEY (id),
    CONSTRAINT donations_creator_id_fk FOREIGN KEY (creator_id) REFERENCES creator_profiles(id),
    CONSTRAINT donations_status_ck     CHECK (status IN ('pending', 'paid', 'failed', 'expired', 'refunded')),
    CONSTRAINT donations_amount_ck     CHECK (amount > 0),
    CONSTRAINT donations_net_amount_ck CHECK (net_amount >= 0),
    CONSTRAINT donations_fees_ck       CHECK (platform_fee >= 0 AND payment_fee >= 0)
);

COMMENT ON TABLE  donations IS 'Every donation attempt. Soft-deleted on refund/expired cleanup.';
COMMENT ON COLUMN donations.amount       IS 'Gross amount paid by donor, in IDR.';
COMMENT ON COLUMN donations.platform_fee IS 'Kreatip platform fee (IDR).';
COMMENT ON COLUMN donations.payment_fee  IS 'Payment gateway fee passed to creator (IDR).';
COMMENT ON COLUMN donations.net_amount   IS 'amount - platform_fee - payment_fee. Credited to creator wallet.';

-- =============================================================
-- TABLE: payments
-- =============================================================
CREATE TABLE payments (
    id              UUID         NOT NULL DEFAULT gen_random_uuid(),
    donation_id     UUID         NOT NULL,
    provider        VARCHAR(30)  NOT NULL,
    method          VARCHAR(30)  NOT NULL,
    external_id     VARCHAR(255) NOT NULL DEFAULT '',
    external_status VARCHAR(50)  NOT NULL DEFAULT '',
    checkout_url    VARCHAR(500),
    raw_request     JSONB        NOT NULL DEFAULT '{}',
    raw_response    JSONB        NOT NULL DEFAULT '{}',
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT payments_pkey            PRIMARY KEY (id),
    CONSTRAINT payments_donation_id_fk  FOREIGN KEY (donation_id) REFERENCES donations(id),
    CONSTRAINT payments_donation_id_uq  UNIQUE (donation_id),
    CONSTRAINT payments_provider_ck     CHECK (provider IN ('xendit', 'midtrans'))
);

COMMENT ON TABLE  payments IS 'Payment gateway data per donation. 1-to-1 with donations.';
COMMENT ON COLUMN payments.method      IS 'qris | va_bca | va_bni | va_bri | va_mandiri | gopay | ovo | dana | shopeepay';
COMMENT ON COLUMN payments.raw_request  IS 'Serialized request body sent to provider (for audit/replay).';
COMMENT ON COLUMN payments.raw_response IS 'Raw response from provider (for audit).';

-- =============================================================
-- TABLE: ledger_entries
-- =============================================================
CREATE TABLE ledger_entries (
    id          UUID         NOT NULL DEFAULT gen_random_uuid(),
    wallet_id   UUID         NOT NULL,
    account     VARCHAR(50)  NOT NULL,
    direction   VARCHAR(10)  NOT NULL,
    amount      BIGINT       NOT NULL,
    ref_type    VARCHAR(30)  NOT NULL,
    ref_id      UUID         NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT ledger_entries_pkey        PRIMARY KEY (id),
    CONSTRAINT ledger_entries_wallet_id_fk FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    CONSTRAINT ledger_entries_direction_ck CHECK (direction IN ('debit', 'credit')),
    CONSTRAINT ledger_entries_account_ck  CHECK (account IN (
        'creator_balance',
        'pending_withdrawal',
        'platform_revenue',
        'payment_gateway',
        'pg_fee_expense'
    )),
    CONSTRAINT ledger_entries_ref_type_ck CHECK (ref_type IN ('donation', 'withdrawal', 'adjustment')),
    CONSTRAINT ledger_entries_amount_ck   CHECK (amount > 0)
);

COMMENT ON TABLE  ledger_entries IS 'Immutable double-entry ledger. Source of truth for all money movement.';
COMMENT ON COLUMN ledger_entries.account    IS 'Virtual account name. See entity/ledger_entry.go for constants.';
COMMENT ON COLUMN ledger_entries.direction  IS 'debit = money leaves account, credit = money enters account.';
COMMENT ON COLUMN ledger_entries.ref_type   IS 'donation | withdrawal | adjustment';
COMMENT ON COLUMN ledger_entries.ref_id     IS 'FK to donations.id or withdrawals.id depending on ref_type.';

-- =============================================================
-- TABLE: withdrawals
-- =============================================================
CREATE TABLE withdrawals (
    id              UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL,
    bank_account_id UUID        NOT NULL,
    amount          BIGINT      NOT NULL,
    fee             BIGINT      NOT NULL DEFAULT 0,
    net_amount      BIGINT      NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'requested',
    external_id     VARCHAR(255),
    failure_reason  TEXT,
    approved_by     UUID,
    approved_at     TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    CONSTRAINT withdrawals_pkey              PRIMARY KEY (id),
    CONSTRAINT withdrawals_user_id_fk        FOREIGN KEY (user_id)        REFERENCES users(id),
    CONSTRAINT withdrawals_bank_account_id_fk FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id),
    CONSTRAINT withdrawals_approved_by_fk    FOREIGN KEY (approved_by)    REFERENCES users(id),
    CONSTRAINT withdrawals_status_ck         CHECK (status IN (
        'requested', 'approved', 'processing', 'success', 'failed', 'rejected', 'cancelled'
    )),
    CONSTRAINT withdrawals_amount_ck     CHECK (amount > 0),
    CONSTRAINT withdrawals_net_amount_ck CHECK (net_amount >= 0),
    CONSTRAINT withdrawals_fee_ck        CHECK (fee >= 0)
);

COMMENT ON TABLE  withdrawals IS 'Creator payout requests. Flow: requested→approved→processing→success|failed.';
COMMENT ON COLUMN withdrawals.external_id    IS 'Xendit disbursement ID, set after Disburse() call.';
COMMENT ON COLUMN withdrawals.approved_by    IS 'Admin user ID who approved this withdrawal.';

-- =============================================================
-- TABLE: webhook_events
-- =============================================================
CREATE TABLE webhook_events (
    id              UUID         NOT NULL DEFAULT gen_random_uuid(),
    provider        VARCHAR(30)  NOT NULL,
    event_type      VARCHAR(50)  NOT NULL,
    external_id     VARCHAR(255) NOT NULL DEFAULT '',
    idempotency_key VARCHAR(255) NOT NULL,
    payload         JSONB        NOT NULL DEFAULT '{}',
    processed       BOOLEAN      NOT NULL DEFAULT FALSE,
    error           TEXT,
    received_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMPTZ,

    CONSTRAINT webhook_events_pkey               PRIMARY KEY (id),
    CONSTRAINT webhook_events_idempotency_key_uq UNIQUE (idempotency_key)
);

COMMENT ON TABLE  webhook_events IS 'Raw inbound webhooks log. Idempotency key = provider+external_id+event_type.';
COMMENT ON COLUMN webhook_events.idempotency_key IS 'Dedup key: "<provider>:<external_id>:<event_type>"';
COMMENT ON COLUMN webhook_events.processed       IS 'TRUE once business logic has been applied.';

-- =============================================================
-- TABLE: audit_logs
-- =============================================================
CREATE TABLE audit_logs (
    id            UUID        NOT NULL DEFAULT gen_random_uuid(),
    actor_user_id UUID,
    action        VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id   UUID,
    metadata      JSONB       NOT NULL DEFAULT '{}',
    ip            VARCHAR(45) NOT NULL DEFAULT '',
    occurred_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT audit_logs_pkey              PRIMARY KEY (id),
    CONSTRAINT audit_logs_actor_user_id_fk  FOREIGN KEY (actor_user_id) REFERENCES users(id) ON DELETE SET NULL
);

COMMENT ON TABLE  audit_logs IS 'Append-only audit trail. actor_user_id NULL = system action.';
COMMENT ON COLUMN audit_logs.action        IS 'e.g. login | logout | withdraw_request | admin_approve | ...';
COMMENT ON COLUMN audit_logs.resource_type IS 'e.g. user | donation | withdrawal';

-- =============================================================
-- INDEXES
-- =============================================================

-- users
CREATE INDEX idx_users_deleted_at        ON users (deleted_at) WHERE deleted_at IS NULL;

-- creator_profiles
CREATE INDEX idx_creator_profiles_active ON creator_profiles (is_active) WHERE is_active = TRUE;

-- donations
CREATE INDEX idx_donations_creator_status_created ON donations (creator_id, status, created_at DESC);
CREATE INDEX idx_donations_status_created          ON donations (status, created_at DESC);
CREATE INDEX idx_donations_deleted_at              ON donations (deleted_at) WHERE deleted_at IS NULL;

-- payments
CREATE INDEX idx_payments_external_id ON payments (external_id);

-- ledger_entries
CREATE INDEX idx_ledger_entries_wallet_occurred ON ledger_entries (wallet_id, occurred_at DESC);
CREATE INDEX idx_ledger_entries_ref             ON ledger_entries (ref_type, ref_id);

-- withdrawals
CREATE INDEX idx_withdrawals_user_id_status    ON withdrawals (user_id, status);
CREATE INDEX idx_withdrawals_status_created    ON withdrawals (status, created_at ASC);
CREATE INDEX idx_withdrawals_deleted_at        ON withdrawals (deleted_at) WHERE deleted_at IS NULL;

-- refresh_tokens
CREATE INDEX idx_refresh_tokens_user_id    ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);

-- alert_tokens
CREATE INDEX idx_alert_tokens_user_id ON alert_tokens (user_id);

-- bank_accounts
CREATE INDEX idx_bank_accounts_user_id    ON bank_accounts (user_id);
CREATE INDEX idx_bank_accounts_deleted_at ON bank_accounts (deleted_at) WHERE deleted_at IS NULL;

-- webhook_events
CREATE INDEX idx_webhook_events_provider_external ON webhook_events (provider, external_id);
CREATE INDEX idx_webhook_events_unprocessed        ON webhook_events (received_at) WHERE processed = FALSE;

-- audit_logs
CREATE INDEX idx_audit_logs_actor       ON audit_logs (actor_user_id);
CREATE INDEX idx_audit_logs_action      ON audit_logs (action);
CREATE INDEX idx_audit_logs_occurred_at ON audit_logs (occurred_at DESC);
